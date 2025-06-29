package main

import (
	"flag"
	"fmt"
	"io"
	"ocala/internal/core"
	"ocala/internal/mos6502"
	"ocala/internal/z80"
	"os"
	"path/filepath"
	"regexp"
)

var version = "dev"
var appRoot = ""

var errNoInputFile = fmt.Errorf("input file required(-h for help)")
var errInvalidTarget = fmt.Errorf("invalid target arch(-h for help)")
var errInvalidInstallation = fmt.Errorf("invalid installation")

var Archs = map[string]func() *core.Compiler{
	"z80":              z80.BuildCompiler,
	"z80+undocumented": z80.BuildCompilerUndocumented,
	"6502":             mos6502.BuildCompiler,
	"mos6502":          mos6502.BuildCompiler,
}

type CLI struct {
	inReader   io.Reader
	outWriter  io.Writer
	errWriter  io.Writer
	executable string
}

type StringListValue []string

func (v *StringListValue) String() string {
	return fmt.Sprintf("%v", *v)
}

func (v *StringListValue) Set(s string) error {
	*v = append(*v, s)
	return nil
}

func (cli *CLI) findAppRoot() error {
	path, err := cli.executable, error(nil)
	if path == "" {
		path, err = os.Executable()
		if err != nil {
			return err
		}

		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
	}

	path = filepath.Join(path, "../..")
	stat, err := os.Stat(filepath.Join(path, "share/ocala/include"))
	if err != nil || !stat.IsDir() {
		return errInvalidInstallation
	}

	appRoot = path
	return nil
}

var reDefName = regexp.MustCompile(`^[_A-Za-z][-_A-Za-z0-9]*$`)

func (*CLI) ParseCommandLineOptions(g *core.Generator, args []string) (string, error) {
	arch := ""
	incPaths := []string{}
	printVersion := false

	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.SetOutput(g.ErrWriter)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "Usage: %s [options] file\nOptions:\n", args[0])
		flags.PrintDefaults()
	}
	flags.StringVar(&arch, "t", "", "Specify the target arch")
	flags.Var((*StringListValue)(&incPaths), "I", "Add the directory to the include path")
	flags.Var((*StringListValue)(&g.Defs), "D", "Define the symbol")
	flags.StringVar(&g.OutPath, "o", "", "Specify the output file name")
	flags.StringVar(&g.ListPath, "L", "", "Specify the list file name")
	flags.BoolFunc("l", "Generate a list file", func(string) error {
		g.GenList = true
		return nil
	})
	flags.BoolVar(&printVersion, "V", false, "Display the version information")

	if err := flags.Parse(args[1:]); err != nil {
		return "", err
	} else if printVersion {
		fmt.Fprintf(g.OutWriter, "ocala %s\n", version)
		return "", flag.ErrHelp
	}
	if len(flags.Args()) == 0 {
		err := errNoInputFile
		fmt.Fprintln(g.ErrWriter, err.Error())
		return "", err
	}

	if arch != "" {
		builder, ok := Archs[arch]
		if !ok {
			err := errInvalidTarget
			fmt.Fprintln(g.ErrWriter, err.Error())
			return "", err
		}
		g.SetCompiler(builder())
	}

	if g.ListPath != "" {
		g.GenList = true
	}

	g.AppendIncPath(filepath.Join(appRoot, "share/ocala/include"))
	for _, i := range incPaths {
		if err := g.AppendIncPath(i); err != nil {
			fmt.Fprintln(g.ErrWriter, err.Error())
			return "", err
		}
	}

	for _, i := range g.Defs {
		if !reDefName.MatchString(i) {
			err := fmt.Errorf("invalid constant name `%s` from -D option", i)
			fmt.Fprintln(g.ErrWriter, err.Error())
			return "", err
		}
	}

	return flags.Args()[0], nil
}

func (cli *CLI) Run(args []string) int {
	g := &core.Generator{
		InReader:  cli.inReader,
		OutWriter: cli.outWriter,
		ErrWriter: cli.errWriter,
		DebugMode: os.Getenv("OCALADEBUG") == "1",
		Archs:     Archs,
		ListText:  &[]byte{},
	}
	if appRoot == "" {
		if err := cli.findAppRoot(); err != nil {
			fmt.Fprintln(g.ErrWriter, err.Error())
			return 1
		}
	}

	path, err := cli.ParseCommandLineOptions(g, args)
	if err == flag.ErrHelp {
		return 0
	} else if err != nil {
		return 1
	}

	if !g.CompileAndGenerate(path) {
		g.ErrWriter.Write(g.FullErrorMessage())
		return 1
	}

	return 0
}
