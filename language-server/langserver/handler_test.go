package langserver

import (
	"context"
	"encoding/json"
	"net"
	"ocala/language-server/file"
	tt "ocala/testutil"
	"os"
	"testing"

	"github.com/sourcegraph/jsonrpc2"
)

type cs struct {
	ctx context.Context
	h   *langHandler
	c   *jsonrpc2.Conn
	s   *jsonrpc2.Conn
	fn  func(*jsonrpc2.Request) (any, error)
}

func (cs *cs) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if cs.fn != nil {
		return cs.fn(req)
	}
	return nil, nil
}

func newTestHandler() *langHandler {
	return &langHandler{
		files:          make(map[DocumentURI]*file.File),
		request:        make(chan lintRequest),
		rootMarkers:    []string{".git"},
		triggerChars:   []string{"."},
		compileOptions: file.NewCompileOptions(),
	}
}

func newConn() *cs {
	r, w := net.Pipe()
	ctx := context.Background()
	cs := &cs{ctx: ctx, h: newTestHandler()}
	cs.c = jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(w, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.HandlerWithError(cs.handle).SuppressErrClosed())
	cs.s = jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(r, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.HandlerWithError(cs.h.handle).SuppressErrClosed())
	cs.h.conn = cs.s
	return cs
}

func setFileText(h *langHandler, uri DocumentURI, text []byte) {
	tt.MustOk(h.openFile(uri, "ocala", 0))
	h.files[uri].Update(text)
}

func TestMain(m *testing.M) {
	file.Init("../../bin/a.out")
	m.Run()
}

func TestConnect(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		rwc := tt.NewBytesReadWriteCloser(nil)
		Connect(rwc, &Config{}, nil)
	})

	t.Run("ok: initialize", func(t *testing.T) {
		cs := newConn()
		result := InitializeResult{}
		err := cs.c.Call(cs.ctx, "initialize", InitializeParams{}, &result)
		tt.Eq(t, nil, err)
		tt.EqStruct(t, InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:           TDSKIncremental,
				DocumentFormattingProvider: false,
				RangeFormattingProvider:    false,
				DocumentSymbolProvider:     true,
				DefinitionProvider:         true,
				CompletionProvider: &CompletionProvider{
					TriggerCharacters: []string{".", ":"},
				},
				HoverProvider:      false,
				CodeActionProvider: false,
				Workspace: &ServerCapabilitiesWorkspace{
					WorkspaceFolders: WorkspaceFoldersServerCapabilities{
						Supported:           true,
						ChangeNotifications: true,
					},
				},
			},
		}, result)
	})

	t.Run("ok: initialized", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "initialized", nil, nil)
		tt.Eq(t, nil, err)
	})

	t.Run("ok: shutdown", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "shutdown", 0, nil)
		tt.Eq(t, jsonrpc2.ErrClosed, err)
	})

	t.Run("error: unsupported", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "unknown/unknown", nil, nil)
		tt.Eq(t, "jsonrpc2: code -32601 message: method not supported: unknown/unknown", err.Error())
	})

	t.Run("error: initialize", func(t *testing.T) {
		cs := newConn()
		result := InitializeResult{}
		err := cs.c.Call(cs.ctx, "initialize", nil, &result)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)
	})
}

