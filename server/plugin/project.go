package plugin

import (
	"fmt"
	"strconv"

	"github.com/jszwedko/go-circleci"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"
)

const (
	projectTrigger  = "project"
	projectHint     = "<" + projectListTrigger + "|" + projectRecentBuildsTrigger + ">"
	projectHelpText = "View informations about your CircleCI projects"

	projectListTrigger  = "list-followed"
	projectListHint     = ""
	projectListHelpText = "List followed projects"

	projectRecentBuildsTrigger = "recent-build"
	// TODO rename in all files (code and UI strings) 'username' to 'owner' for repository
	projectRecentBuildsHint     = "<username> <repository> <branch>"
	projectRecentBuildsHelpText = "List the 10 last builds for a project"

	projectEnvVarTrigger  = "envvar"
	projectEnvVarHint     = "<" + projectEnvVarListTrigger + "|" + projectEnvVarAddTrigger + "|" + projectEnvVarAddTrigger + ">"
	projectEnvVarHelpText = "get, add or remove environment varibales for given project"

	projectEnvVarListTrigger  = "list"
	projectEnvVarListHint     = "<vcs-slug/org-name/repo-name>"
	projectEnvVarListHelpText = "List all environment variables for given project"

	projectEnvVarAddTrigger  = "add"
	projectEnvVarAddHint     = "<vcs-slug/org-name/repo-name> <env var name> <value>"
	projectEnvVarAddHelpText = "Add a new environment varibale for a project"

	projectEnvVarDelTrigger  = "remove"
	projectEnvVarDelHint     = "<vcs-slug/org-name/repo-name> <env var name>"
	projectEnvVarDelHelpText = "Delete an environment varibale for a project"
)

func getProjectAutoComplete() *model.AutocompleteData {
	project := model.NewAutocompleteData(projectTrigger, projectHint, projectHelpText)

	projectList := model.NewAutocompleteData(projectListTrigger, projectListHint, projectListHelpText)
	projectRecentBuild := model.NewAutocompleteData(projectRecentBuildsTrigger, projectRecentBuildsHint, projectRecentBuildsHelpText)
	projectRecentBuild.AddTextArgument("Owner of the project's repository", "[username]", "")
	projectRecentBuild.AddDynamicListArgument("", routeAutocomplete+subrouteFollowedProjects, true)
	projectRecentBuild.AddTextArgument("Branch name", "[branch]", "")

	envvar := model.NewAutocompleteData(projectEnvVarTrigger, projectEnvVarHint, projectEnvVarHelpText)
	list := model.NewAutocompleteData(projectEnvVarListTrigger, projectEnvVarListHint, projectEnvVarListHelpText)
	list.AddTextArgument("<vcs-slug/org-name/repo-name>", "The repo to get env vars of. Ex: gh/mattermost/mattermost-server", "")
	add := model.NewAutocompleteData(projectEnvVarAddTrigger, projectEnvVarAddHint, projectEnvVarAddHelpText)
	add.AddTextArgument("<vcs-slug/org-name/repo-name>", "Project slug. Ex:gh/mattermost/mattermost-server", "")
	add.AddTextArgument("<env var name> <value>", "Name and value of environment variable to add. Ex: testVar testVal", "")
	del := model.NewAutocompleteData(projectEnvVarDelTrigger, projectEnvVarDelHint, projectEnvVarDelHelpText)
	del.AddTextArgument("<vcs-slug/org-name/repo-name>", "Project slug. Ex:gh/mattermost/mattermost-server", "")
	del.AddTextArgument("<env var name>", "Name and value of environment variable to remove. Ex: testVar", "")

	envvar.AddCommand(list)
	envvar.AddCommand(add)
	envvar.AddCommand(del)

	project.AddCommand(projectRecentBuild)
	project.AddCommand(projectList)
	project.AddCommand(envvar)

	return project
}

func (p *Plugin) executeProject(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case projectListTrigger:
		return p.executeProjectList(args, circleciToken)

	case projectRecentBuildsTrigger:
		return p.executeProjectRecentBuilds(args, circleciToken, split[1:])

	case commandHelpTrigger:
		return p.sendHelpResponse(args, projectTrigger)

	case projectEnvVarTrigger:
		subsubcmd := "list"
		if len(split) > 1 {
			subsubcmd = split[1]
		}
		switch subsubcmd {
		case projectEnvVarListTrigger:
			return p.executeProjectListEnvVars(args, circleciToken, split[2:])
		case projectEnvVarAddTrigger:
			return p.executeProjectAddEnvVar(args, circleciToken, split[1:])
		case projectEnvVarDelTrigger:
			return p.executeProjectDelEnvVar(args, circleciToken, split[1:])
		default:
			return p.sendIncorrectSubcommandResponse(args, projectEnvVarListTrigger)
		}
	default:
		return p.sendIncorrectSubcommandResponse(args, projectTrigger)
	}
}

