# opinionated-actionlint
This is an extension to actionlint to add rules and configuration which is more opinionated than the default actionlint.

It is highly recommended that a workflow running [ossf/scorecard-action](https://github.com/ossf/scorecard-action) is set up for any repository that also uses this action, as that action detects additional security-related issues. This action does not detect anything the `ossf/scorecard-action` is capable of identifying.

## TODO
* Dependabot doesn't support updating Docker images referenced in `uses`. Convert to immutable GitHub action when it gets released - https://github.com/github/roadmap/issues/592

## Differences
### Only available as a container image
Simplifies distribution. Run locally by running `docker run -v $(pwd):/repo --workdir /repo ghcr.io/wjam/opinionated-actionlint:v0.1.0` or adding the following as GitHub workflow
```yaml
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@ # etc
      - uses: docker://ghcr.io/wjam/opinionated-actionlint:v0.1.0
```

### Unable to ignore issues
The only reason I came across for wanting to ignore issues was to format the output, which isn't required.

### Unable to format the output
Automatically produces messages that are [parsed by GitHub as error messages](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions#setting-an-error-message) if detected to be running under CI.

### Don't allow the use of actions/github-script action
While it makes it easy to create a new action, it's easy to introduce security problems as the script is passed *twice* - once when the GitHub context expressions are rendered and again when the script is executed.

The solution for this rule is to just create a 'real' action as it's not that difficult and provides benefits such as easier automated testing. Another solution is to use an existing open source action.

### Use environment variables to pass values from GitHub context expressions into run blocks
Shell scripts have a similar problem to `actions/github-script` in that they will be parsed twice - when expanding the GitHub context expressions and when the script is run. Small shell scripts are too useful to be completely avoided, so an outright ban is impossible.

While it is possible to spot GitHub context expressions that come directly from something the user controls (e.g. PR title), it's impossible to do this for indirect values (e.g. the output of an action which is the directory that was added in this PR).

The solution to this is to pass the values of GitHub context expressions into the script using environment variables.

### Don't use default workflow permissions
The default permissions given to a workflow is generally more than most workflow jobs will need.

To solve this, just define what permissions the workflow (or job) actually needs.

