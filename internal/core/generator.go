package core

import (
	"fmt"
	"io"
	"runtime"
)

type InternalError struct {
	message   string
	tag       string
	at        []Value
	DebugMode bool
}

func (err *InternalError) Error() string {
	return err.tag + err.message
}

func (err *InternalError) With(message string, args ...any) {
	err.message = fmt.Sprintf(message, args...)
	if err.DebugMode {
		for i := 1; ; i++ {
			if _, file, line, ok := runtime.Caller(i); ok {
				err.message += fmt.Sprintf("\n-- %s:%d", file, line)
				continue
			}
			break
		}
	}
	raiseError(err)
}

type Generator struct {
	DebugMode bool
	GenList   bool
	IsSub     bool
	Optimizer Optimizer
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer
	Archs     map[string]func() *Compiler
	IncPaths  []string
	Defs      []string
	ListText  *[]byte
	ListPath  string
	OutPath   string
	Err       *InternalError
	changes   int
	cc        *Compiler
}

type Optimizer struct {
	OptimizeBCode func(*Compiler, *Inst, []BCode, bool) []BCode
}

func (g *Generator) ErrorAt(values ...Value) *InternalError {
	return &InternalError{
		tag:       "generate error: ",
		at:        values,
		DebugMode: g.DebugMode,
	}
}

func (g *Generator) ErrorWith(message string, args ...any) {
	g.ErrorAt().With(message, args...)
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

func FindToken(v Value) *Token {
	switch v := v.(type) {
	case *Vec:
		return FindToken(v.AtOrUndef(0))
	case *Inst:
		return FindToken(v.From.AtOrUndef(0))
	case *Operand:
		return FindToken(v.From)
	case *Token:
		return v
	case *Identifier:
		return v.Token
	case *Constexpr:
		return v.Token
	case *Named:
		return v.Token
	}
	return nil
}

func (g *Generator) FullErrorMessage() []byte {
	if g.Err == nil {
		return []byte{}
	}

	message := g.ErrorMessage() + string('\n')
	x := 0
	for _, i := range g.Err.at {
		if token := FindToken(i); token != nil {
			message += fmt.Sprintf("[error #%d]\n", x)
			message += token.FormatAsErrorLine("at")
			if id, ok := token.Value.(*Identifier); ok {
				for ; id.ExpandedBy != nil; id = id.ExpandedBy {
					message += id.ExpandedBy.Token.FormatAsErrorLine("from")
				}
			}
			x++
		}
	}
	return []byte(message)
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
		g.ErrorWith("the first statement must be an `arch` directive unless the `-t` option is specified")
	}

	builder, ok := g.Archs[arch]
	if !ok {
		g.ErrorWith("unknown arch: %s", arch)
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

func (g *Generator) prependList(list []byte) {
	*g.ListText = append(list, *g.ListText...)
}

func (g *Generator) Compile(path string, text []byte) []*Inst {
	return g.cc.Compile(path, text)
}
