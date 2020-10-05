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
	pipelineTriggerHint     = "<" + branchTrigger + "|" + tagTrigger + ">"
	pipelineTriggerHelpText = "Trigger pipeline for a project"

	branchTrigger         = "branch"
	branchTriggerHint     = "<branch name>"
	branchTriggerHelpText = "Provide branch name for which you want to trigger the pipeline"

	tagTrigger         = "tag"
	tagTriggerHint     = "<tag value>"
	tagTriggerHelpText = "Provide tag value for which you want to trigger the pipeline"

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
	branch := model.NewAutocompleteData(branchTrigger, branchTriggerHint, branchTriggerHelpText)
	branch.AddTextArgument("<branch>", "The branch for which pipeline will be trigeered. Leave empty for master", "")
	branch.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)
	tag := model.NewAutocompleteData(tagTrigger, tagTriggerHint, tagTriggerHelpText)
	tag.AddTextArgument("<tag>", "The tag for which pipeline will be trigeered.", "")
	tag.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)
	trigger.AddCommand(branch)
	trigger.AddCommand(tag)

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

func (p *Plugin) executePipelineTrigger(args *model.CommandArgs, circleciToken string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
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
		return p.executePipelineGetAllForProject(args, circleciToken, project)
	case pipelineGetRecentTrigger:
		return p.executePipelineGetRecent(args, circleciToken, argument)
	case pipelineGetMineTrigger:
		return p.executePipelineGetAllForProjectByMe(args, circleciToken, project)
	case pipelineWorkflowTrigger:
		return p.executePipelineGetWorkflowByID(args, circleciToken, argument)
	case pipelineTriggerTrigger:
		return p.executeTriggerPipeline(args, circleciToken, project, split[1:])
	case pipelineGetSingleTrigger:
		return p.executePipelineGetSingle(args, circleciToken, project, argument)

	case commandHelpTrigger:
		return p.sendHelpResponse(args, pipelineTrigger)
	default:
		return p.sendIncorrectSubcommandResponse(args, pipelineTrigger)
	}
}

func (p *Plugin) executePipelineGetRecent(args *model.CommandArgs, token string,
	orgSlug string) (*model.CommandResponse, *model.AppError) {
	if orgSlug == "" {
		return p.sendEphemeralResponse(args, "Please provide org slug in the form of vcs/orgname."), nil
	}
	pipelines, err := circle.GetRecentlyBuiltPipelines(token, orgSlug, false)
	if err != nil {
		p.API.LogError("Failed to fetch info for pipeline", "org", orgSlug, "error", err.Error())
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

func (p *Plugin) executePipelineGetAllForProject(args *model.CommandArgs, token string, project *store.ProjectIdentifier) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetAllPipelinesForProject(token, project.ToSlug())
	if err != nil {
		p.API.LogError("Failed to fetch info for pipeline", "project", project.ToSlug(), "error", err)
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
		fmt.Sprintf("Recently built pipelines for project %s.", project.ToMarkdown()),
		[]*model.SlackAttachment{
			{
				Fallback: "Pipelines list",
				Text:     projectsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executePipelineGetAllForProjectByMe(args *model.CommandArgs, token string, project *store.ProjectIdentifier) (*model.CommandResponse, *model.AppError) {
	pipelines, err := circle.GetAllMyPipelinesForProject(token, project.ToSlug())
	if err != nil {
		p.API.LogError("Failed to fetch info for pipeline", "project", project.ToMarkdown(), "error", err.Error())
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
		fmt.Sprintf("Pipelines recently ran by you for project %s", project.ToMarkdown()),
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
	project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
	var params circleci.TriggerPipelineParameters
	subcmd := "branch"
	if len(split) > 0 {
		subcmd = split[0]
	}
	input := ""

	switch subcmd {
	case branchTrigger:
		branch := "master"
		if len(split) > 1 {
			branch = split[1]
		}
		params = circleci.TriggerPipelineParameters{Branch: branch}
		input = branch

	case tagTrigger:
		if len(split) < 2 {
			return p.sendEphemeralResponse(args, ":red_circle: Please provide a tag value."), nil
		}
		input = split[1]
		params = circleci.TriggerPipelineParameters{Tag: input}
	}

	pl, err := circle.TriggerPipeline(token, project.ToSlug(), params)
	if err != nil {
		p.API.LogError("Could not trigger pipeline", "project", project.ToSlug(), "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Could not trigger pipeline for project %s on %s: `%s` ", project.ToMarkdown(), subcmd, input),
		), nil
	}

	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				Fallback: fmt.Sprintf("Pipeline triggered successfully for project %s for %s: %s", project.ToMarkdown(), subcmd, input),
				Pretext:  fmt.Sprintf(":white_check_mark: Triggered pipeline for project %s, %s: `%s`", project.ToMarkdown(), subcmd, input),
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Number",
						Value: strconv.FormatInt(pl.Number, 10),
						Short: true,
					},
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
	project *store.ProjectIdentifier, num string) (*model.CommandResponse, *model.AppError) {
	var isUUID bool
	var err error
	var pl circleci.Pipeline
	_, err = strconv.ParseInt(num, 10, 64)
	if err != nil {
		isUUID = true
	}
	if !isUUID && project.ToSlug() == "" {
		return p.sendEphemeralResponse(args,
			"Please provide project slug via `--project` flag. i.e. `--project <vcs/org-name/project-name` or configure default project"), nil
	}

	if isUUID {
		pl, err = circle.GetPipelineByID(token, num)
	} else {
		pl, err = circle.GetPipelineByNum(token, project.ToSlug(), num)
	}

	if err != nil {
		p.API.LogError("Could not get info about pipeline", "pipelineNumber", num, "error", err)
		return p.sendEphemeralResponse(args, ":red_circle: Could not get info about pipeline"), nil
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
						Value: strconv.FormatInt(pl.Number, 10),
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
