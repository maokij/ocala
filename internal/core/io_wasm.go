//go:build wasm

package core

import (
	"fmt"
	"strings"
	"syscall/js"
)

var IncMap map[string][]byte
var OnComplete js.Value

func (g *Generator) raiseError(token *Token, tag string, message string, args ...any) {
	message = tag + fmt.Sprintf(message, args...)
	if token != nil {
		message = FormatErrorLine(token, true, message)
	}
	OnComplete.Invoke("", strings.TrimRight(message, "\n")+"\n")
	panic("error")
}

func (g *Generator) AppendIncPath(path string) error {
	return nil
}

func (g *Generator) CompileAndGenerate(path string) bool {
	return false
}

// SPECIAL: (__FILE__)
func (cc *Compiler) sFilename(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	cc.RaiseCompileError(etag, "not supported for browser wasm")
	return NIL
}

// SYNTAX: (include path)
func (cc *Compiler) sInclude(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	path := CheckConst(e.At(1), StrT, "include path", etag, cc)

	rpath := path.String()
	text, ok := IncMap[rpath]
	if !ok {
		cc.RaiseCompileError(etag, "`%s` not found", rpath)
	}

	return cc.CompileIncluded(rpath, text)
}

// SYNTAX: (load-file path)
func (cc *Compiler) sLoadFile(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	cc.RaiseCompileError(etag, "not supported for browser wasm")
	return NIL
}

// SYNTAX: (compile-file path)
func sCompileFile(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	cc.RaiseCompileError(etag, "not supported for browser wasm")
	return NIL
}