func (p *Plugin) executeProjectList(args *model.CommandArgs, circleciToken string) (*model.CommandResponse, *model.AppError) {
	projects, err := v1.GetCircleciUserProjects(circleciToken)
	if err != nil {
		return p.sendEphemeralResponse(args, errorConnectionText), nil
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

	_ = p.sendEphemeralPost(
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

func (p *Plugin) executeProjectRecentBuilds(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	client := &circleci.Client{Token: circleciToken}

	if len(split) < 3 {
		return p.sendEphemeralResponse(args, "Please provide the project username, repository and branch name)"), nil
	}

	account, repo, branch := split[0], split[1], split[2]
	builds, err := client.ListRecentBuildsForProject(account, repo, branch, "", 10, 0)
	if err != nil {
		p.API.LogError("Unable to get recent build from CircleCI", "CircleCI error", err)
		return p.sendEphemeralResponse(args, errorConnectionText), nil
	}

	text := "| Workflow | Job | Build | Subject | Start time | Status | Duration | Triggered by|\n| :---- | :----- | :----- | :----- | :----- | :----- | :---- | \n"
	for _, build := range builds {
		buildStartTime := v1.BuildStartTimeToString(build)

		buildTime := "/"
		if build.BuildTimeMillis != nil {
			buildTime = strconv.Itoa(*build.BuildTimeMillis/1000) + "s"
		}

		statusImageMarkdown := v1.BuildStatusToMarkdown(build, badgePassedURL, badgeFailedURL)

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

	_ = p.sendEphemeralPost(
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

func (p *Plugin) executeProjectListEnvVars(args *model.CommandArgs,
	token string, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 1 {
		return p.sendEphemeralResponse(args, "Project Slug cannot be empty for list command"),
			&model.AppError{Message: "received empty project slug"}
	}
	envvars, err := circle.GetEnvVarsList(token, split[0])
	if err != nil {
		return p.sendEphemeralResponse(args, fmt.Sprintf("Could not list environment varibales for ptoject %s", split[0])),
			&model.AppError{Message: "Could not list env vars for project" + split[0] + "err: " + err.Error()}
	}
	envVarListString := "| Name | Value |\n| :---- | :----- | \n"
	for _, env := range envvars {
		envVarListString += fmt.Sprintf(
			"| %s | %s |\n",
			env.Name,
			env.Value,
		)
	}

	_ = p.sendEphemeralPost(
		args,
		"Environment variables for project "+split[0],
		[]*model.SlackAttachment{
			{
				Fallback: "Environment Varibale List",
				Text:     envVarListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeProjectAddEnvVar(args *model.CommandArgs,
	token string, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 3 {
		return p.sendEphemeralResponse(args, "Please provide project slug, varibale name and value"),
			&model.AppError{Message: "received empty project slug or variable name or value"}
	}
	err := circle.AddEnvVar(token, split[0], split[1], split[2])
	if err != nil {
		return p.sendEphemeralResponse(args, fmt.Sprintf("Could not add environment varibale",
				"`%s: %s` for ptoject %s", split[1], split[2], split[0])),
			&model.AppError{Message: "Could not add env var %s:%s for project %s" + split[1] + split[2] +
				split[0] + "err: " + err.Error()}
	}

	return p.sendEphemeralResponse(args, fmt.Sprintf("Succesfully added environment variable `%s:%s` for project %s",
		split[1], split[2], split[0])), nil
}

func (p *Plugin) executeProjectDelEnvVar(args *model.CommandArgs,
	token string, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(args, "Please provide project slug and varibale name"),
			&model.AppError{Message: "received empty project slug or variable name"}
	}
	err := circle.DelEnvVar(token, split[0], split[1])
	if err != nil {
		return p.sendEphemeralResponse(args, fmt.Sprintf("Could not remove environment varibale",
				"`%s` for ptoject %s", split[1], split[0])),
			&model.AppError{Message: "Could not remove env var %s for project %s" + split[1] +
				split[0] + "err: " + err.Error()}
	}

	return p.sendEphemeralResponse(args, fmt.Sprintf("Succesfully removed environment variable `%s` for project %s",
		split[1], split[2], split[0])), nil
}
