package main

import "os"

func main() {
	cli := &CLI{inReader: os.Stdin, outWriter: os.Stdout, errWriter: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