func TestTextDocument(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "textDocument/didOpen", DidOpenTextDocumentParams{
			TextDocument: TextDocumentItem{
				URI: "file:///test",
			},
		}, nil)
		tt.Eq(t, nil, err)

		err = cs.c.Call(cs.ctx, "textDocument/didChange", DidChangeTextDocumentParams{
			TextDocument: VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: TextDocumentIdentifier{
					URI: "file:///test",
				},
			},
		}, nil)
		tt.Eq(t, nil, err)

		err = cs.c.Call(cs.ctx, "textDocument/didChange", DidChangeTextDocumentParams{
			TextDocument: VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: TextDocumentIdentifier{
					URI: "file:///test",
				},
			},
			ContentChanges: []TextDocumentContentChangeEvent{
				{
					Text: "test\ntest\n",
				},
			},
		}, nil)
		tt.Eq(t, nil, err)
		tt.Eq(t, "test\ntest\n", string(cs.h.files["file:///test"].Text))

		err = cs.c.Call(cs.ctx, "textDocument/didChange", DidChangeTextDocumentParams{
			TextDocument: VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: TextDocumentIdentifier{
					URI: "file:///test",
				},
			},
			ContentChanges: []TextDocumentContentChangeEvent{
				{
					Text: "abc\n",
					Range: &Range{
						Start: Position{Line: 0, Character: 0},
						End:   Position{Line: 0, Character: 0},
					},
				},
				{
					Text: "def\n",
					Range: &Range{
						Start: Position{Line: 2, Character: 2},
						End:   Position{Line: 2, Character: 4},
					},
				},
			},
		}, nil)
		tt.Eq(t, nil, err)
		tt.Eq(t, "abc\ntest\ntedef\n\n", string(cs.h.files["file:///test"].Text))

		err = cs.c.Call(cs.ctx, "textDocument/didSave", DidSaveTextDocumentParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///test",
			},
			Text: nil,
		}, nil)
		tt.Eq(t, nil, err)

		text := "didsavetext"
		err = cs.c.Call(cs.ctx, "textDocument/didSave", DidSaveTextDocumentParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///test",
			},
			Text: &text,
		}, nil)
		tt.Eq(t, nil, err)
		tt.Eq(t, "didsavetext", string(cs.h.files["file:///test"].Text))

		err = cs.c.Call(cs.ctx, "textDocument/didClose", DidCloseTextDocumentParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///test",
			},
		}, nil)
		tt.Eq(t, nil, err)
	})

	t.Run("error", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "textDocument/didOpen", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "textDocument/didOpen", DidOpenTextDocumentParams{
			TextDocument: TextDocumentItem{
				URI: "\x01",
			},
		}, nil)
		tt.Eq(t,
			"jsonrpc2: code 0 message: parse \"\\x01\": net/url: invalid control character in URL",
			err.Error())

		err = cs.c.Call(cs.ctx, "textDocument/didOpen", DidOpenTextDocumentParams{
			TextDocument: TextDocumentItem{
				URI: "invalid:///a",
			},
		}, nil)
		tt.Eq(t,
			"jsonrpc2: code 0 message: only file URIs are supported, got invalid",
			err.Error())

		err = cs.c.Call(cs.ctx, "textDocument/didChange", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "textDocument/didChange", DidChangeTextDocumentParams{
			TextDocument: VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: TextDocumentIdentifier{
					URI: "file:///unknown",
				},
			},
			ContentChanges: []TextDocumentContentChangeEvent{
				{
					Text: "test\ntest\n",
				},
			},
		}, nil)
		tt.Eq(t, "jsonrpc2: code 0 message: document not found: file:///unknown", err.Error())

		err = cs.c.Call(cs.ctx, "textDocument/didChange", DidChangeTextDocumentParams{
			TextDocument: VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: TextDocumentIdentifier{
					URI: "file:///unknown",
				},
			},
		}, nil)
		tt.Eq(t, "jsonrpc2: code 0 message: document not found: file:///unknown", err.Error())

		err = cs.c.Call(cs.ctx, "textDocument/didSave", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "textDocument/didSave", DidSaveTextDocumentParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///unknown",
			},
		}, nil)
		tt.Eq(t, "jsonrpc2: code 0 message: document not found: file:///unknown", err.Error())

		err = cs.c.Call(cs.ctx, "textDocument/didClose", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "textDocument/didClose", DidCloseTextDocumentParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///unknown",
			},
		}, nil)
		tt.Eq(t, "jsonrpc2: code 0 message: document not found: file:///unknown", err.Error())
	})
}

