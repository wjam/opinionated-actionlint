package main

import (
	"strings"

	"github.com/rhysd/actionlint"
)

func newBanCheckoutWithPersistCredentials() actionlint.Rule {
	return &banCheckoutWithPersistCredentials{
		actionlint.NewRuleBase(
			"ban-checkout-with-persist-credentials",
			"persist-credentials should be set to 'false'",
		),
	}
}

type banCheckoutWithPersistCredentials struct {
	actionlint.RuleBase
}

func (a *banCheckoutWithPersistCredentials) VisitStep(node *actionlint.Step) error {
	action, ok := node.Exec.(*actionlint.ExecAction)
	if !ok {
		return nil
	}
	if !strings.HasPrefix(action.Uses.Value, "actions/checkout@") {
		return nil
	}

	if v, ok := action.Inputs["persist-credentials"]; ok && v.Value.Value == "false" {
		return nil
	}

	a.Error(node.Pos, "persist-credentials should be set to 'false' when using actions/checkout")
	return nil
}
