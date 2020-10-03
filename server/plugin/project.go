package plugin

import (
	"fmt"
	"strconv"

	"github.com/jszwedko/go-circleci"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	projectTrigger  = "project"
	projectHint     = "<" + projectListTrigger + "|" + projectRecentBuildsTrigger + "|" + projectEnvVarTrigger + ">"
	projectHelpText = "View informations about your CircleCI projects"

	projectListTrigger  = "list-followed"
	projectListHint     = ""
	projectListHelpText = "List followed projects"

	projectRecentBuildsTrigger = "recent-build"
	// TODO rename in all files (code and UI strings) 'username' to 'organization' or 'org' for repository
	projectRecentBuildsHint     = "<username> <repository> <branch>"
	projectRecentBuildsHelpText = "List the 10 last builds for a project"

	projectEnvVarTrigger  = "env"
	projectEnvVarHint     = "<" + projectEnvVarListTrigger + "|" + projectEnvVarAddTrigger + "|" + projectEnvVarAddTrigger + ">"
	projectEnvVarHelpText = "get, add or remove environment variables for given project"

	projectEnvVarListTrigger  = "list"
	projectEnvVarListHint     = "<vcs-slug/org-name/repo-name>"
	projectEnvVarListHelpText = "List all environment variables for given project"

	projectEnvVarAddTrigger  = "add"
	projectEnvVarAddHint     = "<vcs-slug/org-name/repo-name> <env var name> <value>"
	projectEnvVarAddHelpText = "Add a new environment variable for a project"

	projectEnvVarDelTrigger  = "remove"
	projectEnvVarDelHint     = "<vcs-slug/org-name/repo-name> <env var name>"
	projectEnvVarDelHelpText = "Delete an environment variable for a project"
)

func getProjectAutoComplete() *model.AutocompleteData {
	// TODO : update autocomplete
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
	add.AddTextArgument("<env var name> ", "Name of environment variable to add. Ex: testVar", "")
	add.AddTextArgument("<env var value> ", "Value of environment variable to add. Ex: testVal", "")
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

func (p *Plugin) executeProject(args *model.CommandArgs, circleciToken string, config *store.Config, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case projectListTrigger:
		return p.executeProjectList(args, circleciToken)

	case projectRecentBuildsTrigger:
		return p.executeProjectRecentBuilds(args, circleciToken, config, split[1:])

	case commandHelpTrigger:
		return p.sendHelpResponse(args, projectTrigger)

	case projectEnvVarTrigger:
		subsubcmd := "list"
		if len(split) > 1 {
			subsubcmd = split[1]
		}
		switch subsubcmd {
		case projectEnvVarListTrigger:
			return p.executeProjectListEnvVars(args, circleciToken, config)
		case projectEnvVarAddTrigger:
			return p.executeProjectAddEnvVar(args, circleciToken, config, split[2:])
		case projectEnvVarDelTrigger:
			return p.executeProjectDelEnvVar(args, circleciToken, config, split[2:])
		default:
			return p.sendIncorrectSubcommandResponse(args, projectEnvVarTrigger)
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

func (p *Plugin) executeProjectRecentBuilds(args *model.CommandArgs, circleciToken string, config *store.Config, split []string) (*model.CommandResponse, *model.AppError) {
	client := &circleci.Client{Token: circleciToken}

	if len(split) < 1 {
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Please precise the branch name. Selected project: %s", config.ToMarkdown()),
		), nil
	}

	branch := split[0]
	builds, err := client.ListRecentBuildsForProject(config.Org, config.Project, branch, "", 10, 0)
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
		fmt.Sprintf("Recent builds for %s branch `%s`", config.ToMarkdown(), branch),
		[]*model.SlackAttachment{
			{
				Fallback: "Recent builds list",
				Text:     text,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeProjectListEnvVars(args *model.CommandArgs, token string, config *store.Config) (*model.CommandResponse, *model.AppError) {
	envvars, err := circle.GetEnvVarsList(token, config.ToSlug())
	if err != nil {
		p.API.LogError("Could not list env vars", "error", err.Error(), "project", config.ToSlug())
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Could not list environment variables for project %s", config.ToMarkdown()),
		), nil
	}

	if len(envvars) == 0 {
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Project %s does not have any environment variables", config.ToMarkdown()),
		), nil
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
		fmt.Sprintf("Environment variables for project %s", config.ToMarkdown()),
		[]*model.SlackAttachment{
			{
				Fallback: "Environment Variable List",
				Text:     envVarListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeProjectAddEnvVar(args *model.CommandArgs, token string, config *store.Config, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(args, "Please provide the variable name and value"), nil
	}

	varName := split[0]
	varValue := split[1]

	err := circle.AddEnvVar(token, config.ToSlug(), varName, varValue)
	if err != nil {
		p.API.LogError("Unable to set CircleCI envVar", "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Could not add environment variable `%s: %s` for project %s", varName, varValue, config.ToMarkdown()),
		), nil
	}

	return p.sendEphemeralResponse(args,
		fmt.Sprintf("Successfully added environment variable `%s:%s` for project %s", varName, varValue, config.ToMarkdown()),
	), nil
}

func (p *Plugin) executeProjectDelEnvVar(args *model.CommandArgs, token string, config *store.Config, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 1 {
		return p.sendEphemeralResponse(args, "Please provide the variable name"), nil
	}

	varName := split[0]

	err := circle.DelEnvVar(token, config.ToSlug(), varName)
	if err != nil {
		p.API.LogError("Could not remove env var for project", "error", err, "env var", varName)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Could not remove environment variable `%s` for project %s", varName, config.ToMarkdown()),
		), nil
	}

	return p.sendEphemeralResponse(args,
		fmt.Sprintf("Successfully removed environment variable `%s` for project %s", varName, config.ToMarkdown()),
	), nil
}
