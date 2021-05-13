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
		a bool
		s bool
		e bool
		c bool

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&a, "a", false, "")
	flags.BoolVar(&a, "a", false, "(Short)")
	flags.BoolVar(&s, "s", false, "")
	flags.BoolVar(&s, "s", false, "(Short)")
	flags.BoolVar(&e, "e", false, "")
	flags.BoolVar(&e, "e", false, "(Short)")
	flags.BoolVar(&c, "c", false, "")
	flags.BoolVar(&c, "c", false, "(Short)")

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

	_ = a

	_ = s

	_ = e

	_ = c

	return ExitCodeOK
}
