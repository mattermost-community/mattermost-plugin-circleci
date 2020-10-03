# Mattermost Plugin CircleCI [![CircleCI branch](https://img.shields.io/circleci/project/github/nathanaelhoun/mattermost-plugin-circleci/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-circleci)

A [Mattermost Plugin](https://developers.mattermost.com/extend/plugins/) for [CircleCI](https://circleci.com) to interact with jobs, builds or workflows and receive notifications in Mattermost channels.

This plugin uses the CircleCI Orb for Mattermost Plugin by **[@nathanaelhoun](https://github.com/nathanaelhoun)**: [check it out here](https://github.com/nathanaelhoun/circleci-orb-mattermost-plugin-notify).

## Features

#### Connect to your CircleCI account

-   `/circleci account view` - Get informations about yourself
-   `/circleci account connect <API token>` - Connect your Mattermost account to CircleCI
-   `/circleci account disconnect` - Disconnect your Mattermost account from CircleCI

#### Manage CircleCI projects

-   `/circleci project list-followed` - List followed projects
-   `/circleci project recent-build <owner-name> <project-name> <branch>` - List the 10 last builds for a project

#### Subscribe to notifications projects

-   `/circleci subscription list` — List the CircleCI subscriptions for the current channel
-   `/circleci subscription add <owner> <repository> [--flags]` — Subscribe the current channel to CircleCI notifications for a repository
-   `/circleci subscription remove <owner> <repository> [--flags]` — Unsubscribe the current channel to CircleCI notifications for a repository
-   `/circleci subscription list-channels <owner> <repository>` — List all channels subscribed to this repository in the current team

#### Config

-   `/circleci config <vcs/org-name/project-name>` - Allows you to set a default project to run your commands against

#### Pipeline

-   `/circleci pipeline recent vcs/orgname` - Lists recently built pipelines in a given org
-   `/circleci pipeline all vcs/orgname/project` - Lists all pipelines for a given project
-   `/circleci pipeline mine vcs/orgname/project` - Lists all pipelines triggered by you in a given project
-   `/circleci pipeline workflows pipelineId` - Lists all workflows for a given pipeline

## Installation instructions

_Coming_

## How to use this plugin

See [`HOW_TO.md`](./docs/HOW_TO.md)

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md)

## License

Apache License.

## Thanks to

-   **[@jszwedko](https://github.com/jszwedko)** for his [CircleCI v1 Go API](https://github.com/jszwedko/go-circleci)
-   **[@darkLord19](https://github.com/darkLord19)** for this [CircleCI v2 Go API](https://github.com/darkLord19/circleci-v2)
-   [Mattermost](https://mattermost.org) for providing a good software and having a great community
