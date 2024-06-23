package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunLinter(t *testing.T) {
	tests := []struct {
		file string
		err  error
	}{
		{
			file: "./testdata/default-permissions.yaml",
			err: lintErrors{
				{
					Message:  "job should not use default or *-all workflow permissions",
					Filepath: "./testdata/default-permissions.yaml",
					Line:     5,
					Column:   3,
					Kind:     "ban-default-workflow-permissions",
				},
			},
		},
		{
			file: "./testdata/permissions-defined-on-job.yaml",
		},
		{
			file: "./testdata/read-all-permission.yaml",
			err: lintErrors{
				{
					Message:  "job should not use default or *-all workflow permissions",
					Filepath: "./testdata/read-all-permission.yaml",
					Line:     7,
					Column:   3,
					Kind:     "ban-default-workflow-permissions",
				},
			},
		},
		{
			file: "./testdata/run-block-with-expression.yaml",
			err: lintErrors{
				{
					Message:  "use of GitHub context expressions within run blocks is not allowed",
					Filepath: "./testdata/run-block-with-expression.yaml",
					Line:     10,
					Column:   9,
					Kind:     "ban-run-block-with-github-context",
				},
			},
		},
		{
			file: "./testdata/unpinned-action.yaml",
			err: lintErrors{
				{
					Message:  "actions should be pinned",
					Filepath: "./testdata/unpinned-action.yaml",
					Line:     10,
					Column:   9,
					Kind:     "pin-actions",
				},
			},
		},
		{
			file: "./testdata/use-of-github-script.yaml",
			err: lintErrors{
				{
					Message:  "use of actions/github-script is not allowed",
					Filepath: "./testdata/use-of-github-script.yaml",
					Line:     10,
					Column:   9,
					Kind:     "ban-github-script",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			var buf strings.Builder
			err := runLinter(&buf, "actionlint", test.file)
			t.Log(buf.String())

			if test.err == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, test.err, err)
			}
		})
	}
}
