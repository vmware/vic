# Contributing to VIC Engine

The VIC project team welcomes contributions from the community. If you wish to
contribute code and you have not signed our contributor license agreement (CLA),
our bot will update the issue when you open a Pull Request.
For any questions about the CLA process, please refer to our
[FAQ][cla].

[cla]:https://cla.vmware.com/faq


## Community

To connect with the community, please [join][slack] our public Slack workspace,
which includes a [#vic-engine slack channel][slack-channel] for this project.

[slack]:https://code.vmware.com/join
[slack-channel]:https://vmwarecode.slack.com/messages/vic-engine


## Getting started

Go is the primary programming languge used by the project. These instructions
will assume that you are familiar with Go and have [installed][go-install] it.
If you are interested in contributing to ancillary portions of the project, such
as documentation, tests, or scripts, this may not be strictly necessary, however
it will provide consistency between your development environment and that of
other contributors, simplifying subsequent instructions.

To begin contributing, please create your own [fork][fork] of this repository.
This will allow you to share proposed changes with the community for review.

The [hub][hub] utility can be used to do this from the command line.

``` shell
go get github.com/vmware/vic
cd $(go env GOPATH)/src/github.com/vmware/vic
hub fork
```

See the [README](README.md#building) for build instructions.

[fork]:https://help.github.com/articles/fork-a-repo/
[installed]:https://golang.org/doc/install


## Contribution flow

A rough outline of a contributor's workflow might look like:

1. Create a topic branch from where you want to base your work
2. Make commits of logical units
3. Make sure your commit messages are in the proper format (see below)
4. Push your changes to a topic branch in your fork of the repository
5. [Test your changes](#automated-testing)
6. Submit a pull request

Example:

``` shell
git checkout -b my-new-topic-branch master
# (make changes)
git commit
git push -u YOUR_USER my-new-topic-branch
```

Note: You should push your topic branches to your fork, not origin, even if you
are an existing contributor with write access to the repository.

### Stay in sync with upstream

When your branch gets out of sync with the vmware/master branch, use the
following to update:

``` shell
git checkout my-new-topic-branch
git remote update
git pull --rebase origin master
git push --force-with-lease my-new-topic-branch
```

Note: In this case, we are able to invoke `git push` without specifying a remote
because we previously invoked `git push` with the `-u` (`--set-upstream` flag).

To make `--rebase` the default behavior when invoking `git pull`, you can use
`git config pull.rebase true`. This makes the history of your topic branch
easier to read by avoiding merge commits.

### Code style

Writing code that is easy to understand, debug, and modify is important to the
long-term health of a project. To help achieve these goals, we have adopted
coding style recommendations of the broader software community.

[Effective Go][effective-go] introduces and discusses many key ideas and the
[Code Review Comments wiki][crc] lists many common mistakes.

The [Robot Frameowrk User Guide][robot-user-guide] includes stylistic tips and
[How To Write Good Test Cases][good-test-cases] givs additional recommendations.

Docker provides [Best practices for writing Dockerfiles][dockerfiles].

For Bash, [unofficial Bash Strict Mode][bash-strict-mode] makes scripts behave
more reliable and [shellcheck][shellcheck] can catch many common mistakes.

[effective-go]:https://golang.org/doc/effective_go.html
[crc]:https://github.com/golang/go/wiki/CodeReviewComments
[robot-user-guide]:http://robotframework.org/robotframework/latest/RobotFrameworkUserGuide.html
[good-test-cases]:https://github.com/robotframework/HowToWriteGoodTestCases/blob/master/HowToWriteGoodTestCases.rst
[dockerfiles]:https://docs.docker.com/develop/develop-images/dockerfile_best-practices/
[bash-strict-mode]:http://redsymbol.net/articles/unofficial-bash-strict-mode/
[shellcheck]:https://github.com/koalaman/shellcheck

### Formatting commit messages

While the contents of your changes are easily improved in the future, your
commit message becomes part the permanent historical record for the repository.
Please take the time to craft meaningful commits with useful messages.

[How to Write a Git Commit Message][commitmsg] provides helpful conventions.

To be reminded when you may be making a common commit message mistake, you can
use the [git-good-commit][commithook] commit hook.

Example:
```shell
curl https://cdn.rawgit.com/tommarshall/git-good-commit/v0.6.1/hook.sh > .git/hooks/commit-msg && chmod +x .git/hooks/commit-msg
```

Please include any related GitHub issue references in the body of the pull
request, but not the commit message. See [GFM syntax][gfmsyntax] for referencing
issues and commits.

[commitmsg]:http://chris.beams.io/posts/git-commit/
[commithook]:https://github.com/tommarshall/git-good-commit
[gfmsyntax]:https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown

### Updating pull requests

If your PR fails to pass CI or needs changes based on code review, you'll want
to make additional commits to address these and push them to your topic branch
on your fork.

Providing updates this way instead of amending your existing commit makes it
easier for reviewers to see what has changed since they last looked at your
pull request.

You can use the `--fixup` and `--squash` options of `git commit` to communicate
your intent to combine these changes with a previous commit before merging.

Be sure to add a comment to the PR indicating your new changes are ready to
review, as GitHub does not generate a notification when you push to your topic
branch to update your pull request.

### Preparing to merge

After the review process is complete and you are ready to merge your changes,
you should rebase your changes into a series of meaningful, atomic commits.

If you have used the `--fixup` and `--squash` options suggested above, you can
leverage `git rebase -i --autosquash` to re-organize some of your history
automatically based on the intent you previously communicated.

If you have multiple commits on your topic branch, update the first line of
each commit's message to include your PR number. If you have a single commit,
you can use the "Squash & Merge" operation to do this automatically.

Once you've cleaned up the history on your topic branch, it's best practice to
wait for CI to run one last time before merging.

### Merging

Generally, we avoid merge commits on `master`. We suggest using "Squash & Merge"
if you are merging a single commit or "Rebase & Merge" if you are merging a
series of related commits. If you believe creating a merge commit is the right
operation for your change (e.g., because you're merging a long-lived feature
branch), please note that in your pull request.


## Automated Testing

Several kinds of automated testing are used by the project.

1. Compile-time checks are used to statically analyze the product code. These
   are run via `make check`.
    - `goimports`, `gofmt`, and `golint` are used for Go linting.
    - `gas` is used to check for potential security issues.
    - `missspell.sh` is used to check for common spelling mistakes.
    - `header-check.sh` is used to check for copyright headers.
    - `whitespace-check.sh` is used to check for trailing whitespace.
2. Compile-time testing is used to verify functionality of some individual
   components. These tests do not depend on external systems to function. These
   are run via `make test` (or `make focused-test` for pending changes).
    - Unit tests are used to verify the functionality of individual functions,
      methods, structs, and packages.
    - Simulated tests leverage [vcsim][vcsim] to verify behavior of packages
      which interact with vSphere.
3. [Integration tests](tests/README.md) are used to verify the behavior of the
   product against a live environment. These tests can be run against a VMware
   vSphere ESXi host or a VMware vCenter Server (tests only applicable to one
   environment or the other should automatically be skipped). All integration
   tests are automatically against vCenter Server when a change is merged to
   `master` or a `release/*` branch. A configurable set of integration tests are
   run when a pull request is submitted or updated. `local-integration-test.sh`
   can be used to run tests from a development environment.
    - Integration tests are written using the Robot Framework and divided into
      various Groups covering areas of product functionality and Suites within
      those Groups.
4. Scenario tests are used for more complex verification. These are periodically
   run using internal VMware infrastructure.
    - Interoperability tests are used to verify functionality against various
      supported versions of infrastructure components.
    - Workload tests mimic a realistic customer workload.

[vcsim]:https://github.com/vmware/govmomi/tree/master/vcsim

### Drone

[Drone][dronesrc] is used to run compile-time checks, compile-time tests, and
integration tests on pull requests and pushed commits.

By default, a pull request builds the project and runs compile-time tests,
compile-time checks, and the "regression" integration test group. To customize
the tests that run on a pull request, directives can be included in the body.
The [pull request template](PULL_REQUEST_TEMPLATE.md) documents the supported directives and their use.

Links to Drone builds results can be found within a pull request and on the list
of changes pushed to a branch. Results can also be browsed [directly][dronevic].

As the Drone environment is a shared resource, it is best to run tests locally
before publishing a pull request. If you don't have an environment suitable for
running the tests, you may leverage the Drone environment. When doing so, it is
best to include `[WIP]` (work in progress) at the beginning of the title of the
pull request to alert readers that the change is not ready for review.

Drone builds can be restarted via the web interface or the [CLI][dronecli]:
```shell
export DRONE_TOKEN=<Drone Token, from https://ci-vic.vmware.com/account/token>
export DRONE_SERVER=https://ci-vic.vmware.com

drone build start vmware/vic <Build Number>
```

For security reasons, your pull request build may fail if you are not a member
of the `vmware` organization in GitHub. If this occurs, leave a comment on your
pull request asking that the build be restarted by an organization member.

When an organization member restarts a build submitted by a user who is not an
organization member, they should include the `SKIP_CHECK_MEMBERSHIP` parameter:
```shell
drone build start --param SKIP_CHECK_MEMBERSHIP=true vmware/vic <Build Number>
```

[dronevic]:https://ci-vic.vmware.com/vmware/vic
[dronesrc]:https://github.com/drone/drone
[dronecli]:http://docs.drone.io/cli-installation/


## Reporting bugs and creating issues

Communicating clearly helps with efficient triage and resolution of reported
issues.

The summary of each issue will likely be read by many people. Quickly conveying
the essence of the problem you are experiencing helps get the right people
involved. Reports which are vague or unclear may take longer to be routed to
a domain expert.

The body of an issue should communicate what you are trying to accomplish and
why; understanding your goal allows others to suggest potential workarounds. It
should include specific details about what is (or isn't happening).

Proactively including screenshots and logs can be very helpful. When including
log snippets in the body of an issue or a comment instead of as an attachment,
please ensure that formatting is preserved by using [code blocks][code].
Consider formatting longer logs so that they are not shown by default.

Example:
```
<detail><summary>View Logs</summary>
<pre><code>
... (log content)
</code></pre>
</detail>
```

[code]:https://help.github.com/articles/creating-and-highlighting-code-blocks/


## Browsing and managing issues

We use [ZenHub][zenhub] for project management on top of GitHub issues. Boards
can be viewed [via the ZenHub website][board-zenhub] or, once you have installed
the [brower extension][zenhub-plugin], directly [within GitHub][board-github].
ZenHub integrates tightly with GitHub to provide additional project management
[functionality][zenhub-features] intended to facilitate an iterative development
approach, such as support for Epics and [related concepts][zenhub-agile].

[zenhub]:https://www.zenhub.io/
[zenhub-plugin]:https://www.zenhub.com/extension
[zenhub-features]:https://help.zenhub.com/support/solutions/articles/43000010337-take-a-tour-of-zenhub-s-key-features
[zenhub-agile]:https://help.zenhub.com/support/solutions/articles/43000010338-agile-concepts-in-github-and-zenhub
[board-zenhub]:https://app.zenhub.com/workspace/o/vmware/vic
[board-github]:https://github.com/vmware/vic/issues#boards

### Pipelines

ZenHub organizes issues into pipelines, which represent the status of an issue.
We use the following pipelines:

1. **New Issues**: The default pipeline, for issues which need to be reviewed by
   a member of the team.
2. **Not Ready**: Issues that have undergone initial review, but require
   additional refinement or information before they could be worked on.
3. **Backlog**: Issues that include all necessary information for work to begin.
4. **To Do**: Issues that will be worked on next.
5. **In Progress**: Issues that are being worked on.
6. **Verify**: Issues that are being reviewed.
7. **Closed**: Issues that are closed.

#### New Issues

The New Issues are triaged by the team at least once a week.  We try to keep issues from staying in this pipeline for
too long.  After triaging and issue, it will likely be moved to the backlog or stay under [Not Ready](#not-ready) for deferred
discussion.

For VIC engineers, you should set the priority based on the below guidelines. Everyone else, do not set the priority of a new issue.

##### Priorities

| Priority | Bugs | Features | Non Bugs |
| -------- | ---- | -------- | -------- |
| priority/p0 | Bugs that NEED to be fixed immediately as they either block meaningful testing or are release stoppers for the current release. | No Feature should be p0. | An issue that is not a bug and is blocking meaningful testing. eg. builds are failing because the syslog server is out of space. |
| priority/p1 | Bugs that NEED to be fixed by the assigned phase of the current release. | A feature that is required for the next release, typically an anchor feature; a large feature that is the focus for the release and drives the release date. | An issue that must be fixed for the next release. eg. Track build success rates. |
| priority/p2 | Bugs that SHOULD be fixed by the assigned phase of the current release, time permitting. | A feature that is desired for the next release, typically a pebble; a feature that has been approved for inclusion but is not considered the anchor feature or is considered good to have for the anchor feature. | An issue that we should fix in the next release. eg. A typo in the UI. |
| priority/p3 | Bugs that SHOULD be fixed by a given release, time permitting. | A feature that can be fixed in the next release. eg. Migrate to a new kernel version. Or a feature that is nice to have for a pebble. | An issue that can be fixed in the next release. eg. Low hanging productivity improvements. |
| priority/p4 | Bugs that SHOULD be fixed in a future (to be determined) release. | An issue or feature that will be fixed in a future release. | An issue or feature that will be fixed in a future release. |

#### Not Ready

The Not Ready column is for issues that need more discussion, details and/or triaging before being put in the [Backlog](#backlog). Issues in Not Ready should have assignee(s) to track whose input is needed to put the issue in the Backlog. For issues reported by VIC engineers: if the issue's details aren't fleshed out, the reporter should set themselves as the assignee.

#### Backlog

Issues in Backlog should be ready to be worked on in future sprints. For example, they may be feature requests or ideas for a future version of
the project. When moving issues to the Backlog, add more information (like requirements and outlines) into each issue. It's useful to
get ideas out of your head, even if you will not be touching them for a while.

To move an issue into the Backlog swim lane, it must have:

1. a `priority/...` label
2. a `team/...` label
3. an estimated level of effort (see [Story point estimates](#story-point-estimates) for guidance for mapping effort to story points)
4. no assignee (assignees are set when the issue is selected to work on)

Other labels should be added as needed.

Prioritize issues by dragging and dropping their placement in the pipeline. Issues higher in the pipeline are higher
priority; accordingly, they should contain all the information necessary to get started when the time
comes. Low-priority issues should still contain at least a short description.

#### To Do

This is the team's current focus and the issues should be well-defined. This pipeline should contain the high-priority
items for the current milestone. These issues must have an assignee, milestone, estimate and tags. Items are moved
from this pipeline to In Progress when work has been started.

To move an issue into the To Do swim lane, the assignee and milestone fields should be set.

#### In Progress

This is the answer to, "What are you working on right now?" Ideally, this pipeline will not contain more issues than
members of the team; each team member should be working on one thing at a time.

Issues in the In Progress swim lane must have an assignee.

After an issue is In Progress, it is best practice to update the issue with current progress and any discussions that may occur via the various collaboration tools used. An issue that is in progress should not go more than 2 days without updates.

Note: Epics should never be In Progress.

#### Verify

A "Verify" issue normally means the feature or fix is in code review and/or awaiting further testing. These issues require one final QE sign off or at the end of a sprint another dev that didn't work on the issue can verify the issue.

In most cases, an issue should be in Verify _before_ the corresponding PR is merged. The developer can then close the issue while merging the PR.

#### Closed

This pipeline includes all closed issues. It can be filtered like the rest of the Board – by Label, Assignee or Milestone.

This pipeline is also interactive: dragging issues into this pipeline will close them, while dragging them out will re-open them.

### Story point estimates

* Use the fibonacci pattern
* All bugs are a 2 unless we know it is significantly more or less work than the average bug
* 1 is easier than the average bug
* 3 is slightly more work than the average bug and probably should be about an average feature work for an easy feature (which includes design doc, implementation, testing, review)
* 5 is about 2x more work than the average bug and the highest single issue value we want
* Issues with an estimate higher than 5 should be decomposed further
* Unless otherwise necessary, estimates for EPICs are the sum of their sub-issues' estimates - EPICs aren't assigned an estimate themselves

### High level project planning

We use the following structure for higher level project management:
* Epic (zenhub) - implements a functional change - for example 'attach, stdout only', may span milestones and releases. Expected to be broken down from larger Epics into smaller epics prior to commencement.
* Milestones - essentially higher level user stories
* Labels - either by functional area (`component/...`) or feature (`feature/...`)

## Repository structure

The layout within the repository is as follows:
* `cmd` - the main packages for compiled components
* `doc` - all project documentation other than the standard files in the root
* `infra` - supporting scripts, utilities, et al.
* `isos` - ISO mastering scripts and uncompiled content
* `lib` - common library packages that are tightly coupled to vmware/vic
* `pkg` - packages that are not tightly coupled to vmware/vic and could be usefully consumed in other projects. There is still some sanitization to do here.
* `tests` - integration and scenario test code that doesn't use go test
* `vendor` - standard Go vendor model

## Troubleshooting

* If you're building the project in a VM, ensure that it has at least 4GB memory to avoid memory issues during a build.

