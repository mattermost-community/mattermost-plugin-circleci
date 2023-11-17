# Mattermost CircleCI Plugin 

[![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-circleci/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-circleci)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-circleci)](https://github.com/mattermost/mattermost-plugin-circleci/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattermost/mattermost-plugin-circleci)](https://goreportcard.com/report/github.com/mattermost/mattermost-plugin-circleci)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-circleci/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-circleci/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)
[![Mattermost Community Channel](https://img.shields.io/badge/Mattermost%20Community-~Plugin%3A%20CircleCI-blue)](https://community.mattermost.com/core/channels/plugin-circleci)

**Help Wanted Tickets [here](https://github.com/mattermost/mattermost-plugin-circleci/issues)**

# Contents

- [Overview](#overview)
- [Features](#features)
- [Admin Guide](docs/admin-guide.md)
- [End User Guide](#end-user-guide)
- [Contribute](#contribute)
- [License](#license)
- [Security Vulnerability Disclosure](#security-vulnerability-disclosure)
- [Get Help](#get-help)

## Overview

The [CircleCI Orb for Mattermost Plugin](https://github.com/nathanaelhoun/circleci-orb-mattermost-plugin-notify) by [@nathanaelhoun](https://github.com/nathanaelhoun) interacts with jobs, builds, or workflows, and receives notifications in Mattermost channels. The Mattermost CircleCI plugin uses a personal API token to connect your Mattermost account to CircleCI to interact with the API.

### Thanks to

-   **[@jszwedko](https://github.com/jszwedko)** for his [CircleCI v1 Go API](https://github.com/jszwedko/go-circleci)
-   **[@TomTucka](https://github.com/TomTucka)** and **[@darkLord19](https://github.com/darkLord19)** for this [CircleCI v2 Go API](https://github.com/darkLord19/circleci-v2)
-   [Mattermost](https://mattermost.org) for providing a good software and maintaining a great community.

## Features

Use the Circle CI plugin for:

- **Pipeline and workflow management:** Get information about pipelines or workflows, or trigger new ones.
- **Workflows notifications:** Receive workflows notifications, including their status, directly in your Mattermost channel.
- **Slash commands:** Interact with the CircleCI plugin using the `/circleci` slash command. 
- **Metrics:** Get summary metrics for a project's workflows or for a project workflow's jobs. Metrics are refreshed daily, and thus may not include executions from the last 24 hours.
- **Manage environment variables:** Set CircleCI [environment variables](https://circleci.com/docs/2.0/env-vars/) directly from Mattermost.
- **Add notifications:** Receive notifications on [held workflows](https://circleci.com/docs/2.0/workflows/#holding-a-workflow-for-a-manual-approval) and approve them directly from your Mattermost channel. Learn how to set up a held workflow notification in [the Orb documentation](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#jobs-approval-notification).

For more information about contributing to this plugin, visit the Development section.



## [Admin Guide](docs/admin-guide.md)

## End User Guide

### Get Started

### Use the Plugin

### Slash commands

After your System Admin has configured the CircleCI plugin, run `/circleci account connect` in a Mattermost channel to connect your Mattermost and CircleCI accounts.
By default, the commands use the project set by `/circleci config`, unless a specific project is specified by the argument `--project <vcs/org-name/project-name>` (possible on all commands).

#### Connect to your CircleCI account

|                                |                                                   |
| ------------------------------ | ------------------------------------------------- |
| `/circleci account view`       | Get information about yourself.                   |
| `/circleci account connect`    | Connect your Mattermost account to CircleCI.      |
| `/circleci account disconnect` | Disconnect your Mattermost account from CircleCI. |

#### Set your default project

|                                                 |                                                                                     |
| ----------------------------------------------- | ----------------------------------------------------------------------------------- |
| `/circleci default`                             | View your currently configured default project.                                     |
| `/circleci default [vcs/org-name/project-name]` | Set new default project by passing value in the form `<vcs/org-name/project-name>`. |

#### Subscribe your channel to notifications

|                                           |                                                                    |
| ----------------------------------------- | ------------------------------------------------------------------ |
| `/circleci subscription list`             | List the CircleCI subscriptions for the current channel.           |
| `/circleci subscription add [--flags]`    | Subscribe the current channel to CircleCI notifications.           |
| `/circleci subscription remove [--flags]` | Unsubscribe the current channel to CircleCI notifications.         |
| `/circleci subscription list-channels`    | List all channels in the current team subscribed to notifications. |

#### Get insights about workflows and jobs

|                               |                                                                         |
| ----------------------------- | ----------------------------------------------------------------------- |
| `/circleci insight workflows` | Get summary of metrics for workflows over past 90 days.                 |
| `/circleci insight jobs`      | Get summary of metrics for jobs over past 90 days for a given workflow. |

#### Manage CircleCI projects

|                                        |                                 |
| -------------------------------------- | ------------------------------- |
| `/circleci project list-followed`      | List followed projects.         |
| `/circleci project recent-build`       | List the 10 last builds.        |
| `/circleci project env list`           | List environment variables.     |
| `/circleci project env add name value` | Add an environment variable.    |
| `/circleci project env remove name`    | Remove an environment variable. |

#### Manage pipelines

**Note:** Use the `all`, `recent`, and `mine` subcommands get the pipelineID for a particular pipeline.

|                                     |                                                                                                      |
| ----------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `/circleci pipeline trigger branch` | Trigger pipeline for a given branch.                                                                 |
| `/circleci pipeline trigger tag`    | Trigger pipeline for a given tag.                                                                    |
| `/circleci pipeline workflows`      | Get list of workflows for given pipeline.                                                            |
| `/circleci pipeline recent`         | Get list of all recently ran pipelines.                                                              |
| `/circleci pipeline all`            | Get list of all pipelines for a project.                                                             |
| `/circleci pipeline mine`           | Get list of all pipelines triggered by you for a project.                                            |
| `/circleci pipeline get`            | Get information about a single pipeline given pipeline ID.                                           |
| `/circleci pipeline get`            | Get information about a single pipeline for a given project by passing the unique pipelineID number. |

#### Manage workflows

**Note:** Use the `/circleci pipeline workflows` command to get the workflowID.

|                             |                                     |
| --------------------------- | ----------------------------------- |
| `/circleci workflow get`    | Get information about a workflow.   |
| `/circleci workflow jobs`   | Get jobs list for a given workflow. |
| `/circleci workflow rerun`  | Rerun a given workflow.             |
| `/circleci workflow cancel` | Cancel a given workflow.            |



### Frequently asked questions

#### How does the plugin save user data for each connected CircleCI user?

CircleCI user tokens are AES encrypted with an At Rest Encryption Key configured in the plugin's settings page. Once encrypted, the tokens are saved in the `PluginKeyValueStore` table in your Mattermost database.

#### How do I share feedback on this plugin?

Wanting to share feedback on this plugin?

Feel free to create a [GitHub Issue](https://github.com/mattermost/mattermost-plugin-circleci/issues) or join the [CircleCI Plugin channel](https://community.mattermost.com/core/channels/plugin-circleci) on the Mattermost Community server to discuss.

## Contribute

### I saw a bug, I have a feature request or a suggestion

Please fill a [GitHub issue](https://github.com/mattermost/mattermost-plugin-circleci/issues/new/choose), it will be very useful!

### Development

Pull Requests are welcome! You can contact us on the [Mattermost Community ~Plugin: CircleCI channel](https://community.mattermost.com/core/channels/plugin-circleci).

This plugin only contains a server portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploy with local mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. After configuring it, just run:

## Security vulnerability disclosure

Please report any security vulnerability to [https://mattermost.com/security-vulnerability-report/](https://mattermost.com/security-vulnerability-report/).

## Get Help

For questions, suggestions, and help, visit the  [Circleci Plugin channel](https://community.mattermost.com/core/channels/plugin-circleci) on our Community server.