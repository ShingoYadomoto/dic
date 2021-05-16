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

func color(t string, code int) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, t)
}

func red(t string) string {
	return color(t, 31)
}

func cyan(t string) string {
	return color(t, 36)
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		s bool
		e bool
		c bool

		p int

		w string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&s, "s", false, "Match specification: STARTWITH (HEADWORD)")
	flags.BoolVar(&e, "e", false, "Match specification: ENDWITH   (HEADWORD)")
	flags.BoolVar(&c, "c", false, "Match specification: CONTAIN   (ANYWHERE)")
	flags.IntVar(&p, "p", 0, "Max result count(default: 3)")
	flags.StringVar(&w, "w", "", "Word translated by")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	err := flags.Parse(args[1:])
	if err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	// Validate Flags
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
		fmt.Fprint(cli.errStream, "match flag(-s, -e, -c) counts must be 0 or 1")
		return ExitCodeError
	}

	// Run
	err = cli.run(s, e, c, p, w)
	if err != nil {
		fmt.Fprint(cli.errStream, err.Error())
		return ExitCodeError
	}

	return ExitCodeOK
}

func (cli *CLI) run(s, e, c bool, p int, w string) error {
	// New Translator
	var translator Translator

	opts := []EdictJ2EOption{
		EdictJ2EMatchScope(s, e, c),
	}
	if p != 0 {
		opts = append(opts, EdictJ2EPageSize(p))
	}
	translator = NewEdictJ2E(opts...)

	// Call Translate
	resultList, err := translator.Translate(w)
	if err != nil {
		return err
	}

	// Echo Results
	for _, result := range resultList {
		fmt.Fprintf(cli.outStream, "%s: %s\n", result.Origin, cyan(result.Dist))
	}

	return nil
}
