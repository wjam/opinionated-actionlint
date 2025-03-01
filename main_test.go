package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunLinter(t *testing.T) {
	tests := []struct {
		file   string
		output string
		err    error
		env    map[string]string
	}{
		{
			file: "./testdata/default-permissions.yaml",
			output: `./testdata/default-permissions.yaml:5:3: job should not use default or *-all workflow permissions [ban-default-workflow-permissions]
  |
5 |   ci:
  |   ^~~
`,
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
			output: `./testdata/read-all-permission.yaml:7:3: job should not use default or *-all workflow permissions [ban-default-workflow-permissions]
  |
7 |   ci:
  |   ^~~
`,
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
			output: `./testdata/run-block-with-expression.yaml:10:9: use of GitHub context expressions within run blocks is not allowed [ban-run-block-with-github-context]
   |
10 |       - run: |
   |         ^~~~
`,
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
			file: "./testdata/use-of-github-script.yaml",
			env: map[string]string{
				"CI": "true",
			},
			output: "::error file=./testdata/use-of-github-script.yaml,line=10,col=9::use of actions/github-script is not allowed%0A```%0A      - uses: actions/github-script@v7.0.1%0A        ^~~~~%0A```\n",
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
			err := runLinter(&buf, func(s string) (string, bool) {
				v, ok := test.env[s]
				return v, ok
			}, "actionlint", test.file)

			assert.Equal(t, test.output, buf.String())

			if test.err == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, test.err, err)
			}
		})
	}
}
