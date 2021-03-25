package plugin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jszwedko/go-circleci"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-circleci/server/circle"
	v1 "github.com/mattermost/mattermost-plugin-circleci/server/circle/v1"
	"github.com/mattermost/mattermost-plugin-circleci/server/store"
)

const (
	projectTrigger  = "project"
	projectHint     = "<" + projectListTrigger + "|" + projectRecentBuildsTrigger + "|" + projectEnvVarTrigger + ">"
	projectHelpText = "View informations about your CircleCI projects"

	projectListTrigger  = "list-followed"
	projectListHint     = ""
	projectListHelpText = "List followed projects"

	projectRecentBuildsTrigger  = "recent-build"
	projectRecentBuildsHint     = "<branch>"
	projectRecentBuildsHelpText = "List the 10 last builds"

	projectEnvVarTrigger  = "env"
	projectEnvVarHint     = "<" + projectEnvVarListTrigger + "|" + projectEnvVarAddTrigger + "|" + projectEnvVarAddTrigger + ">"
	projectEnvVarHelpText = "get, add or remove environment variables for given project"

	projectEnvVarListTrigger  = "list"
	projectEnvVarListHint     = ""
	projectEnvVarListHelpText = "List all environment variables"

	projectEnvVarAddTrigger  = "add"
	projectEnvVarAddHint     = "<env-var name> <value>"
	projectEnvVarAddHelpText = "Add a new environment variable"

	projectEnvVarDelTrigger  = "remove"
	projectEnvVarDelHint     = "<env-var name>"
	projectEnvVarDelHelpText = "Delete an environment variable"
)

func getProjectAutoComplete() *model.AutocompleteData {
	project := model.NewAutocompleteData(projectTrigger, projectHint, projectHelpText)

	projectList := model.NewAutocompleteData(projectListTrigger, projectListHint, projectListHelpText)

	projectRecentBuild := model.NewAutocompleteData(projectRecentBuildsTrigger, projectRecentBuildsHint, projectRecentBuildsHelpText)
	projectRecentBuild.AddTextArgument("The branch to get recent builds from", "<branch>", "")
	projectRecentBuild.AddNamedDynamicListArgument(namedArgProjectName, namedArgProjectHelpText, routeAutocomplete+subrouteFollowedProjects, false)

	envvar := model.NewAutocompleteData(projectEnvVarTrigger, projectEnvVarHint, projectEnvVarHelpText)
	envvarList := model.NewAutocompleteData(projectEnvVarListTrigger, projectEnvVarListHint, projectEnvVarListHelpText)
	envvarList.AddNamedDynamicListArgument(namedArgProjectName, namedArgProjectHelpText, routeAutocomplete+subrouteFollowedProjects, false)
	envvarAdd := model.NewAutocompleteData(projectEnvVarAddTrigger, projectEnvVarAddHint, projectEnvVarAddHelpText)
	envvarAdd.AddTextArgument("<env-var name> ", "Name of environment variable to add. Ex: testVar", "")
	envvarAdd.AddTextArgument("<env-var value> ", "Value of environment variable to add. Ex: testVal", "")
	envvarAdd.AddNamedDynamicListArgument(namedArgProjectName, namedArgProjectHelpText, routeAutocomplete+subrouteFollowedProjects, false)
	envvarDel := model.NewAutocompleteData(projectEnvVarDelTrigger, projectEnvVarDelHint, projectEnvVarDelHelpText)
	envvarDel.AddTextArgument("<env-var name>", "Name and value of environment variable to remove. Ex: testVar", "")
	envvarDel.AddNamedDynamicListArgument(namedArgProjectName, namedArgProjectHelpText, routeAutocomplete+subrouteFollowedProjects, false)

	envvar.AddCommand(envvarList)
	envvar.AddCommand(envvarAdd)
	envvar.AddCommand(envvarDel)

	project.AddCommand(envvar)
	project.AddCommand(projectRecentBuild)
	project.AddCommand(projectList)

	return project
}

