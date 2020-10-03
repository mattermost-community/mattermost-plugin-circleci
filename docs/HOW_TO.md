# Use the Mattermost Plugin for CircleCI

## Table of Contents

-   [Audience](#audience)
-   [About the CircleCI Plugin](#about-the-circleci-plugin)
-   [Before You Start](#before-you-start)
-   [Configuration](#configuration)
-   [Onboarding Your Users](#onboarding-your-users)
-   [Slash Commands](#slash-commands)
-   [Frequently Asked Questions](#frequently-asked-questions)
-   [License](#license)
-   [Development](#development)

// TODO Add a screenshot here

## Audience

This guide is intended for Mattermost System Admins setting up the CircleCI plugin and Mattermost users who want information about the plugin functionality. For more information about contributing to this plugin, visit the [Development section](#development).

## About the CircleCI Plugin

The Mattermost CircleCI plugin uses a personal API token to connect your Mattermost account to CircleCI to interact with the API.

After your System Admin has [configured the CircleCI plugin](#configuration), run `/circleci account connect` in a Mattermost channel to connect your Mattermost and CircleCI accounts.

Once connected, you'll have access to the following features:

-   **Pipeline and workflow management** - Get informations about pipelines or workflows or triggering new ones
-   **Workflows notifications** - Receive worflows notifications directly in your Mattermost channel, including their status
-   **Slash commands** - Interact with the CircleCI plugin using the `/circleci` slash command. Read more about slash commands [here](#slash-commands).

## Before You Start

This guide assumes you have:

-   A project, hosted on Github or Bitbucket,
-   A CircleCI account, which can access to the project,
-   A Mattermost server running v5.12 or higher, with a configured [Site URL](https://docs.mattermost.com/administration/config-settings.html?highlight=site%20url#site-url).

## Configuration

### Step 1: Configure the Bot account in Mattermost

If you have an existing Mattermost user account with the name `circleci`, the plugin will post using the `circleci` account but without a `BOT` tag.

To prevent this, either:

-   Convert the `circleci` user to a bot account by running `mattermost user convert circleci --bot` in the CLI, or:
-   If the user is an existing user account you want to preserve, change its username and restart the Mattermost server. Once restarted, the plugin will create a bot account with the name `circleci`.

### Step 2: Configure the plugin in Mattermost

**Generate the keys**

Open **System Console > Plugins > CircleCI**:

-   Generate a new value for **Webhooks Secret**
-   Generate a new value for **At Rest Encryption Key**.
-   Hit **Save**.
-   Go to **System Console > Plugins > Management** and click **Enable** to enable the CircleCI plugin.

You're all set!

## Onboarding Your Users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their CircleCI account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team,
>
> We've set up the Mattermost CircleCI plugin, so you can get notifications from CircleCI in Mattermost. To get started, run the `/circleci account connect` slash command from any channel within Mattermost to connect your Mattermost account with CircleCI. Then, take a look at the [slash commands](#slash-commands) section for details about how to use the plugin.

## Slash Commands

### Subscribe to webhooks notifications

Subscribe a channel to notifications from a CircleCI project.

![Success notifications received in Mattermost](./successful-notification.jpg)

#### Steps

1.  In the channel you want to subscribe to notifications, type `/circleci subscription add <org-name> <project-name>`.
2.  Install the [Mattermost Plugin Notify Orb](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify) for CircleCI in your project. You usually do this by modifing the `.circleci/config.yml`.

    -   You can add the command [status](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#usage-status) in your existing jobs to get a notification when this job is finished
    -   Or you can setup the [approval-notification](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#jobs-approval-notification) job in a workflow to warn that you have a workflow waiting for approval.

3.  Add the webhook URL given by `/circleci subscription add` to your CircleCI project.

    -   You may add it to the orb as a parameter, but this is discouraged as it should be treated like a secret
    -   You should add it as a Environment Variable named `MM_WEBHOOK`, through the [CircleCI UI](https://circleci.com/docs/2.0/env-vars/#setting-an-environment-variable-in-a-project) or using the plugin: `/circleci project env add projectSlug MM_WEBHOOK <webhook-url>`

## Frequently Asked Questions

### Any information missing?

Feel free to fill a [Github Issue](https://github.com/nathanaelhoun/mattermost-plugin-circleci/issues/new/choose) and we'll add your information to this guide!

## License

Apache License.

## Development

This plugin only contains a server portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.
