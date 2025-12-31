package langserver

import (
	"context"
	"fmt"
	"io"
	"log"
	"ocala/core"
	"ocala/language-server/file"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

type Duration time.Duration

type Settings struct {
	Config Config `json:"ocala"`
}

// Config is
type Config struct {
	IncPaths []string `json:"incPaths"`
	Defs     []string `json:"defs"`
	Logger   *log.Logger
}

type langHandler struct {
	logger         *log.Logger
	files          map[DocumentURI]*file.File
	request        chan lintRequest
	lintDebounce   time.Duration
	lintTimer      *time.Timer
	conn           *jsonrpc2.Conn
	rootPath       string
	folders        []string
	rootMarkers    []string
	triggerChars   []string
	compileOptions *file.CompileOptions
	mu             sync.Mutex
}

type ReadWriteCloser struct {
	In  io.ReadCloser
	Out io.WriteCloser
}

func (rwc *ReadWriteCloser) Read(p []byte) (int, error) {
	return rwc.In.Read(p)
}

func (rwc *ReadWriteCloser) Write(p []byte) (int, error) {
	return rwc.Out.Write(p)
}

func (rwc *ReadWriteCloser) Close() error {
	if err := rwc.In.Close(); err != nil {
		return err
	}
	return rwc.Out.Close()
}

func Connect(rwc io.ReadWriteCloser, config *Config, connOpt []jsonrpc2.ConnOpt) {
	handler := &langHandler{
		logger:         config.Logger,
		files:          make(map[DocumentURI]*file.File),
		request:        make(chan lintRequest),
		lintDebounce:   300 * time.Millisecond,
		lintTimer:      nil,
		conn:           nil,
		rootMarkers:    []string{".git"},
		triggerChars:   []string{"."},
		compileOptions: file.NewCompileOptions(),
	}
	go handler.linter()

	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(rwc, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.HandlerWithError(handler.handle),
		connOpt...).DisconnectNotify()
}

func tokenToRange(t *core.Token) Range {
	if int(t.Line) >= len(t.From.Lines) {
		return Range{}
	}

	a := t.From.Lines[t.Line]
	b := t.From.Lines[t.End.Line]

	return Range{
		Start: Position{Line: int(t.Line), Character: int(t.Pos - a)},
		End:   Position{Line: int(t.End.Line), Character: int(t.End.Pos - b)},
	}
}

func positionToPt(s *file.File, r Position) core.Pt {
	if len(s.Text) == 0 {
		return core.Pt{}
	}
	if r.Line >= len(s.Lines)-1 {
		return core.Pt{Pos: int32(len(s.Text)), Line: int32(len(s.Lines) - 2)}
	}

	line := int32(r.Line)
	a, b := s.Lines[line], s.Lines[line+1]
	if pos := a + int32(r.Character); pos < b {
		return core.Pt{Pos: pos, Line: int32(line)}
	}
	if b > a && s.Text[b-1] == '\n' {
		b--
	}
	return core.Pt{Pos: b, Line: int32(line)}
}

func (h *langHandler) snapshot(f *file.File) ([]byte, *file.CompileOptions) {
	h.mu.Lock()
	options := *h.compileOptions
	h.mu.Unlock()
	return f.SafeGetText(), &options
}

func (h *langHandler) linter() {
	running := make(map[DocumentURI]context.CancelFunc)

	for {
		lintReq, ok := <-h.request
		if !ok {
			break
		}

		cancel, ok := running[lintReq.URI]
		if ok {
			cancel()
		}

		ctx, cancel := context.WithCancel(context.Background())
		running[lintReq.URI] = cancel

		go func() {
			f := h.files[lintReq.URI]
			diagnostics := []Diagnostic{}
			source := "syntax"

			text, options := h.snapshot(f)
			errors := file.CheckCode(f.Path, text, options)
			select {
			case <-ctx.Done():
				// canceled
			default:
				for _, i := range errors {
					for _, e := range file.ErrorInfoOf(i) {
						diagnostics = append(diagnostics, Diagnostic{
							Range:   tokenToRange(e.Token),
							Message: e.Message,
							Source:  &source,
						})
					}
				}
			}
			h.conn.Notify(
				ctx,
				"textDocument/publishDiagnostics",
				&PublishDiagnosticsParams{
					URI:         DocumentURI(lintReq.URI),
					Diagnostics: diagnostics,
				})
		}()
	}
}

func (h *langHandler) closeFile(uri DocumentURI) error {
	_, ok := h.files[uri]
	if !ok {
		return fmt.Errorf("document not found: %v", uri)
	}
	delete(h.files, uri)
	return nil
}

func (h *langHandler) saveFile(uri DocumentURI) error {
	_, ok := h.files[uri]
	if !ok {
		return fmt.Errorf("document not found: %v", uri)
	}
	h.lintRequest(uri, eventTypeSave)
	return nil
}

func (h *langHandler) openFile(uri DocumentURI, languageID string, version int) error {
	path, err := fromURI(uri)
	if err != nil {
		return err
	}

	f, ok := h.files[uri]
	if !ok {
		f = file.NewFile(path, version, h.compileOptions)
		h.files[uri] = f
	}
	f.Update([]byte{})
	return nil
}

func (h *langHandler) updateFile(uri DocumentURI, text string, version *int, eventType eventType) error {
	f, ok := h.files[uri]
	if !ok {
		return fmt.Errorf("document not found: %v", uri)
	}
	if version != nil {
		f.Version = *version
	}
	f.Update([]byte(text))
	h.lintRequest(uri, eventType)
	return nil
}

func (h *langHandler) addFolder(folder string) {
	folder = filepath.Clean(folder)
	if !slices.Contains(h.folders, folder) {
		h.folders = append(h.folders, folder)
	}
}

func (h *langHandler) addFile(uri DocumentURI, text []byte) (*file.File, error) {
	path, err := fromURI(uri)
	if err != nil {
		return nil, err
	}

	f := file.NewFile(path, 0, h.compileOptions)
	f.Update(text)
	h.files[uri] = f
	return f, nil
}

func (h *langHandler) addIncludedFile(path string) *file.File {
	uri := toURI(path)
	f, ok := h.files[uri]
	if !ok {
		text, err := os.ReadFile(path)
		if err != nil {
			h.logMessage(LogError, err.Error())
			return nil
		}

		f, err = h.addFile(uri, text)
		if err != nil {
			h.logMessage(LogError, err.Error())
			return nil
		}
	}
	f.Analyze(h.addIncludedFile)
	return f
}

func (h *langHandler) updateFileIncremental(uri DocumentURI, changes []TextDocumentContentChangeEvent, version *int, eventType eventType) error {
	f, ok := h.files[uri]
	if !ok {
		return fmt.Errorf("document not found: %v", uri)
	}
	if version != nil {
		f.Version = *version
	}

	for _, i := range changes {
		text := f.Scanner.Text
		buf := []byte{}
		start := positionToPt(f, i.Range.Start)
		end := positionToPt(f, i.Range.End)
		buf = append(buf, text[:start.Pos]...)
		buf = append(buf, []byte(i.Text)...)
		buf = append(buf, text[end.Pos:]...)
		f.Update(buf)
	}

	h.lintRequest(uri, eventType)
	return nil
}
