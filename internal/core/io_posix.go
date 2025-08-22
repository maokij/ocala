//go:build !wasm

package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func raiseError(err error) {
	panic(err)
}

func replacePathExt(path, ext string) string {
	return path[:len(path)-len(filepath.Ext(path))] + ext
}

var reInvalidPath = regexp.MustCompile(`^/|^\.\.+/|/\.+/|/\.*$|^\.+$`)

func regularizePath(s, dir string, paths []string) (string, error) {
	b := filepath.ToSlash(s)

	if filepath.IsAbs(s) || reInvalidPath.MatchString(b) {
		return "", fmt.Errorf("invalid path `%s`", s)
	}
	if strings.HasPrefix(b, "./") {
		paths = []string{dir}
	}
	for _, a := range paths {
		candidate := filepath.Join(a, b)
		if i, err := os.Stat(candidate); err == nil && !i.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("the file `%s` not found", s)
}

func (g *Generator) AppendIncPath(path string) error {
	a, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	g.IncPaths = append(g.IncPaths, a)
	return nil
}

func (g *Generator) CompileAndGenerate(path string) bool {
	defer g.HandlePanic()

	var err error
	var text []byte

	if path == "-" {
		text, err = io.ReadAll(g.InReader)
	} else {
		path, err = filepath.Abs(path)
		if err == nil {
			text, err = os.ReadFile(path)
		}
	}
	if err != nil {
		g.ErrorWith(err.Error())
	}

	if g.OutPath == "" {
		if strings.HasSuffix(path, ".oc") {
			g.OutPath = replacePathExt(path, ".bin")
		} else {
			g.ErrorWith("output file name required")
		}
	}

	if g.GenList && g.ListPath == "" {
		if strings.HasSuffix(path, ".oc") {
			g.ListPath = replacePathExt(path, ".lst")
		} else {
			g.ErrorWith("list file name required")
		}
	}

	if g.cc == nil {
		g.SetCompilerFromSource(text)
	}

	binary := g.GenerateBin(g.Compile(path, text))
	if g.OutPath == "-" {
		_, err = g.OutWriter.Write(binary)
	} else {
		err = os.WriteFile(g.OutPath, binary, 0o644)
	}
	if err == nil && g.GenList {
		err = os.WriteFile(g.ListPath, *g.ListText, 0o644)
	}
	if err != nil {
		g.ErrorWith(err.Error())
	}
	return true
}

// SPECIAL: (__FILE__)
func (cc *Compiler) sFilename(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 1, 1, CtConstexpr, cc)
	CheckPhase(PhCompile, etag, cc)

	path := cc.InPath
	if path != "-" {
		path = "./" + filepath.Base(path)
	}
	return NewStr(path)
}

// SYNTAX: (include path)
func (cc *Compiler) sInclude(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	path := CheckConst(e.At(1), StrT, "include path", etag, cc)

	CheckToplevelEnvIfCtProc(env, etag, cc)
	rpath, err := regularizePath(string(*path), filepath.Dir(cc.InPath), cc.g.IncPaths)
	if err != nil {
		cc.ErrorAt(etag).With(err.Error())
	}

	if slices.Index(cc.loaded, rpath) > -1 {
		return NIL
	}
	cc.loaded = append(cc.loaded, rpath)

	text, err := os.ReadFile(rpath)
	if err != nil {
		cc.ErrorAt(etag).With(err.Error())
	}

	return cc.CompileIncluded(rpath, text)
}

// SYNTAX: (load-file path)
func (cc *Compiler) sLoadFile(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	path := EvalConstAs(e.At(1), env, StrT, "path", etag, cc)
	CheckPhase(PhCompile, etag, cc)

	rpath, err := regularizePath(string(*path), filepath.Dir(cc.InPath), cc.g.IncPaths)
	if err != nil {
		cc.ErrorAt(etag).With(err.Error())
	}

	data, err := os.ReadFile(rpath)
	if err != nil {
		cc.ErrorAt(etag).With(err.Error())
	}

	return &Blob{data: data, path: rpath, origPath: string(*path)}
}

// SYNTAX: (compile-file path)
func sCompileFile(cc *Compiler, env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtConstexpr, cc)
	path := EvalConstAs(e.At(1), env, StrT, "path", etag, cc)

	blob := cc.sLoadFile(env, &Vec{etag, path}).(*Blob)
	g := CopyPtr(cc.g)
	g.IsSub = true
	g.SetCompiler(NewCompiler(cc.FullArchName()))
	data := g.GenerateBin(g.Compile(blob.path, blob.data))
	return &Blob{data: data, path: blob.path, origPath: blob.origPath, compiled: true}
}
