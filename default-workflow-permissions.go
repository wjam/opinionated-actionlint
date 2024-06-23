package main

import (
	"github.com/rhysd/actionlint"
)

func newBanDefaultWorkflowPermissions() actionlint.Rule {
	return &banDefaultWorkflowPermissions{
		actionlint.NewRuleBase("ban-default-workflow-permissions", "workflows should define what permissions they need"),
	}
}

type banDefaultWorkflowPermissions struct {
	actionlint.RuleBase
}

func (a *banDefaultWorkflowPermissions) VisitWorkflowPost(node *actionlint.Workflow) error {
	if node.Permissions != nil && node.Permissions.All == nil {
		return nil
	}

	for _, job := range node.Jobs {
		pos := job.Pos
		if job.Permissions != nil {
			pos = job.Permissions.Pos
		}
		if job.Permissions == nil || job.Permissions.All != nil {
			a.Error(pos, "job should not use default or *-all workflow permissions")
		}
	}

	return nil
}
