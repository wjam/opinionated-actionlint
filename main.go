package main

import (
	_ "embed"
	"errors"
	"flag"
	"io"
	"os"

	"github.com/rhysd/actionlint"
)

func main() {
	if err := runLinter(os.Stdout, os.LookupEnv, os.Args...); err != nil {
		if errors.Is(err, lintErrors{}) {
			os.Exit(1)
		}
		panic(err)
	}
}

func runLinter(stdout io.Writer, env func(string) (string, bool), args ...string) error {
	opts := &actionlint.LinterOptions{
		LogWriter: stdout,
		OnRulesCreated: func(rules []actionlint.Rule) []actionlint.Rule {
			return append(rules,
				newBanGitHubScriptAction(),
				newBanRunBlockWithGitHubContext(),
				newBanAllWorkflowPermissions(),
				newBanCheckoutWithPersistCredentials(),
			)
		},
	}

	if _, ok := env("CI"); ok {
		// https://github.com/actions/toolkit/issues/193#issuecomment-605394935
		opts.Format = githubActionTmpl
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.SetOutput(stdout)
	flags.StringVar(&opts.Shellcheck, "shellcheck", "shellcheck", "executable to use to run shellcheck")
	flags.StringVar(&opts.ConfigFile, "config-file", "", "actionlint config file location")

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

//go:embed githubAction.go.tmpl
var githubActionTmpl string

var _ error = lintErrors{}

type lintErrors []*actionlint.Error

func (l lintErrors) Error() string {
	return "failed"
}

func (l lintErrors) Is(err error) bool {
	_, ok := err.(lintErrors)
	return ok
}