func TestWorkspace(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "workspace/didChangeConfiguration", DidChangeConfigurationParams{
			Settings: Config{},
		}, nil)
		tt.Eq(t, nil, err)

		err = cs.c.Call(cs.ctx, "workspace/didChangeWorkspaceFolders", DidChangeWorkspaceFoldersParams{
			Event: WorkspaceFoldersChangeEvent{
				Added: []WorkspaceFolder{
					{URI: "file:///test01"},
					{URI: "file:///test02"},
					{URI: "file:///test02"},
					{URI: "file:///test03"},
				},
			},
		}, nil)
		tt.Eq(t, nil, err)
		tt.Eq(t, 3, len(cs.h.folders))

		err = cs.c.Call(cs.ctx, "workspace/didChangeWorkspaceFolders", DidChangeWorkspaceFoldersParams{
			Event: WorkspaceFoldersChangeEvent{
				Removed: []WorkspaceFolder{
					{URI: "file:///test01"},
					{URI: "file:///test02"},
					{URI: "file:///test02"},
					{URI: "file:///test03"},
				},
			},
		}, nil)
		tt.Eq(t, nil, err)
		tt.Eq(t, 0, len(cs.h.folders))

		workspaces := []WorkspaceFolder{}
		tt.MustOk(cs.c.Call(cs.ctx, "initialize", InitializeParams{
			RootURI: "file:///rootpath/test",
		}, nil))
		err = cs.c.Call(cs.ctx, "workspace/workspaceFolders", 0, &workspaces)
		tt.Eq(t, nil, err)
		tt.EqSlice(t, []WorkspaceFolder{
			WorkspaceFolder{
				URI:  "file:///rootpath/test",
				Name: "test",
			},
		}, workspaces)
	})

	t.Run("error", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "workspace/didChangeConfiguration", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "workspace/didChangeWorkspaceFolders", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "workspace/workspaceFolders", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)
	})
}

func TestTextDocumentSymbol(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cs := newConn()
		setFileText(cs.h, "file:///test", tt.UnindentBytes(`
			arch z80
			const c001 = 1
			data d001 = byte
			proc p001() { return }
			macro m001() { }
			module n001 { }
		`))
		result := []DocumentSymbol{}
		err := cs.c.Call(cs.ctx, "textDocument/documentSymbol", DocumentSymbolParams{
			TextDocument: TextDocumentIdentifier{URI: "file:///test"},
		}, &result)
		tt.Eq(t, nil, err)
		tt.Eq(t, 5, len(result))
	})

	t.Run("error", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "textDocument/documentSymbol", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "textDocument/documentSymbol", DocumentSymbolParams{
			TextDocument: TextDocumentIdentifier{URI: "file:///unknown"},
		}, nil)
		tt.Eq(t, "jsonrpc2: code 0 message: document not found: file:///unknown", err.Error())
	})
}

func TestTextDocumentCompletion(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cs := newConn()
		setFileText(cs.h, "file:///test", tt.UnindentBytes(`
			arch z80
			const i001 = 1
			data i002 = byte
			proc i003() { return }
			macro i004() { }
			module i005 { }; i0
		`))
		result := []CompletionItem{}
		err := cs.c.Call(cs.ctx, "textDocument/completion", CompletionParams{
			TextDocumentPositionParams: TextDocumentPositionParams{
				TextDocument: TextDocumentIdentifier{URI: "file:///test"},
				Position:     Position{Line: 5, Character: 20},
			},
		}, &result)
		tt.Eq(t, nil, err)
		tt.Eq(t, 5, len(result))
	})

	t.Run("error", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "textDocument/completion", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "textDocument/completion", CompletionParams{
			TextDocumentPositionParams: TextDocumentPositionParams{
				TextDocument: TextDocumentIdentifier{URI: "file:///unknown"},
			},
		}, nil)
		tt.Eq(t, "jsonrpc2: code 0 message: document not found: file:///unknown", err.Error())
	})
}

