package main

import (
	"bytes"
	"ocala/internal/tt"
	"os"
	"path/filepath"
	"testing"
)

func setupTestCLI() *CLI {
	appRoot, _ = filepath.Abs("../..")
	return &CLI{outWriter: &bytes.Buffer{}, errWriter: &bytes.Buffer{}}
}

func TestRun(t *testing.T) {
	b := "../../internal/"
	d := "../../internal/core/testdata/"

	t.Run("ok", func(t *testing.T) {
		cli := setupTestCLI()
		data := []struct {
			expected int
			args     []string
		}{
			{0, []string{"-o", d + "output.bin", b + "z80/testdata/empty.oc"}},
			{0, []string{"-I", d + "z80/testdata", b + "z80/testdata/empty.oc"}},
			{0, []string{"-t", "z80", b + "z80/testdata/empty.oc"}},
			{0, []string{"-t", "6502", b + "mos6502/testdata/empty.oc"}},
			{0, []string{"-t", "z80", "-l", "-L", d + "output.lst", b + "z80/testdata/empty.oc"}},
			{0, []string{"-D", "CONST", b + "z80/testdata/empty.oc"}},
			{0, []string{"-V"}},
			{0, []string{"-h"}},
			{1, []string{}},
			{1, []string{"-error"}},
			{1, []string{"-t", "unknown", b + "z80/testdata/empty.oc"}},
			{1, []string{"-t", "z80", b + "mos6502/testdata/empty.oc"}},
			{1, []string{"-D", "!CONST", b + "z80/testdata/empty.oc"}},
		}
		for x, i := range data {
			actual := cli.Run(append([]string{"cmd"}, i.args...))
			tt.Eq(t, i.expected, actual, x, i.args)
		}
	})

	t.Run("ok: version", func(t *testing.T) {
		cli := setupTestCLI()
		actual := cli.Run([]string{"cmd", "-V"})
		tt.Eq(t, 0, actual)
		tt.Eq(t, "ocala dev\n", tt.FlushString(cli.outWriter))
	})

	t.Run("ok: findAppRoot", func(t *testing.T) {
		cli := setupTestCLI()
		cli.executable, _ = filepath.Abs("../cmd")
		appRoot = ""
		actual := cli.Run([]string{"cmd", b + "z80/testdata/empty.oc"})
		tt.Eq(t, 0, actual)
	})

	t.Run("error: findAppRoot", func(t *testing.T) {
		cli := setupTestCLI()
		cli.executable, _ = filepath.Abs("./cmd")
		appRoot = ""
		actual := cli.Run([]string{"cmd", b + "z80/testdata/empty.oc"})
		tt.Eq(t, 1, actual)
		tt.Eq(t, "invalid installation\n", tt.FlushString(cli.errWriter))
	})
}

func TestCompileExamples(t *testing.T) {
	cli := setupTestCLI()

	paths := []string{
		"../../examples/z80/msx-hello-world/main",
		"../../examples/z80/msx-hello-world-vdp/main",
		"../../examples/z80/msx-keytest/main",
		"../../examples/z80/msx-simple-game/main",
		"../../examples/z80/msx2-scroll/main",
	}
	for _, i := range paths {
		actual := cli.Run([]string{"cmd", i + ".oc"})
		tt.Eq(t, 0, actual, 0, tt.FlushString(cli.errWriter))

		a, _ := os.ReadFile(i + ".bin")
		tt.Eq(t, 1024*32, len(a))
	}
}

func TestFindAppRoot(t *testing.T) {
	cli := &CLI{outWriter: &bytes.Buffer{}, errWriter: &bytes.Buffer{}}

	cli.executable, _ = filepath.Abs("../../test.bin/test")
	err := cli.findAppRoot()
	tt.Eq(t, nil, err)

	cli.executable, _ = filepath.Abs("../../test.bin/test/test")
	err = cli.findAppRoot()
	tt.Eq(t, "invalid installation", err.Error())

	cli.executable = ""
	err = cli.findAppRoot()
	tt.Eq(t, "invalid installation", err.Error())
}