func (p *Plugin) executeProject(args *model.CommandArgs, circleciToken string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case projectListTrigger:
		return p.executeProjectList(args, circleciToken)

	case projectRecentBuildsTrigger:
		return p.executeProjectRecentBuilds(args, circleciToken, project, split[1:])

	case commandHelpTrigger:
		return p.sendHelpResponse(args, projectTrigger)

	case projectEnvVarTrigger:
		subsubcmd := "list"
		if len(split) > 1 {
			subsubcmd = split[1]
		}
		switch subsubcmd {
		case projectEnvVarListTrigger:
			return p.executeProjectListEnvVars(args, circleciToken, project)
		case projectEnvVarAddTrigger:
			return p.executeProjectAddEnvVar(args, circleciToken, project, split[2:])
		case projectEnvVarDelTrigger:
			return p.executeProjectDelEnvVar(args, circleciToken, project, split[2:])
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
		VCSType := "gh"
		if strings.Contains(project.VCSURL, "https://bitbucket.org") {
			VCSType = "bb"
		}

		projectsListString += fmt.Sprintf(
			"| [%s/%s](%s) | %s | %t |\n",
			project.Username,
			project.Reponame,
			project.VCSURL,
			fmt.Sprintf("https://circleci.com/%s/%s/%s", VCSType, project.Username, project.Reponame),
			project.FeatureFlags.OSS,
		)
	}

	p.sendEphemeralPost(
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

func (p *Plugin) executeProjectRecentBuilds(args *model.CommandArgs, circleciToken string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
	client := &circleci.Client{Token: circleciToken}

	if len(split) < 1 {
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Please precise the branch name. Selected project: %s", project.ToMarkdown()),
		), nil
	}

	branch := split[0]
	builds, err := client.ListRecentBuildsForProject(project.Org, project.Project, branch, "", 10, 0)
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

		buildSubject := "/"
		if build.Subject != "" {
			buildSubject = build.Subject
		}

		statusImageMarkdown := v1.BuildStatusToMarkdown(build, badgePassedURL, badgeFailedURL)

		text += fmt.Sprintf("| % s | % s | [%d](%s) | %s | %s | %s | %s | %s |\n",
			build.Workflows.WorkflowName,
			build.Workflows.JobName,
			build.BuildNum,
			build.BuildURL,
			buildSubject,
			buildStartTime,
			statusImageMarkdown,
			buildTime,
			build.Why,
		)
	}

	p.sendEphemeralPost(
		args,
		fmt.Sprintf("Recent builds for %s branch `%s`", project.ToMarkdown(), branch),
		[]*model.SlackAttachment{
			{
				Fallback: "Recent builds list",
				Text:     text,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeProjectListEnvVars(args *model.CommandArgs, token string, project *store.ProjectIdentifier) (*model.CommandResponse, *model.AppError) {
	envvars, err := circle.GetEnvVarsList(token, project.ToSlug())
	if err != nil {
		p.API.LogError("Could not list env vars", "error", err.Error(), "project", project.ToSlug())
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Could not list environment variables for project %s", project.ToMarkdown()),
		), nil
	}

	if len(envvars) == 0 {
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":information_source: Project %s does not have any environment variables", project.ToMarkdown()),
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

	p.sendEphemeralPost(
		args,
		fmt.Sprintf("Environment variables for %s", project.ToMarkdown()),
		[]*model.SlackAttachment{
			{
				Fallback: "Environment Variable List",
				Text:     envVarListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeProjectAddEnvVar(args *model.CommandArgs, token string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(args, "Please provide the variable name and value"), nil
	}

	varName := split[0]
	varValue := split[1]

	val, exist, errr := circle.EnvVarExists(token, project.ToSlug(), varName)
	if errr != nil {
		p.API.LogError("err while getting env var details", errr.Error())
	}
	if exist {
		attach := model.SlackAttachment{}
		attach.Actions = []*model.PostAction{
			{
				Id:   "envoverwrite",
				Name: "Approve",
				Type: model.POST_ACTION_TYPE_BUTTON,
				Integration: &model.PostActionIntegration{
					URL: fmt.Sprintf("/plugins/%s/%s", manifest.Id, routeEnvOverwrite),
					Context: map[string]interface{}{
						"EnvName":     varName,
						"EnvVal":      varValue,
						"ProjectSlug": project.ToSlug(),
						"Action":      "approve",
					},
				},
			},
			{
				Id:   "envoverwritedeny",
				Name: "Deny",
				Type: model.POST_ACTION_TYPE_BUTTON,
				Integration: &model.PostActionIntegration{
					URL: fmt.Sprintf("/plugins/%s/%s", manifest.Id, routeEnvOverwrite),
					Context: map[string]interface{}{
						"EnvName":     varName,
						"EnvVal":      varValue,
						"ProjectSlug": project.ToSlug(),
						"Action":      "deny",
					},
				},
			},
		}
		attach.Title = fmt.Sprintf("Do you want to overwrite environment variable `%s` with masked value `%s`?", val.Name, val.Value)
		attach.Fallback = attach.Title
		attach.Color = "#8267E4" // purple
		attach.Pretext = attach.Title

		p.sendEphemeralPost(args, "", []*model.SlackAttachment{&attach})

		return &model.CommandResponse{}, nil
	}
	err := circle.AddEnvVar(token, project.ToSlug(), varName, varValue)
	if err != nil {
		p.API.LogError("Unable to set CircleCI envVar", "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Could not add environment variable `%s: %s` for project %s", varName, varValue, project.ToMarkdown()),
		), nil
	}

	return p.sendEphemeralResponse(args,
		fmt.Sprintf(":white_check_mark: Successfully added environment variable `%s:%s` for project %s", varName, varValue, project.ToMarkdown()),
	), nil
}

func (p *Plugin) executeProjectDelEnvVar(args *model.CommandArgs, token string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 1 {
		return p.sendEphemeralResponse(args, "Please provide the variable name"), nil
	}

	varName := split[0]

	err := circle.DelEnvVar(token, project.ToSlug(), varName)
	if err != nil {
		p.API.LogError("Could not remove env var for project", "error", err, "env var", varName)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Could not remove environment variable `%s` for project %s", varName, project.ToMarkdown()),
		), nil
	}

	return p.sendEphemeralResponse(args,
		fmt.Sprintf(":white_check_mark: Successfully removed environment variable `%s` for project %s", varName, project.ToMarkdown()),
	), nil
}
