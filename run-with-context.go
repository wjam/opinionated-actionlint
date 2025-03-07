package main

import (
	"github.com/rhysd/actionlint"
)

func newBanRunBlockWithGitHubContext() actionlint.Rule {
	return &banRunBlockWithGitHubContext{
		actionlint.NewRuleBase(
			"ban-run-block-with-github-context",
			"use environment variables to pass GitHub context values into run blocks",
		),
	}
}

type banRunBlockWithGitHubContext struct {
	actionlint.RuleBase
}

func (a *banRunBlockWithGitHubContext) VisitStep(node *actionlint.Step) error {
	run, ok := node.Exec.(*actionlint.ExecRun)
	if !ok {
		return nil
	}

	if actionlint.ContainsExpression(run.Run.Value) {
		a.Error(node.Pos, "use of GitHub context expressions within run blocks is not allowed")
	}

	return nil
}
