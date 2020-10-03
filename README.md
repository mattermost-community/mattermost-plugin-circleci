# Mattermost Plugin CircleCI

[![CircleCI branch](https://img.shields.io/circleci/project/github/nathanaelhoun/mattermost-plugin-circleci/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-circleci)
[![Release](https://img.shields.io/github/v/release/nathanaelhoun/mattermost-plugin-circleci)](https://github.com/nathanaelhoun/mattermost-plugin-circleci/releases/latest)

A [Mattermost Plugin](https://developers.mattermost.com/extend/plugins/) for [CircleCI](https://circleci.com) to interact with jobs, builds or workflows and receive notifications in Mattermost channels.

This plugin uses the CircleCI Orb for Mattermost Plugin by **[@nathanaelhoun](https://github.com/nathanaelhoun)**: [check it out here](https://github.com/nathanaelhoun/circleci-orb-mattermost-plugin-notify).

## Features

### Connect to your CircleCI account

-   `/circleci account view` — Get informations about yourself
-   `/circleci account connect` <API token> — Connect your Mattermost account to CircleCI
-   `/circleci account disconnect` — Disconnect your Mattermost account from CircleCI

### Set your default project

-   `/circleci config <vcs/org-name/project-name>` — View the config. Pass in the project (vcs/org/projectname) to set the default config

### Subscribe your channel to notifications

-   `/circleci subscription list` — List the CircleCI subscriptions for the current channel
-   `/circleci subscription add [--flags]` — Subscribe the current channel to CircleCI notifications for a project
-   `/circleci subscription remove [--flags]` — Unsubscribe the current channel to CircleCI notifications for a project
-   `/circleci subscription list-channels` — List all channels in the current team subscribed to a project

### Manage pipelines

-   `/circleci pipeline trigger <branch>` — Trigger pipeline for a project
-   `/circleci pipeline workflows <pipelineID>` — Get list of workflows for given pipeline
-   `/circleci pipeline recent <vcs-slug/org-name>` — Get list of all recently ran pipelines
-   `/circleci pipeline all` — Get list of all pipelines for a project
-   `/circleci pipeline mine` — Get list of all pipelines triggered by you for a project
-   `/circleci pipeline get` <pipelineID> — Get informations about a single pipeline

### Manage worflows

-   `/circleci workflow get <workflowID>` — Get informations about workflow
-   `/circleci workflow jobs <workflowID>` — Get jobs list of workflow
-   `/circleci workflow rerun <workflowID>` — Rerun a workflow
-   `/circleci workflow cancel <workflowID>` — Cancel a workflow

### Manage CircleCI projects

-   `/circleci project list-followed` — List followed projects
-   `/circleci project recent-build <branch>` — List the 10 last builds for a project
-   `/circleci project env <list|add|add>` — get, add or remove environment variables for given project

## Installation instructions

_Coming_

After installation, generate the secrets in **System Console > Plugins > CircleCI** for `Webhooks Secret` and `At Rest Encryption Key`.

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
