# Webhooks Notificatons

Subscribe a channel to notifications from a CircleCI project.

1. In the channel you want to subscribe to notifications, type `/circleci subscription add`.
  - You can add the optional flag `--only-failed` to only receive notifications about failed jobs.
  - You can temporarily use a project different that the one set with `/circleci config`, using the optional flag `--project <vcs/org-name/project-name>`.

2. Install the [Mattermost Plugin Notify Orb for CircleCI ](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify) in your project. You usually do this by modifying the `.circleci/config.yml`.
  - You can add the command [status](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#usage-status) in your existing jobs to get a notification when this job is finished.
  - Or you can set up the [approval-notification](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify#jobs-approval-notification) job in a workflow to warn that you have a workflow waiting for approval.

3. Add the webhook URL given by `/circleci subscription` add to your CircleCI project.
  - You may add it to the orb as a parameter, but this is discouraged as it should be treated like a secret.
  - You should add it as a Environment Variable named `MM_WEBHOOK`, through the [CircleCI UI](https://circleci.com/docs/2.0/env-vars/#setting-an-environment-variable-in-a-project) or using the plugin: `/circleci project env add MM_WEBHOOK <webhook-url>`.
  
