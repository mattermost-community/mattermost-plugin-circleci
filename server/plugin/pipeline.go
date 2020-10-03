package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

const (
	pipelineTrigger = "pipeline"
	pipelineHint    = "<" + pipelineGetSingleTrigger + "|" + pipelineTriggerTrigger + "|" + pipelineGetAllTrigger +
		"|" + pipelineGetMineTrigger + "|" + pipelineGetRecentTrigger + "|" + pipelineWorkflowTrigger + ">"
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

	pipelineWorkflowTrigger  = "workflows"
	pipelineWorkflowHint     = "<pipeline id>"
	pipelineWorkflowHelpText = "Get list of workflows for given pipeline"

	pipelineTriggerTrigger  = "trigger"
	pipelineTriggerHint     = "<vcs-slug/org-name/repo-name> <branch>"
	pipelineTriggerHelpText = "Trigger pipeline for given project"

	pipelineGetSingleTrigger  = "get"
	pipelineGetSingleHint     = "<pipeline id>"
	pipelineGetSingleHelpText = "Get info about a single pipeline"
)

func getPipelineAutoCompeleteData() *model.AutocompleteData {
	pipeline := model.NewAutocompleteData(pipelineTrigger, pipelineHint, pipelineHelpText)
	all := model.NewAutocompleteData(pipelineGetAllTrigger, pipelineGetAllHint, pipelineGetAllHelpText)
	all.AddTextArgument("< vcs-slug/org-name/repo-name >", pipelineGetAllHint, "")
	recent := model.NewAutocompleteData(pipelineGetRecentTrigger, pipelineGetRecentHint, pipelineGetRecentHelpText)
	recent.AddTextArgument("< vcs-slug/org-name >", pipelineGetRecentHint, "")
	mine := model.NewAutocompleteData(pipelineGetMineTrigger, pipelineGetMineHint, pipelineGetMineHelpText)
	mine.AddTextArgument("< vcs-slug/org-name/repo-name >", pipelineGetMineHint, "")
	wf := model.NewAutocompleteData(pipelineWorkflowTrigger, pipelineWorkflowHint, pipelineWorkflowHelpText)
	wf.AddTextArgument("< pipeline id >", pipelineWorkflowHint, "")
	trigger := model.NewAutocompleteData(pipelineTriggerTrigger, pipelineTriggerHint, pipelineTriggerHelpText)
	trigger.AddTextArgument("< vcs-slug/org-name/repo-name >", "The repo to trigger the pipeline on. Ex: gh/mattermost/mattermost-server", "")
	trigger.AddTextArgument("<branch>", "The branch to trigger the pipeline on", "")
	get := model.NewAutocompleteData(pipelineGetSingleTrigger, pipelineGetSingleHint, pipelineGetSingleHelpText)
	get.AddTextArgument("< pipeline id >", pipelineGetSingleHint, "")
	pipeline.AddCommand(all)
	pipeline.AddCommand(recent)
	pipeline.AddCommand(mine)
	pipeline.AddCommand(wf)
	pipeline.AddCommand(trigger)
	pipeline.AddCommand(get)
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
	case pipelineWorkflowTrigger:
		return p.executePipelineGetWorkflowByID(args, circleciToken, project)
	case pipelineTriggerTrigger:
		branch := ""
		if len(split) > 2 {
			branch = split[2]
		}
		return p.executeTriggerPipeline(args, circleciToken, project, branch)
	case pipelineGetSingleTrigger:
		return p.executePipelineGetSingle(args, circleciToken, project)
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

func (p *Plugin) executePipelineGetWorkflowByID(args *model.CommandArgs,
	token string, pipelineID string) (*model.CommandResponse, *model.AppError) {
	wfs, err := circle.GetWorkflowsByPipeline(token, pipelineID)
	if err != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Failed to fetch wokflows for given pipeline ", pipelineID, err.Error())}
	}
	workflowListString := "| Name | Started By | Status | ID |\n| :---- | :----- | \n"
	for _, wf := range wfs {
		uname, _ := circle.GetNameByID(token, wf.StartedBy)
		workflowListString += fmt.Sprintf(
			"| %s | %s | %s | %s |\n",
			wf.Name,
			uname,
			wf.Status,
			wf.Id,
		)
	}

	_ = p.sendEphemeralPost(
		args,
		"Workflows for given pipeline ID: "+pipelineID,
		[]*model.SlackAttachment{
			{
				Fallback: "Workflow List",
				Text:     workflowListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeTriggerPipeline(args *model.CommandArgs,
	token string, projectSlug string, branch string) (*model.CommandResponse, *model.AppError) {
	p.API.LogError("trigger pipeline called")
	pl, err := circle.TriggerPipeline(token, projectSlug, branch)
	if err != nil {
		p.API.LogError("uparmar", err.Error())
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Could not trigger pipeline for project ", projectSlug, err.Error())}
	}
	if branch == "" {
		branch = "master"
	}
	p.API.LogError("uparmar no error")
	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				Fallback: fmt.Sprintf("Pipeline triggered successfully for project %s for branch: %s", projectSlug, branch),
				Pretext:  fmt.Sprintf("Triggered pipeline for project `%s` branch `%s`", projectSlug, branch),
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Id",
						Value: pl.Id,
						Short: true,
					},
					{
						Title: "CreatedAt",
						Value: pl.CreatedAt.String(),
						Short: true,
					},
					{
						Title: "State",
						Value: pl.State,
						Short: true,
					},
				},
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executePipelineGetSingle(args *model.CommandArgs,
	token string, pipelineID string) (*model.CommandResponse, *model.AppError) {
	pl, err := circle.GetPipelineByID(token, pipelineID)
	if err != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Could not get info about pipeline ", pipelineID, err.Error())}
	}

	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				Fallback: "Pipeline Info",
				Pretext:  "Information about pipeline: " + pipelineID,
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Id",
						Value: pl.Id,
						Short: true,
					},
					{
						Title: "Triggered By",
						Value: pl.Trigger.Actor.Login,
						Short: true,
					},
					{
						Title: "CreatedAt",
						Value: pl.CreatedAt.String(),
						Short: true,
					},
					{
						Title: "UpdatedAt",
						Value: pl.UpdatedAt.String(),
						Short: true,
					},
					{
						Title: "State",
						Value: pl.State,
						Short: true,
					},
				},
			},
		},
	)

	return &model.CommandResponse{}, nil
}
