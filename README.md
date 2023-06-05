[![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-circleci/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-circleci)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-circleci)](https://github.com/mattermost/mattermost-plugin-circleci/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattermost/mattermost-plugin-circleci)](https://goreportcard.com/report/github.com/mattermost/mattermost-plugin-circleci)
[![Mattermost Community Channel](https://img.shields.io/badge/Mattermost%20Community-~Plugin%3A%20CircleCI-blue)](https://community.mattermost.com/core/channels/plugin-circleci)

# Mattermost CircleCI Plugin 

The [CircleCI Orb for Mattermost Plugin](https://github.com/nathanaelhoun/circleci-orb-mattermost-plugin-notify) by [@nathanaelhoun](https://github.com/nathanaelhoun) interacts with jobs, builds, or workflows, and receives notifications in Mattermost channels. The Mattermost CircleCI plugin uses a personal API token to connect your Mattermost account to CircleCI to interact with the API.

Use the Circle CI plugin for:

- **Pipeline and workflow management:** Get information about pipelines or workflows, or trigger new ones.
- **Workflows notifications:** Receive workflows notifications, including their status, directly in your Mattermost channel.
- **Slash commands:** Interact with the CircleCI plugin using the `/circleci` slash command. 
- **Metrics:** Get summary metrics for a project's workflows or for a project workflow's jobs. Metrics are refreshed daily, and thus may not include executions from the last 24 hours.
- **Manage environment variables:** Set CircleCI [environment variables](https://circleci.com/docs/2.0/env-vars/) directly from Mattermost.
- **Add notifications:** Receive notifications on [held workflows](https://circleci.com/docs/2.0/workflows/#holding-a-workflow-for-a-manual-approval) and approve them directly from your Mattermost channel. Learn how to set up a held workflow notification in [the Orb documentation](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#jobs-approval-notification).

For more information about contributing to this plugin, visit the Development section.



## Thanks to

-   **[@jszwedko](https://github.com/jszwedko)** for his [CircleCI v1 Go API](https://github.com/jszwedko/go-circleci)
-   **[@TomTucka](https://github.com/TomTucka)** and **[@darkLord19](https://github.com/darkLord19)** for this [CircleCI v2 Go API](https://github.com/darkLord19/circleci-v2)
-   [Mattermost](https://mattermost.org) for providing a good software and maintaining a great community.

## Admin Guide

### Prerequisites 

This guide is intended for Mattermost System Admins setting up the CircleCI plugin and Mattermost users who want information about the plugin functionality.

This guide assumes you have:

* A project, hosted on github.com or bitbucket.org.
* A CircleCI SaaS account, which has access to the projects/org you want to interact with.
* Mattermost self-hosted: A Mattermost server running v5.12 or higher, with a configured Site URL. v5.24 or higher is recommended to have the autocomplete feature.

### Installation

#### Marketplace Installation

1. Go to **Main Menu > Plugin Marketplace** in Mattermost.
2. Search for "Circle CI" or find the plugin from the list.
3. Select **Install**.
4. When the plugin has downloaded and been installed, select **Configure**.

#### Manual Installation

If your server doesn't have access to the internet, you can download the latest [plugin binary release](https://github.com/mattermost/mattermost-plugin-circleci/releases) and upload it to your server via **System Console > Plugin Management**. The releases on this page are the same used by the Marketplace. To learn more about how to upload a plugin, see [the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

### Configuration

#### Step 1: Configure the bot account in Mattermost

If you have an existing Mattermost user account with the name circleci, the plugin will post using the circleci account but without a BOT tag.

To prevent this, either:

Convert the circleci user to a bot account by running `mattermost user convert circleci --bot` in the Mattermost CLI.

or

If the user is an existing user account you want to preserve, change its username and restart the Mattermost server. Once restarted, the plugin will create a bot account with the name `circleci`.

#### Step 2: Configure the plugin in Mattermost

To generate the keys needed, go to **System Console > Plugins > CircleCI**:

1. Generate a new value for **Webhooks Secret**. If the generated secret contains a forwardslash, please regenerate it.
2. Generate a new value for **At Rest Encryption Key**.
3. Select **Save**.
4. Go to **System Console > Plugins > Management** and choose **Enable** to enable the CircleCI plugin.

You're all set!

### Onboarding Your Users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their CircleCI account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team,
> 
> We've set up the Mattermost CircleCI plugin, so you can get notifications from CircleCI in Mattermost. 
> To get started, run the `/circleci account connect` slash command from any channel within Mattermost to connect your Mattermost account with CircleCI. 
> 
> Then, take a look at the slash commands section ([link](https://mattermost.gitbook.io/circle-ci-plugin/user-guide/slash-commands) or `/circleci help`) for details about how to use the plugin.

### Slash Commands

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

### Webhooks notifications

Subscribe a channel to notifications from a CircleCI project.

1. In the channel you want to subscribe to notifications, type `/circleci subscription add`.

  -   You can add the optional flag `--only-failed` to only receive notifications about failed jobs.
  -   You can temporarily use a project different that the one set with `/circleci default`, using the optional flag `--project <vcs/org-name/project-name>`.

2. Install the [Mattermost Plugin Notify Orb for CircleCI ](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify) in your project. You usually do this by modifying the `.circleci/config.yml`.

  -   You can add the command [status](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#usage-status) in your existing jobs to get a notification when this job is finished.
  -   Or you can set up the [approval-notification](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#jobs-approval-notification) job in a workflow to warn that you have a workflow waiting for approval.

3. Add the webhook URL given by `/circleci subscription` add to your CircleCI project.

  -   You may add it to the orb as a parameter, but this is discouraged as it should be treated like a secret.
  -   You should add it as a Environment Variable named `MM_WEBHOOK`, through the [CircleCI UI](https://circleci.com/docs/2.0/env-vars/#setting-an-environment-variable-in-a-project) or using the plugin: `/circleci project env add MM_WEBHOOK <webhook-url>`.

## Frequently Asked Questions

### How does the plugin save user data for each connected CircleCI user?

CircleCI user tokens are AES encrypted with an At Rest Encryption Key configured in the plugin's settings page. Once encrypted, the tokens are saved in the `PluginKeyValueStore` table in your Mattermost database.

### How do I share feedback on this plugin?

Wanting to share feedback on this plugin?

Feel free to create a [GitHub Issue](https://github.com/mattermost/mattermost-plugin-circleci/issues) or join the [CircleCI Plugin channel](https://community.mattermost.com/core/channels/plugin-circleci) on the Mattermost Community server to discuss.

## Contributing

### I saw a bug, I have a feature request or a suggestion

Please fill a [GitHub issue](https://github.com/mattermost/mattermost-plugin-circleci/issues/new/choose), it will be very useful!

## Development

Pull Requests are welcome! You can contact us on the [Mattermost Community ~Plugin: CircleCI channel](https://community.mattermost.com/core/channels/plugin-circleci).

This plugin only contains a server portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploy with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. After configuring it, just run:


## Security Vulnerability Disclosure

Please report any security vulnerability to [https://mattermost.com/security-vulnerability-report/](https://mattermost.com/security-vulnerability-report/).
