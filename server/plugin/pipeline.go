package plugin

import (
	"fmt"
	"strconv"

	"github.com/darkLord19/circleci-v2/circleci"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	pipelineTrigger = "pipeline"
	pipelineHint    = "<" + pipelineGetSingleTrigger + "|" + pipelineTriggerTrigger + "|" + pipelineGetAllTrigger +
		"|" + pipelineGetMineTrigger + "|" + pipelineGetRecentTrigger + "|" + pipelineWorkflowTrigger + ">"
	pipelineHelpText = "Manage the connection to your CircleCI acccount"

	pipelineGetRecentTrigger  = "recent"
	pipelineGetRecentHint     = "<vcs-slug/org-name>"
	pipelineGetRecentHelpText = "Get list of all recently ran pipelines"

	pipelineGetAllTrigger  = "all"
	pipelineGetAllHint     = ""
	pipelineGetAllHelpText = "Get list of all pipelines for a project"

	pipelineGetMineTrigger  = "mine"
	pipelineGetMineHint     = ""
	pipelineGetMineHelpText = "Get list of all pipelines triggered by you for a project"

	pipelineTriggerTrigger  = "trigger"
	pipelineTriggerHint     = "<branch>"
	pipelineTriggerHelpText = "Trigger pipeline for a project"

	pipelineWorkflowTrigger  = "workflows"
	pipelineWorkflowHint     = "<pipeline number>"
	pipelineWorkflowHelpText = "Get list of workflows for given pipeline id"

	pipelineGetSingleTrigger  = "get"
	pipelineGetSingleHint     = "<pipeline id> or <pipeline number>"
	pipelineGetSingleHelpText = "Get informations about a single pipeline for a given project or a pipeline id"
)

func getPipelineAutoCompeleteData() *model.AutocompleteData {
	pipeline := model.NewAutocompleteData(pipelineTrigger, pipelineHint, pipelineHelpText)

	all := model.NewAutocompleteData(pipelineGetAllTrigger, pipelineGetAllHint, pipelineGetAllHelpText)
	all.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	recent := model.NewAutocompleteData(pipelineGetRecentTrigger, pipelineGetRecentHint, pipelineGetRecentHelpText)
	recent.AddTextArgument("VCS is either bb or gh. Leave blank for default org", pipelineGetRecentHint, "")

	mine := model.NewAutocompleteData(pipelineGetMineTrigger, pipelineGetMineHint, pipelineGetMineHelpText)
	mine.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	wf := model.NewAutocompleteData(pipelineWorkflowTrigger, pipelineWorkflowHint, pipelineWorkflowHelpText)
	wf.AddTextArgument("<pipelineID>", pipelineWorkflowHint, "")

	trigger := model.NewAutocompleteData(pipelineTriggerTrigger, pipelineTriggerHint, pipelineTriggerHelpText)
	trigger.AddTextArgument("<branch>", "The branch to trigger the pipeline on. Leave empty for master", "")
	trigger.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	get := model.NewAutocompleteData(pipelineGetSingleTrigger, pipelineGetSingleHint, pipelineGetSingleHelpText)
	get.AddTextArgument("< pipeline number > or < pipelineID >", pipelineGetSingleHint, "")

	pipeline.AddCommand(trigger)
	pipeline.AddCommand(get)
	pipeline.AddCommand(wf)
	pipeline.AddCommand(all)
	pipeline.AddCommand(mine)
	pipeline.AddCommand(recent)

	return pipeline
}

func (p *Plugin) executePipelineTrigger(args *model.CommandArgs, circleciToken string, config *store.Config, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := commandHelpTrigger
	if len(split) > 0 {
		subcommand = split[0]
	}

	var argument string
	if len(split) > 1 {
		argument = split[1]
	}

	switch subcommand {
	case pipelineGetAllTrigger:
		return p.executePipelineGetAllForProject(args, circleciToken, config)
	case pipelineGetRecentTrigger:
		return p.executePipelineGetRecent(args, circleciToken, config)
	case pipelineGetMineTrigger:
		return p.executePipelineGetAllForProjectByMe(args, circleciToken, config)
	case pipelineWorkflowTrigger:
		return p.executePipelineGetWorkflowByID(args, circleciToken, argument)
	case pipelineTriggerTrigger:
		branch := ""
		if len(split) > 2 {
			branch = split[2]
		}
		return p.executeTriggerPipeline(args, circleciToken, config, branch)
	case pipelineGetSingleTrigger:
		return p.executePipelineGetSingle(args, circleciToken, config, argument)

	case commandHelpTrigger:
		return p.sendHelpResponse(args, pipelineTrigger)
	default:
		return p.sendIncorrectSubcommandResponse(args, pipelineTrigger)
	}
}

