# opinionated-actionlint
This is an extension to actionlint to add rules and configuration which is more opinionated than the default actionlint.


## Don't allow the use of actions/github-script action
While it makes it easy to create a new action, it's easy to introduce security problems as the script is passed *twice* - once when the GitHub context expressions are rendered and again when the script is executed.

The solution for this rule is to just create a 'real' action as it's not that difficult and provides benefits such as easier automated testing. Another solution is to use an existing open source action.

## Pin actions using their Git SHA rather a Git tag or branch
Using a Git tag or branch as the action version is standard practice, but means it's possible for the action to change without you being aware of the change. These changes could be a security issue, or could just lead to wasted time when a workflow suddenly no longer works.

Solution for this is to 'pin' the action using a commit SHA, typically with the Git tag as a comment afterwords.

## Use environment variables to pass values from GitHub context expressions into run blocks
Shell scripts have a similar problem to `actions/github-script` in that they will be parsed twice - when expanding the GitHub context expressions and when the script is run. Small shell scripts are too useful to be completely avoided, so an outright ban is impossible.

While it is possible to spot GitHub context expressions that come directly from something the user controls (e.g. PR title), it's impossible to do this for indirect values (e.g. the output of an action which is the directory that was added in this PR).

The solution to this is to pass the values of GitHub context expressions into the script using environment variables.

## Don't use default workflow permissions
The default permissions given to a workflow is generally more than most workflow jobs will need.

To solve this, just define what permissions the workflow (or job) actually needs.

