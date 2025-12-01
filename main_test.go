package main

import (
	"embed"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*-output.txt
var outputFiles embed.FS

func TestRunLinter(t *testing.T) {
	tests := []struct {
		file   string
		output string
		err    error
		env    map[string]string
	}{
		{
			// Handled by ossf/scorecard-actino
			file:   "testdata/default-permissions.yaml",
			output: "testdata/default-permissions-output.txt",
		},
		{
			file:   "testdata/permissions-defined-on-job.yaml",
			output: "testdata/permissions-defined-on-job-output.txt",
		},
		{
			file:   "testdata/read-all-permission-job.yaml",
			output: "testdata/read-all-permission-job-output.txt",
			err: lintErrors{
				{
					Message:  "job should not use *-all workflow permissions",
					Filepath: "testdata/read-all-permission-job.yaml",
					Line:     6,
					Column:   5,
					Kind:     "ban-all-workflow-permissions",
				},
			},
		},
		{
			file:   "testdata/read-all-permission-workflow.yaml",
			output: "testdata/read-all-permission-workflow-output.txt",
			err: lintErrors{
				{
					Message:  "jobs should not use *-all workflow permissions",
					Filepath: "testdata/read-all-permission-workflow.yaml",
					Line:     4,
					Column:   1,
					Kind:     "ban-all-workflow-permissions",
				},
			},
		},
		{
			file:   "testdata/run-block-with-expression.yaml",
			output: "testdata/run-block-with-expression-output.txt",
			err: lintErrors{
				{
					Message:  "use of GitHub context expressions within run blocks is not allowed",
					Filepath: "testdata/run-block-with-expression.yaml",
					Line:     10,
					Column:   9,
					Kind:     "ban-run-block-with-github-context",
				},
			},
		},
		{
			file: "testdata/use-of-github-script.yaml",
			env: map[string]string{
				"CI": "true",
			},
			output: "testdata/use-of-github-script-output.txt",
			err: lintErrors{
				{
					Message:  "use of actions/github-script is not allowed",
					Filepath: "testdata/use-of-github-script.yaml",
					Line:     10,
					Column:   9,
					Kind:     "ban-github-script",
				},
			},
		},
		{
			file:   "testdata/checkout-no-with.yaml",
			output: "testdata/checkout-no-with-output.txt",
			err: lintErrors{
				{
					Message:  "persist-credentials should be set to 'false' when using actions/checkout",
					Filepath: "testdata/checkout-no-with.yaml",
					Line:     10,
					Column:   9,
					Kind:     "ban-checkout-with-persist-credentials",
				},
			},
		},
		{
			file:   "testdata/checkout-persist-credentials-false.yaml",
			output: "testdata/checkout-persist-credentials-false-output.txt",
		},
		{
			file:   "testdata/checkout-persist-credentials-true.yaml",
			output: "testdata/checkout-persist-credentials-true-output.txt",
			err: lintErrors{
				{
					Message:  "persist-credentials should be set to 'false' when using actions/checkout",
					Filepath: "testdata/checkout-persist-credentials-true.yaml",
					Line:     10,
					Column:   9,
					Kind:     "ban-checkout-with-persist-credentials",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			var buf strings.Builder
			actualErr := runLinter(&buf, func(s string) (string, bool) {
				v, ok := test.env[s]
				return v, ok
			}, "actionlint", test.file)

			output, err := outputFiles.ReadFile(test.output)
			require.NoError(t, err)

			assert.Equal(t, string(output), buf.String())

			if test.err == nil {
				assert.Nil(t, actualErr)
			} else {
				assert.Equal(t, test.err, actualErr)
			}
		})
	}
}

func TestLintErrors_Is(t *testing.T) {
	var err error = lintErrors{
		{},
	}
	v := errors.Is(err, lintErrors{})
	assert.True(t, v)
}
