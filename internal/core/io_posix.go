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

func replacePathExt(path, ext string) string {
	return path[:len(path)-len(filepath.Ext(path))] + ext
}

var reInvalidPath = regexp.MustCompile(`^/|^\.\.+/|/\.+/|/\.*$|^\.+$`)

func regularizePath(s, dir string, paths []string) (string, error) {
	path := filepath.ToSlash(s)

	if reInvalidPath.MatchString(path) {
		return "", fmt.Errorf("invalid path `%s`", s)
	}
	if strings.HasPrefix(path, "./") {
		paths = []string{dir}
	}
	for _, dir := range paths {
		candidate := filepath.Join(dir, path)
		if i, err := os.Stat(candidate); err == nil && !i.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("the file `%s` not found", s)
}

func (g *Generator) raiseError(token *Token, tag string, message string, args ...any) {
	err := &InternalError{message: tag + fmt.Sprintf(message, args...), token: token}
	panic(err)
}

func (g *Generator) AppendIncPath(path string) error {
	a, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	g.IncPaths = append(g.IncPaths, a)
	return nil
}

func generateFiles(g *Generator, insts []*Inst) {
	binary, list := g.GenerateBin(insts)

	var err error
	if g.OutPath == "-" {
		_, err = g.OutWriter.Write(binary)
	} else {
		err = os.WriteFile(g.OutPath, binary, 0o644)
	}
	if err == nil && g.GenList {
		err = os.WriteFile(g.ListPath, []byte(list), 0o644)
	}
	if err != nil {
		g.RaiseGenerateError(err.Error())
	}
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
		if err == nil && g.OutPath == "" && strings.HasSuffix(path, ".oc") {
			g.OutPath = replacePathExt(path, ".bin")
		}
	}
	if err != nil {
		g.RaiseGenerateError(err.Error())
	}

	if g.OutPath == "" {
		g.RaiseGenerateError("output file name required")
	}

	if g.GenList && g.ListPath == "" {
		if strings.HasSuffix(path, ".oc") {
			g.ListPath = replacePathExt(path, ".lst")
		} else {
			g.RaiseGenerateError("list file name required")
		}
	}

	if g.cc == nil {
		g.SetCompilerFromSource(text)
	}

	generateFiles(g, g.Compile(path, text))
	return true
}

// SYNTAX: (include path)
func (cc *Compiler) sInclude(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	path := CheckConst(e.At(1), StrT, "include path", etag, cc)

	if cc.Context() == CtProc {
		CheckToplevelEnv(env, etag, cc)
	}

	rpath, err := regularizePath(string(*path), filepath.Dir(cc.InPath), cc.g.IncPaths)
	if err != nil {
		cc.RaiseCompileError(etag, err.Error())
	}

	if slices.Index(cc.loaded, rpath) > -1 {
		return NIL
	}
	cc.loaded = append(cc.loaded, rpath)

	text, err := os.ReadFile(rpath)
	if err != nil {
		cc.RaiseCompileError(etag, err.Error())
	}

	return cc.CompileIncluded(rpath, text)
}

// SYNTAX: (embed-file path)
func (cc *Compiler) sEmbedFile(env *Env, e *Vec) Value {
	etag, _ := CheckExpr(e, 2, 2, CtModule|CtProc, cc)
	path := CheckConst(e.At(1), StrT, "include path", etag, cc)

	rpath, err := regularizePath(string(*path), filepath.Dir(cc.InPath), cc.g.IncPaths)
	if err != nil {
		cc.RaiseCompileError(etag, err.Error())
	}

	data, err := os.ReadFile(rpath)
	if err != nil {
		cc.RaiseCompileError(etag, err.Error())
	}
	return cc.EmitCode(NewInst(e, InstFile, path, (*Binary)(&data)))
}
