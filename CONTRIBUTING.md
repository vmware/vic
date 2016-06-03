# Contributing to vIC

## Community

Contributors and users are encouraged to collaborate using the following resources in addition to the GitHub issue
tracker:

- [Socialcast](https://vic-vmware.socialcast.com)

- [Gitter](https://gitter.im/vmware/vic)

## Getting started

First, fork the repository on GitHub to your personal account.

Note that _GOPATH_ can be any directory, the example below uses _$HOME/vic_.
Change _$USER_ below to your github username if they are not the same.

``` shell
export GOPATH=$HOME/vic
mkdir -p $GOPATH/src/github.com/vmware
go get github.com/vmware/vic
cd $GOPATH/src/github.com/vmware/vic
git config push.default nothing # anything to avoid pushing to vmware/vic by default
git remote rename origin vmware
git remote add $USER git@github.com:$USER/vic.git
git fetch $USER
```

See the [README](README.md) for build instructions.

## Contribution flow

This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work.
- Make commits of logical units.
- Make sure your commit messages are in the proper format (see below).
- Push your changes to a topic branch in your fork of the repository.
- Submit a pull request to vmware/vic.
- Your PR must receive at least one LGTM from a maintainer before merging.

Example:

``` shell
git checkout -b my-new-feature vmware/master
git commit -a
git push $USER my-new-feature
```

### Stay in sync with upstream

When your branch gets out of sync with the vmware/master branch, use the following to update:

``` shell
git checkout my-new-feature
git fetch -a
git rebase vmware/master
git push --force-with-lease $USER my-new-feature
```

### Updating pull requests

If your PR fails to pass CI or needs changes based on code review, you'll most likely want to squash these changes into
existing commits.

If your pull request contains a single commit or your changes are related to the most recent commit, you can simply
amend the commit.

``` shell
git add .
git commit --amend
git push --force-with-lease $USER my-new-feature
```

If you need to squash changes into an earlier commit, you can use:

``` shell
git add .
git commit --fixup <commit>
git rebase -i --autosquash vmware/master
git push --force-with-lease $USER my-new-feature
```

Be sure to add a comment to the PR indicating your new changes are ready to review, as github does not generate a
notification when you git push.

### Code style

The coding style suggested by the Golang community is used in vIC. See the
[style doc](https://github.com/golang/go/wiki/CodeReviewComments) for details.

Try to limit column width to 120 characters for both code and markdown documents such as this one.

### Format of the Commit Message

We follow the conventions on [How to Write a Git Commit Message](http://chris.beams.io/posts/git-commit/).

Be sure to include any related GitHub issue references in the commit message.  See
[GFM syntax](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown) for referencing issues
and commits.

## Reporting Bugs and Creating Issues

When opening a new issue, try to roughly follow the commit message format conventions above.

We use [Zenhub](https://www.zenhub.io/) for project management on top of GitHub issues.  Once you have the Zenhub
browser plugin installed, click on the [Boards](https://github.com/vmware/vic/issues#boards) tab to open the Zenhub task
board.

Our task board practices are as follows:

### New Issues

The New Issues are triaged by the team at least once a week.  We try to keep issues from staying in this pipeline for
too long.  After triaging and issue, it will likely be moved to the backlog or stay under New Issues for deferred
discussion.

### Backlog

Issues in Backlog are not a current focus. For example, they may be feature requests or ideas for a future version of
your project.

When moving issues to the Backlog, add more information (like requirements and outlines) into each issue. It’s useful to
get ideas out of your head, even if you will not be touching them for a while.

Prioritize issues by dragging and dropping their placement in the pipeline. Issues higher in the pipeline are higher
priority; accordingly, they should contain all the information necessary to get started when the time
comes.  Low-priority issues should still contain at least a short description.

### To Do

This is the team’s current focus and issues should be well-defined.  This pipeline should contain the high-priority
items for the current milestone.  These issues must have an assignee, milestone, estimate and tags.  Items are moved
from this pipeline to In Progress when work has been started.

### In Progress

This is the answer to, "What are you working on right now? Ideally, this pipeline will not contain more issues than
members of the team; each team member should be working on one thing at a time.

This pipeline is a good candidate for WIP (work-in-progress) limits. WIP limits help ensure your work flows smoothly,
and help bring to light any blockers or bottlenecks. Adjust WIP limits according to the size of your team.

To move an issue into the In Progress swim lane several steps must be taken.

1. Assign yourself as the owner.
2. Ensure the milestone is set (if there is one) and also review the labels to ensure they accurately reflect the issue.
3. Assign an estimated level of effort. See the below table for guidance for effort mapping.

After an issue is In Progress it is the best practice to update the issue with current progress and any discussions that may occur via the various collaboration tools used. An issue that is in progress should not go more than 2 days without updates.

Story Points | Story Size
------------ | -------------------------------------------------------
1            | Less than 1 day of effort
2            | 2 - 3 days of effort
3            | 3 - 5 days of effort
5            | 5 - 10 days of effort, consider splitting this if it's 7 - 10 days
8            | > 10 days, anything of this size should be split before moving into In Progress

Note: Epics should never be In Progress

### Done

A "Done" issue normally means the feature or fix is in code review and/or awaiting further testing.  We review these
issues as a team before closing them.

### Closed

This pipeline includes all closed issues, it can be filtered like the rest of the Board – by Label, Assignee or
Milestone.

This pipeline is also interactive: dragging issues into this pipeline will close them, while dragging them out will
re-open them.

## High level project planning

We use the following structure for higher level project management
* Epic (zenhub) - implements a functional change - for example 'attach, stdout only', may span milestones and releases. Expected to be broken down from larger Epics into smaller epics prior to commencement.
* Milestones - essentially higher level user stories
* Labels - either by functional area (`component/...`) or feature (`feature/...`)


## Repository structure

The layout in the repo is as follows - this is a recent reorganisation so there is still some mixing between directories:
* cmd - the main packages for compiled components
* doc - all project documentation other than the standard files in the root
* infra - supporting scripts, utilities, et al
* isos - ISO mastering scripts and uncompiled content
* lib - common library packages that are tightly coupled to vmware/vic
* pkg - packages that are not tightly coupled to vmware/vic and could be usefully consumed in other projects. There is still some sanitization to do here.
* tests - integration and system test code that doesn't use go test
* vendor - standard Go vendor model


