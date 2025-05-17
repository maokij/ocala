package core

import (
	"io"
	"strings"
)

type InternalError struct {
	message string
	token   *Token
}

func (e *InternalError) Error() string {
	return e.message
}

type Generator struct {
	DebugMode bool
	Optimizer Optimizer
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer
	Archs     map[string]func() *Compiler
	IncPaths  []string
	Defs      []string
	GenList   bool
	ListPath  string
	OutPath   string
	Err       *InternalError
	changes   int
	cc        *Compiler
}

type Optimizer struct {
	OptimizeBCode func(*Compiler, *Inst, []BCode, bool) []BCode
}

func (g *Generator) RaiseGenerateError(message string, args ...any) {
	g.raiseError(nil, "generate error: ", message, args...)
}

func (g *Generator) Changed() {
	g.changes += 1
}

func (g *Generator) IsChanged() bool {
	return g.changes > 0
}

func (g *Generator) ErrorMessage() string {
	if g.Err != nil {
		return g.Err.Error()
	}
	return ""
}

func (g *Generator) ErrorMessageWithErrorLine() []byte {
	message := g.ErrorMessage()
	if g.Err != nil && g.Err.token != nil {
		message = FormatErrorLine(g.Err.token, true, message)
	}
	return []byte(strings.TrimRight(message, "\n") + "\n")
}

func (g *Generator) HandlePanic() {
	if err := recover(); err != nil {
		switch err := err.(type) {
		case *InternalError: // ok
			g.Err = err
		default:
			panic(err)
		}
	}
}

func (g *Generator) SetCompiler(cc *Compiler) {
	cc.g = g
	g.cc = cc
}

func (g *Generator) SetCompilerFromSource(text []byte) {
	arch := g.findArchDirective(text)
	if arch == "" {
		g.RaiseGenerateError("the first statement must be an `arch` directive unless the `-t` option is specified")
	}

	builder, ok := g.Archs[arch]
	if !ok {
		g.RaiseGenerateError("unknown arch: %s", arch)
	}
	g.SetCompiler(builder())
}

func (g *Generator) findArchDirective(text []byte) string {
	p := &Parser{Scanner: Scanner{Text: text}}

	p.seekToNextToken(false)
	if !p.Scan(reIdentifier) || p.Matched[1] != "" || p.Matched[2] != "arch" {
		return ""
	}

	_, _, nl := p.seekToNextToken(true)
	if nl || !p.Scan(reIdentifier) || p.Matched[1] != "" {
		return ""
	}

	return p.Matched[2]
}

func (g *Generator) Compile(path string, text []byte) []*Inst {
	return g.cc.Compile(path, text)
}