func TestTextDocumentDefinition(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cs := newConn()
		setFileText(cs.h, "file:///test", tt.UnindentBytes(`
			arch z80
			const c001 = 1
			data d001 = byte
			proc p001() { return }
			macro m001() { }
			module n001 { }; d001
		`))
		result := []Location{}
		err := cs.c.Call(cs.ctx, "textDocument/definition", DocumentDefinitionParams{
			TextDocumentPositionParams: TextDocumentPositionParams{
				TextDocument: TextDocumentIdentifier{URI: "file:///test"},
				Position:     Position{Line: 5, Character: 20},
			},
		}, &result)
		tt.Eq(t, nil, err)
		tt.Eq(t, 1, len(result))
		tt.Eq(t, 2, result[0].Range.Start.Line)
		tt.Eq(t, 5, result[0].Range.Start.Character)
	})

	t.Run("error", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "textDocument/definition", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)

		err = cs.c.Call(cs.ctx, "textDocument/definition", DocumentDefinitionParams{
			TextDocumentPositionParams: TextDocumentPositionParams{
				TextDocument: TextDocumentIdentifier{URI: "file:///unknown"},
			},
		}, nil)
		tt.Eq(t, "jsonrpc2: code 0 message: document not found: file:///unknown", err.Error())
	})
}

func TestTextDocumentSymbolHover(t *testing.T) {
	t.Run("ok: textDocument/hover", func(t *testing.T) {
		cs := newConn()
		setFileText(cs.h, "file:///test", tt.UnindentBytes(`
			arch z80
			const i001 = 1
			data i002 = byte
			proc i003() { return }
			macro i004() { }
			module i005 { }; i0
		`))
		err := cs.c.Call(cs.ctx, "textDocument/hover", HoverParams{
			TextDocumentPositionParams: TextDocumentPositionParams{
				TextDocument: TextDocumentIdentifier{URI: "file:///test"},
				Position:     Position{Line: 2, Character: 8},
			},
		}, nil)
		tt.Eq(t, nil, err)
	})

	t.Run("error", func(t *testing.T) {
		cs := newConn()
		err := cs.c.Call(cs.ctx, "textDocument/hover", nil, nil)
		tt.Eq(t, jsonrpc2.CodeInvalidParams, err.(*jsonrpc2.Error).Code)
	})
}

func TestLint(t *testing.T) {
	cs := newConn()
	go cs.h.linter()
	defer close(cs.h.request)

	var params PublishDiagnosticsParams
	var clientErr error
	c := make(chan bool)
	cs.fn = func(req *jsonrpc2.Request) (any, error) {
		switch req.Method {
		case "textDocument/publishDiagnostics":
			if err := json.Unmarshal(*req.Params, &params); err != nil {
				clientErr = err
			}
			c <- true
		}
		return nil, nil
	}
	text := tt.Must(os.ReadFile("../file/testdata/test2.oc"))
	setFileText(cs.h, "file:///test", text)
	err := cs.c.Call(cs.ctx, "textDocument/didChange", DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{
				URI: "file:///test",
			},
		},
	}, nil)
	tt.Eq(t, nil, err)

	<-c
	tt.Eq(t, nil, clientErr)
	tt.True(t, len(params.Diagnostics) > 0)
	t.Log("?", len(params.Diagnostics))
}

func TestLogMessage(t *testing.T) {
	cs := newConn()
	var params LogMessageParams
	var clientErr error
	cs.fn = func(req *jsonrpc2.Request) (any, error) {
		switch req.Method {
		case "window/logMessage":
			if err := json.Unmarshal(*req.Params, &params); err != nil {
				clientErr = err
			}

		}
		return nil, nil
	}
	cs.h.logMessage(LogInfo, "test")
	cs.s.Close()
	<-cs.c.DisconnectNotify()
	tt.Eq(t, nil, clientErr)
	tt.Eq(t, "test", params.Message)
}
