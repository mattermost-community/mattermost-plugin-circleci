package commands

import (
	"fmt"
	"strconv"

	"github.com/jszwedko/go-circleci"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/plugin"
)

const (
	projectTrigger  = "project"
	projectHint     = "<" + projectListTrigger + "|" + projectRecentBuildsTrigger + ">"
	projectHelpText = "View informations about your CircleCI projects"

	projectListTrigger  = "list-followed"
	projectListHint     = ""
	projectListHelpText = "List followed projects"

	projectRecentBuildsTrigger  = "recent-build"
	projectRecentBuildsHint     = "<org-name> <projet-name> <branch>"
	projectRecentBuildsHelpText = "List the 10 last builds for a project"

	projectHelpMessage = "## Manage your CircleCI projects" // TODO
)

func getProjectAutocompleteData() *model.AutocompleteData {
	project := model.NewAutocompleteData(projectTrigger, projectHint, projectHelpText)

	projectList := model.NewAutocompleteData(projectListTrigger, projectListHint, projectListHelpText)
	projectRecentBuild := model.NewAutocompleteData(projectRecentBuildsTrigger, projectRecentBuildsHint, projectRecentBuildsHelpText)
	projectRecentBuild.AddTextArgument("Owner of the project's repository", "[username]", "")
	projectRecentBuild.AddDynamicListArgument("", plugin.RouteAutocompleteFollowedProjects, true)
	projectRecentBuild.AddTextArgument("Branch name", "[branch]", "")

	project.AddCommand(projectRecentBuild)
	project.AddCommand(projectList)

	return project
}

func executeProject(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case projectListTrigger:
		return executeProjectList(args, circleciToken)

	case projectRecentBuildsTrigger:
		return executeProjectRecentBuilds(args, circleciToken, split[1:])

	case helpTrigger:
		return formatHelpMessage(args, projectTrigger)

	default:
		return formatIncorrectSubcommand(args, projectTrigger)
	}
}

func executeProjectList(args *model.CommandArgs, circleciToken string) (*model.CommandResponse, *model.AppError) {
	projects, ok := getCircleciUserProjects(circleciToken)
	if !ok {
		return sendEphemeralResponse(args, errorConnectionText), nil
	}

	projectsListString := "| Project | CircleCI URL | Is [OSS](https://circleci.com/open-source/) |\n| :---- | :----- | :---- | \n"
	for _, project := range projects {
		// TODO : add environment variables

		projectsListString += fmt.Sprintf(
			"| [%s/%s](%s) | %s | %t |\n",
			project.Username,
			project.Reponame,
			project.VCSURL,
			fmt.Sprintf("https://circleci.com/gh/%s/%s", project.Username, project.Reponame), // TODO : handle bitbucket URL
			project.FeatureFlags.OSS,
		)
	}

	_ = sendEphemeralPost(
		args,
		"Projects you are following on CircleCI",
		[]*model.SlackAttachment{
			{
				Fallback: "Projects list",
				Text:     projectsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func executeProjectRecentBuilds(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	client := &circleci.Client{Token: circleciToken}

	if len(split) < 3 {
		return sendEphemeralResponse(args, "Please provide the project username, repository and branch name)"), nil
	}

	account, repo, branch := split[0], split[1], split[2]
	builds, err := client.ListRecentBuildsForProject(account, repo, branch, "", 10, 0)
	if err != nil {
		API.LogError("Unable to get recent build from CircleCI", "CircleCI error", err)
		return sendEphemeralResponse(args, errorConnectionText), nil
	}

	text := "| Workflow | Job | Build | Subject | Start time | Status | Duration | Triggered by|\n| :---- | :----- | :----- | :----- | :----- | :----- | :---- | \n"
	for _, build := range builds {
		buildStartTime := buildStartTimeToString(build)

		buildTime := "/"
		if build.BuildTimeMillis != nil {
			buildTime = strconv.Itoa(*build.BuildTimeMillis/1000) + "s"
		}

		statusImageMarkdown := buildStatusToMarkdown(build, p)

		text += fmt.Sprintf("| % s | % s | [%d](%s) | `%s` | %s | %s | %s | %s |\n",
			build.Workflows.WorkflowName,
			build.Workflows.JobName,
			build.BuildNum,
			build.BuildURL,
			build.Subject,
			buildStartTime,
			statusImageMarkdown,
			buildTime,
			build.Why,
		)
	}

	_ = sendEphemeralPost(
		args,
		"Recent builds for "+account+"/"+repo+" "+branch,
		[]*model.SlackAttachment{
			{
				Fallback: "Recent builds list",
				Text:     text,
			},
		},
	)

	return &model.CommandResponse{}, nil
}
