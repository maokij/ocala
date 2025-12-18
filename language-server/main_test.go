package main

import (
	tt "ocala/testutil"
	"testing"
)

func setupCLI(s string) *CLI {
	return &CLI{
		inReader:   tt.NewBytesReadWriteCloser([]byte(s)),
		outWriter:  tt.NewBytesReadWriteCloser(nil),
		errWriter:  tt.NewBytesReadWriteCloser(nil),
		executable: "../bin/a.out",
	}
}

func TestRun(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cli := setupCLI("")
		tt.Eq(t, 0, cli.Run())
	})
}
