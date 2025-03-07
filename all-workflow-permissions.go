package main

import (
	"github.com/rhysd/actionlint"
)

func newBanAllWorkflowPermissions() actionlint.Rule {
	return &banAllWorkflowPermissions{
		actionlint.NewRuleBase("ban-all-workflow-permissions", "workflows should define what permissions they need"),
	}
}

type banAllWorkflowPermissions struct {
	actionlint.RuleBase
}

func (a *banAllWorkflowPermissions) VisitWorkflowPost(node *actionlint.Workflow) error {
	if node.Permissions != nil && node.Permissions.All != nil {
		a.Error(node.Permissions.Pos, "jobs should not use *-all workflow permissions")
	}

	for _, job := range node.Jobs {
		if job.Permissions != nil && job.Permissions.All != nil {
			a.Error(job.Permissions.Pos, "job should not use *-all workflow permissions")
		}
	}

	return nil
}
