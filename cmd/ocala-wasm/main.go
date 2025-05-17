//go:build wasm

package main

import (
	"bytes"
	"ocala/internal/core"
	"ocala/internal/mos6502"
	"ocala/internal/z80"
	"syscall/js"
)

func main() {
	args := js.Global().Get("ocala")
	text := args.Get("text")
	core.OnComplete = args.Get("done")

	if len(core.IncMap) == 0 {
		core.IncMap = make(map[string][]byte, 8)
		incMap := args.Get("incMap")
		n := incMap.Length()
		for i := 0; i < n; i += 2 {
			k := incMap.Index(i).String()
			v := incMap.Index(i + 1).String()
			core.IncMap[k] = []byte(v)
		}
	}

	g := &core.Generator{
		InReader:  &bytes.Buffer{},
		OutWriter: &bytes.Buffer{},
		ErrWriter: &bytes.Buffer{},
		ListPath:  "/dev/null",
		GenList:   true,
		OutPath:   "-",
		Archs: map[string]func() *core.Compiler{
			"z80":     z80.BuildCompiler,
			"mos6502": mos6502.BuildCompiler,
		},
	}

	_, list := func(src []byte) ([]byte, []byte) {
		defer g.HandlePanic()

		g.SetCompilerFromSource(src)
		return g.GenerateBin(g.Compile("-", src))
	}([]byte(text.String()))
	core.OnComplete.Invoke(string(list), "")
}