func (p *Plugin) executePipelineGetRecent(args *model.CommandArgs, token string,
	config *store.Config) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetRecentlyBuiltPipelines(token, fmt.Sprintf("%s/%s", config.VCSType, config.Org), false)
	if err != nil {
		p.API.LogError("Failed to fetch info for pipeline", "org", config.ToSlug(), "error", err.Error())
		return p.sendEphemeralResponse(args, "Failed to fetch info for pipeline"), nil
	}

	pipelineListString := "| Pipeline No. | Pipeline ID | State |\n| :---- | :----- | \n"
	for _, pipeline := range pipelines {
		pipelineListString += fmt.Sprintf(
			"| %d | %s | %s |\n",
			pipeline.Number,
			pipeline.Id,
			pipeline.State,
		)
	}

	_ = p.sendEphemeralPost(
		args,
		"Recently built pipelines in your organizaition",
		[]*model.SlackAttachment{
			{
				Fallback: "Pipelines list",
				Text:     pipelineListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executePipelineGetAllForProject(args *model.CommandArgs, token string, config *store.Config) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetAllPipelinesForProject(token, config.ToSlug())
	if err != nil {
		p.API.LogError("Failed to fetch info for pipeline", "project", config.ToSlug(), "error", err)
		return p.sendEphemeralResponse(args, "Failed to fetch info for pipeline"), nil
	}

	projectsListString := "| Pipeline No. | Pipeline ID | State |\n| :---- | :----- | \n"
	for _, pipeline := range pipelines {
		projectsListString += fmt.Sprintf(
			"| %d | %s | %s |\n",
			pipeline.Number,
			pipeline.Id,
			pipeline.State,
		)
	}

	_ = p.sendEphemeralPost(
		args,
		fmt.Sprintf("Recently built pipelines for project %s.", config.ToMarkdown()),
		[]*model.SlackAttachment{
			{
				Fallback: "Pipelines list",
				Text:     projectsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executePipelineGetAllForProjectByMe(args *model.CommandArgs, token string, config *store.Config) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetAllMyPipelinesForProject(token, config.ToSlug())
	if err != nil {
		p.API.LogError("Failed to fetch info for pipeline", "project", config.ToMarkdown(), "error", err.Error())
		return p.sendEphemeralResponse(args, "Failed to fetch info for pipeline"), nil
	}

	projectsListString := "| Pipeline No. | Pipeline ID | State |\n| :---- | :----- | \n"
	for _, pipeline := range pipelines {
		projectsListString += fmt.Sprintf(
			"| %d | %s | %s |\n",
			pipeline.Number,
			pipeline.Id,
			pipeline.State,
		)
	}

	_ = p.sendEphemeralPost(
		args,
		fmt.Sprintf("Pipelines recently ran by you for project %s", config.ToMarkdown()),
		[]*model.SlackAttachment{
			{
				Fallback: "Pipelines list",
				Text:     projectsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executePipelineGetWorkflowByID(args *model.CommandArgs,
	token string, pipelineID string) (*model.CommandResponse, *model.AppError) {
	if pipelineID == "" {
		return p.sendEphemeralResponse(args, "Please provide the pipelineID"), nil
	}

	wfs, err := circle.GetWorkflowsByPipeline(token, pipelineID)
	if err != nil {
		p.API.LogError("Failed to fetch wokflows for given pipeline", "pipelineID", pipelineID, "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Failed to fetch workflows for pipeline `%s`", pipelineID),
		), nil
	}

	workflowListString := "| Name | Started By | Status | ID |\n| :---- | :----- | \n"
	for _, wf := range wfs {
		uname, err := circle.GetNameByID(token, wf.StartedBy)
		if err != nil {
			uname = wf.StartedBy
		}
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
		fmt.Sprintf("Workflows for given pipeline ID: `%s`", pipelineID),
		[]*model.SlackAttachment{
			{
				Fallback: "Workflow List",
				Text:     workflowListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeTriggerPipeline(args *model.CommandArgs, token string,
	config *store.Config, branch string) (*model.CommandResponse, *model.AppError) {
	pl, err := circle.TriggerPipeline(token, config.ToSlug(), branch)
	if branch == "" {
		branch = "master"
	}
	if err != nil {
		p.API.LogError("Could not trigger pipeline", "project", config.ToSlug(), "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("Could not trigger pipeline for project %s on `%s` branch", config.ToSlug(), branch),
		), nil
	}

	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				Fallback: fmt.Sprintf("Pipeline triggered successfully for project %s for branch: %s", config.ToMarkdown(), branch),
				Pretext:  fmt.Sprintf("Triggered pipeline for project %s branch `%s`", config.ToMarkdown(), branch),
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

func (p *Plugin) executePipelineGetSingle(args *model.CommandArgs, token string,
	config *store.Config, num string) (*model.CommandResponse, *model.AppError) {
	var isUUID bool
	var err error
	var pl circleci.Pipeline
	_, err = strconv.ParseInt(num, 10, 64)
	if err != nil {
		isUUID = true
	}
	if !isUUID && config.ToSlug() == "" {
		return p.sendEphemeralResponse(args,
			"Please provide project slug via --project flag. i.e. --project vcs/org/repo or configure default project"), nil
	}

	if isUUID {
		pl, err = circle.GetPipelineByID(token, num)
	} else {
		pl, err = circle.GetPipelineByNum(token, config.ToSlug(), num)
	}

	if err != nil {
		p.API.LogError("Could not get info about pipeline", "pipelineNumber", num, "error", err)
		return p.sendEphemeralResponse(args, "Could not get info about pipeline"), nil
	}

	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				Fallback: "Pipeline Info",
				Pretext:  fmt.Sprintf("Informations about pipeline `%s`", num),
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Number",
						Value: num,
						Short: true,
					},
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
