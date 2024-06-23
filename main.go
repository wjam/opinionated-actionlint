package main

import (
	"flag"
	"io"
	"os"

	"github.com/rhysd/actionlint"
)

func main() {
	if err := runLinter(os.Stdout, os.Args...); err != nil {
		if _, ok := err.(lintErrors); ok {
			os.Exit(1)
		}
		panic(err)
	}
}

func runLinter(stdout io.Writer, args ...string) error {
	opts := &actionlint.LinterOptions{
		LogWriter: stdout,
		OnRulesCreated: func(rules []actionlint.Rule) []actionlint.Rule {
			return append(rules,
				newBanGitHubScriptAction(),
				newBanRunBlockWithGitHubContext(),
				newBanDefaultWorkflowPermissions(),
				newPinAction(),
			)
		},
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.SetOutput(stdout)
	flags.StringVar(&opts.Shellcheck, "shellcheck", "shellcheck", "")
	flags.StringVar(&opts.ConfigFile, "config", "", "")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	linter, err := actionlint.NewLinter(stdout, opts)
	if err != nil {
		return err
	}

	var lintErr []*actionlint.Error
	if files := flags.Args(); len(files) > 0 {
		lintErr, err = linter.LintFiles(files, nil)
	} else {
		lintErr, err = linter.LintRepository(".")
	}
	if err != nil {
		return err
	}

	if len(lintErr) != 0 {
		return lintErrors(lintErr)
	}

	return nil
}

var _ error = lintErrors{}

type lintErrors []*actionlint.Error

func (l lintErrors) Error() string {
	return "failed"
}
