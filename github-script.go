package main

import (
	"strings"

	"github.com/rhysd/actionlint"
)

func newBanGitHubScriptAction() actionlint.Rule {
	return &banGitHubScriptAction{
		actionlint.NewRuleBase("ban-github-script", "ban use of action/github-script action"),
	}
}

type banGitHubScriptAction struct {
	actionlint.RuleBase
}

func (a *banGitHubScriptAction) VisitStep(node *actionlint.Step) error {
	action, ok := node.Exec.(*actionlint.ExecAction)
	if !ok {
		return nil
	}
	if strings.HasPrefix(action.Uses.Value, "actions/github-script@") {
		a.Error(node.Pos, "use of actions/github-script is not allowed")
	}
	return nil
}
