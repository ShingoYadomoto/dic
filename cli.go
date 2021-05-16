package main

import (
	"flag"
	"fmt"
	"io"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		s bool
		e bool
		c bool

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&s, "s", false, "Match specification: STARTWITH (HEADWORD)")
	flags.BoolVar(&e, "e", false, "Match specification: ENDWITH   (HEADWORD)")
	flags.BoolVar(&c, "c", false, "Match specification: CONTAIN   (ANYWHERE)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	err := cli.validateFlags(s, e, c)
	if err != nil {
		fmt.Fprint(cli.errStream, err.Error())
		return ExitCodeError
	}

	translator := NewEdictJ2E(
		EdictJ2EMatchScope(s, e, c),
	)

	//buf, err := ioutil.ReadAll(os.Stdin)
	//if err != nil {
	//	fmt.Fprint(cli.errStream, err.Error())
	//	return ExitCodeError
	//}
	buf := []byte("日本")

	out, err := cli.callTranslate(translator, string(buf))
	if err != nil {
		fmt.Fprint(cli.errStream, err.Error())
		return ExitCodeError
	}

	fmt.Fprint(cli.outStream, out)

	return ExitCodeOK
}

func (cli *CLI) validateFlags(s, e, c bool) error {
	var trueCount int

	if s {
		trueCount++
	}
	if e {
		trueCount++
	}
	if c {
		trueCount++
	}

	if trueCount > 1 {
		return fmt.Errorf("match flag(-s, -e, -c) counts must be 0 or 1")
	}

	return nil
}

func (cli *CLI) callTranslate(t Translator, origin string) (string, error) {
	return t.Translate(origin)
}
