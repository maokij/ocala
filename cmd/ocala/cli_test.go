package main

import (
	"bytes"
	"hash/crc32"
	tt "ocala/testutil"
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
	d := "../../core/testdata/"

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
	es := []struct {
		expect uint32
		name   string
	}{
		{0x2f378a54, "../../examples/z80/msx-colortext/main"},
		{0x7638c5b1, "../../examples/z80/msx-hello-world-bsave/main"},
		{0x5d0d5139, "../../examples/z80/msx-hello-world-com/main"},
		{0x09541103, "../../examples/z80/msx-hello-world-vdp/main"},
		{0x6a96ebe2, "../../examples/z80/msx-hello-world/main"},
		{0x297dd6d6, "../../examples/z80/msx-keytest/main"},
		{0xd1e86ef0, "../../examples/z80/msx-simple-game/main"},
		{0x8ca3a036, "../../examples/z80/msx-sprite32/main"},
		{0x42ac396e, "../../examples/z80/msx2-scroll/main"},
	}
	for _, i := range es {
		r := cli.Run([]string{"cmd", i.name + ".oc"})
		tt.Eq(t, 0, r, i.name, tt.FlushString(cli.errWriter))

		dat, _ := os.ReadFile(i.name + ".bin")
		actual := crc32.ChecksumIEEE(dat)
		tt.Eq(t, i.expect, actual, i.name)
	}
}
