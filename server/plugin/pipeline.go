package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

const (
	pipelineTrigger  = "pipeline"
	pipelineHint     = "<" + pipelineGetAllTrigger + "|" + pipelineGetMineTrigger + "|" + pipelineGetRecentTrigger + ">"
	pipelineHelpText = "Manage the connection to your CircleCI acccount"

	pipelineGetRecentTrigger  = "recent"
	pipelineGetRecentHint     = "<vcs-slug/org-name>"
	pipelineGetRecentHelpText = "Get list of all recently run pipelines"

	pipelineGetAllTrigger  = "all"
	pipelineGetAllHint     = "<vcs-slug/org-name/repo-name>"
	pipelineGetAllHelpText = "Get list of all pipelines for a given project"

	pipelineGetMineTrigger  = "mine"
	pipelineGetMineHint     = "<vcs-slug/org-name/repo-name>"
	pipelineGetMineHelpText = "Get list of all my pipelines triggered by you"
)

func getPipelineAutoCompeleteData() *model.AutocompleteData {
	pipeline := model.NewAutocompleteData(pipelineTrigger, pipelineHint, pipelineHelpText)
	all := model.NewAutocompleteData(pipelineGetAllTrigger, pipelineGetAllHint, pipelineGetAllHelpText)
	all.AddTextArgument("< vcs-slug/org-name/repo-name >", pipelineGetAllHint, "")
	recent := model.NewAutocompleteData(pipelineGetRecentTrigger, pipelineGetRecentHint, pipelineGetRecentHelpText)
	recent.AddTextArgument("< vcs-slug/org-name >", pipelineGetRecentHint, "")
	mine := model.NewAutocompleteData(pipelineGetMineTrigger, pipelineGetMineHint, pipelineGetMineHelpText)
	mine.AddTextArgument("< vcs-slug/org-name/repo-name >", pipelineGetMineHint, "")
	pipeline.AddCommand(all)
	pipeline.AddCommand(recent)
	pipeline.AddCommand(mine)
	return pipeline
}

func (p *Plugin) executePipelineTrigger(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	var project string
	if len(split) > 1 {
		project = split[1]
	} else {
		return p.sendIncorrectSubcommandResponse(args, pipelineTrigger)
	}

	switch subcommand {
	case pipelineGetAllTrigger:
		return p.executePipelineGetAllForProject(args, circleciToken, project)
	case pipelineGetRecentTrigger:
		return p.executePipelineGetRecent(args, circleciToken, project)
	case pipelineGetMineTrigger:
		return p.executePipelineGetAllForProjectByMe(args, circleciToken, project)
	default:
		return p.sendIncorrectSubcommandResponse(args, pipelineTrigger)
	}
}

func (p *Plugin) executePipelineGetRecent(args *model.CommandArgs,
	token string, orgSlug string) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetRecentlyBuiltPipelines(token, orgSlug, false)
	if err != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Failed to fetch info for pipeline", orgSlug, err.Error())}
	}
	pipelineListString := "| Pipeline ID | State |\n| :---- | :----- | \n"
	for _, pipeline := range pipelines {
		pipelineListString += fmt.Sprintf(
			"| %s | %s |\n",
			pipeline.Id,
			pipeline.State,
		)
	}

	_ = p.sendEphemeralPost(
		args,
		"Recently built pipelines in your org",
		[]*model.SlackAttachment{
			{
				Fallback: "Pipelines list",
				Text:     pipelineListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executePipelineGetAllForProject(args *model.CommandArgs,
	token string, projectSlug string) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetAllPipelinesForProject(token, projectSlug)
	if err != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Failed to fetch info for pipeline", projectSlug, err.Error())}
	}
	projectsListString := "| Pipeline ID | State |\n| :---- | :----- | \n"
	for _, pipeline := range pipelines {
		projectsListString += fmt.Sprintf(
			"| %s | %s |\n",
			pipeline.Id,
			pipeline.State,
		)
	}

	pr := strings.Split(projectSlug, "/")

	_ = p.sendEphemeralPost(
		args,
		"Recently built pipelines for project "+pr[2],
		[]*model.SlackAttachment{
			{
				Fallback: "Projects list",
				Text:     projectsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil

}

func (p *Plugin) executePipelineGetAllForProjectByMe(args *model.CommandArgs,
	token string, projectSlug string) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetAllMyPipelinesForProject(token, projectSlug)
	if err != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Failed to fetch info for pipeline", projectSlug, err.Error())}
	}
	projectsListString := "| Pipeline ID | State |\n| :---- | :----- | \n"
	for _, pipeline := range pipelines {
		projectsListString += fmt.Sprintf(
			"| %s | %s |\n",
			pipeline.Id,
			pipeline.State,
		)
	}

	pr := strings.Split(projectSlug, "/")

	_ = p.sendEphemeralPost(
		args,
		"Recently built pipelines by you for project "+pr[2],
		[]*model.SlackAttachment{
			{
				Fallback: "Projects list",
				Text:     projectsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}
