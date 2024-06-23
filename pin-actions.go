package main

import (
	"strings"

	"github.com/rhysd/actionlint"
)

func newPinAction() actionlint.Rule {
	return &pinAction{
		actionlint.NewRuleBase("pin-actions", "actions should pinned to a SHA"),
	}
}

type pinAction struct {
	actionlint.RuleBase
}

func (a *pinAction) VisitStep(node *actionlint.Step) error {
	action, ok := node.Exec.(*actionlint.ExecAction)
	if !ok {
		return nil
	}

	parts := strings.SplitN(action.Uses.Value, "@", 2)
	if len(parts) != 2 {
		return nil
	}

	if len(parts[1]) < 40 {
		a.Error(node.Pos, "actions should be pinned")
	}
	return nil
}
