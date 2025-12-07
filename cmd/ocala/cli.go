package main

import (
	"flag"
	"fmt"
	"io"
	_ "ocala"
	"ocala/core"
	"os"
	"path/filepath"
	"regexp"
)

var version = "dev"
var appRoot = ""

var errNoInputFile = fmt.Errorf("input file required(-h for help)")
var errInvalidTarget = fmt.Errorf("invalid target arch(-h for help)")

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

var reDefName = regexp.MustCompile(`^[_A-Za-z][-_A-Za-z0-9]*$`)

func printError(g *core.Generator, err error) error {
	fmt.Fprintln(g.ErrWriter, err.Error())
	return err
}

func (*CLI) ParseCommandLineOptions(g *core.Generator, args []string) (string, error) {
	arch := ""
	incPaths := []string{}
	printVersion := false
	printHelp := false

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
	flags.BoolVar(&printHelp, "h", false, "Display this message")

	if err := flags.Parse(args[1:]); err != nil {
		return "", err
	} else if printHelp {
		flags.SetOutput(g.OutWriter)
		flags.Usage()
		return "", flag.ErrHelp
	} else if printVersion {
		fmt.Fprintf(g.OutWriter, "ocala %s\n", version)
		return "", flag.ErrHelp
	}
	if len(flags.Args()) == 0 {
		return "", printError(g, errNoInputFile)
	}

	if arch != "" {
		cc := core.NewCompiler(arch)
		if cc == nil {
			return "", printError(g, errInvalidTarget)
		}
		g.SetCompiler(cc)
	}

	if g.ListPath != "" {
		g.GenList = true
	}

	g.AppendIncPath(filepath.Join(appRoot, "share/ocala/include"))
	for _, i := range incPaths {
		if err := g.AppendIncPath(i); err != nil {
			return "", printError(g, err)
		}
	}

	for _, i := range g.Defs {
		if !reDefName.MatchString(i) {
			err := fmt.Errorf("invalid constant name `%s` from -D option", i)
			return "", printError(g, err)
		}
	}

	return flags.Args()[0], nil
}

func (cli *CLI) Run(args []string) int {
	core.Debug.Enabled = os.Getenv("OCALADEBUG") == "1"

	g := &core.Generator{
		InReader:  cli.inReader,
		OutWriter: cli.outWriter,
		ErrWriter: cli.errWriter,
		ListText:  &[]byte{},
	}
	if appRoot == "" {
		path, err := core.FindAppRoot(cli.executable)
		if err != nil {
			fmt.Fprintln(g.ErrWriter, err.Error())
			return 1
		}
		appRoot = path
	}

	path, err := cli.ParseCommandLineOptions(g, args)
	if err == flag.ErrHelp {
		return 0
	} else if err != nil {
		return 1
	}

	g.CompileAndGenerate(path)
	g.FlushMessages()
	if g.Err != nil {
		return 1
	}
	return 0
}
