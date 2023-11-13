## Admin Guide

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Webhooks](#web-hooks)
- [Slash Commands](../README.md/#slash-commands)
- [Onboard Users](#onboard-users)
- [FAQ](#faq)
- [Get Help](#get-help)

## Prerequisites 

This guide is intended for Mattermost System Admins setting up the CircleCI plugin and Mattermost users who want information about the plugin functionality.

This guide assumes you have:

* A project, hosted on github.com or bitbucket.org.
* A CircleCI SaaS account, which has access to the projects/org you want to interact with.
* Mattermost self-hosted: A Mattermost server running v5.12 or higher, with a configured Site URL. v5.24 or higher is recommended to have the autocomplete feature.

## Installation

### Manual installation

You can download the latest [plugin binary release](https://github.com/mattermost/mattermost-plugin-circleci/releases) and upload it to your server via **System Console > Plugin Management**.

## Configuration

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

## Web Hooks

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

## [Slash Commands](../README.md/#slash-commands)

## Onboard Users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their CircleCI account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team,
> 
> We've set up the Mattermost CircleCI plugin, so you can get notifications from CircleCI in Mattermost. 
> To get started, run the `/circleci account connect` slash command from any channel within Mattermost to connect your Mattermost account with CircleCI. 
> 
> Then, take a look at the slash commands section below or `/circleci help` for details about how to use the plugin.