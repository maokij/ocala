package main

import (
	"io"
	"log"
	"ocala/language-server/file"
	"ocala/language-server/langserver"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

type CLI struct {
	inReader   io.ReadCloser
	outWriter  io.WriteCloser
	errWriter  io.Writer
	executable string
}

func (cli *CLI) Run() int {
	file.Init(cli.executable)

	config := &langserver.Config{
		Logger: log.New(cli.errWriter, "", log.LstdFlags),
	}
	connOpt := []jsonrpc2.ConnOpt{}
	rwc := &langserver.ReadWriteCloser{In: cli.inReader, Out: cli.outWriter}
	langserver.Connect(rwc, config, connOpt)
	return 0
}

func main() {
	cli := &CLI{inReader: os.Stdin, outWriter: os.Stdout, errWriter: os.Stderr}
	os.Exit(cli.Run())
}
