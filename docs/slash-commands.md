# Slash Commands

After your System Admin has configured the CircleCI plugin, run `/circleci account connect` in a Mattermost channel to connect your Mattermost and CircleCI accounts.
By default, the commands use the project set by `/circleci config`, unless a specific project is specified by the argument `--project <vcs/org-name/project-name>` (possible on all commands).

## Connect to your CircleCI account

|                                |                                                   |
| ------------------------------ | ------------------------------------------------- |
| `/circleci account view`       | Get information about yourself.                   |
| `/circleci account connect`    | Connect your Mattermost account to CircleCI.      |
| `/circleci account disconnect` | Disconnect your Mattermost account from CircleCI. |

## Set your default project

|                                                 |                                                                                     |
| ----------------------------------------------- | ----------------------------------------------------------------------------------- |
| `/circleci default`                             | View your currently configured default project.                                     |
| `/circleci default [vcs/org-name/project-name]` | Set new default project by passing value in the form `<vcs/org-name/project-name>`. |

## Subscribe your channel to notifications

|                                           |                                                                    |
| ----------------------------------------- | ------------------------------------------------------------------ |
| `/circleci subscription list`             | List the CircleCI subscriptions for the current channel.           |
| `/circleci subscription add [--flags]`    | Subscribe the current channel to CircleCI notifications.           |
| `/circleci subscription remove [--flags]` | Unsubscribe the current channel to CircleCI notifications.         |
| `/circleci subscription list-channels`    | List all channels in the current team subscribed to notifications. |

## Get insights about workflows and jobs

|                               |                                                                         |
| ----------------------------- | ----------------------------------------------------------------------- |
| `/circleci insight workflows` | Get summary of metrics for workflows over past 90 days.                 |
| `/circleci insight jobs`      | Get summary of metrics for jobs over past 90 days for a given workflow. |

## Manage CircleCI projects

|                                        |                                 |
| -------------------------------------- | ------------------------------- |
| `/circleci project list-followed`      | List followed projects.         |
| `/circleci project recent-build`       | List the 10 last builds.        |
| `/circleci project env list`           | List environment variables.     |
| `/circleci project env add name value` | Add an environment variable.    |
| `/circleci project env remove name`    | Remove an environment variable. |

## Manage pipelines

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

## Manage workflows

**Note:** Use the `/circleci pipeline workflows` command to get the workflowID.

|                             |                                     |
| --------------------------- | ----------------------------------- |
| `/circleci workflow get`    | Get information about a workflow.   |
| `/circleci workflow jobs`   | Get jobs list for a given workflow. |
| `/circleci workflow rerun`  | Rerun a given workflow.             |
| `/circleci workflow cancel` | Cancel a given workflow.            |
