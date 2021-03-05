# Mattermost/CircleCI Plugin

A [Mattermost Plugin](https://developers.mattermost.com/extend/plugins/) for CircleCI, which uses the [CircleCI Orb for Mattermost Plugin](https://github.com/nathanaelhoun/circleci-orb-mattermost-plugin-notify) by [@nathanaelhoun](https://github.com/nathanaelhoun) to interact with jobs, builds, or workflows and receive notifications in Mattermost channels. The Mattermost CircleCI plugin uses a personal API token to connect your Mattermost account to CircleCI to interact with the API.

Use the Circle CI plugin for:

- **Pipeline and workflow management:** Get information about pipelines or workflows, or trigger new ones.
- **Workflows notifications:** Receive workflows notifications, including their status, directly in your Mattermost channel.
- **Slash commands:** Interact with the CircleCI plugin using the `/circleci` slash command. 
- **Metrics:** Get summary metrics for a project's workflows or for a project workflow's jobs. Metrics are refreshed daily, and thus may not include executions from the last 24 hours.
- **Manage environment variables:** Set CircleCI [environment variables](https://circleci.com/docs/2.0/env-vars/) directly from Mattermost.
- **Add notifications:** Receive notifications on [held workflows](https://circleci.com/docs/2.0/workflows/#holding-a-workflow-for-a-manual-approval) and approve them directly from your Mattermost channel. Learn how to set up a held workflow notification in [the Orb documentation](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#jobs-approval-notification).

For more information about contributing to this plugin, visit the Development section.
