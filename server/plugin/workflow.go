package plugin

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

const (
	workflowTrigger = "workflow"
	workflowHint    = "<" + workflowGetJobsTrigger + "|" + workflowGetTrigger + "|" + workflowRerunTrigger +
		"|" + workflowCancelTrigger + ">"
	workflowHelpText = "Manage the connection to your CircleCI acccount"

	workflowGetTrigger  = "get"
	workflowGetHint     = "<workflowID>"
	workflowGetHelpText = "Get informations about workflow"

	workflowGetJobsTrigger  = "jobs"
	workflowGetJobsHint     = "<workflowID>"
	workflowGetJobsHelpText = "Get jobs list of workflow"

	workflowRerunTrigger  = "rerun"
	workflowRerunHint     = "<workflowID>"
	workflowRerunHelpText = "Rerun a workflow"

	workflowCancelTrigger  = "cancel"
	workflowCancelHint     = "<workflowID>"
	workflowCancelHelpText = "Cancel a workflow"
)

func getWorkflowAutoCompeleteData() *model.AutocompleteData {
	workflow := model.NewAutocompleteData(workflowTrigger, workflowHint, workflowHelpText)

	workflowGet := model.NewAutocompleteData(workflowGetTrigger, workflowGetHint, workflowGetHelpText)
	workflowGet.AddTextArgument("<workflowID>", workflowGetHint, "")

	workflowGetJobs := model.NewAutocompleteData(workflowGetJobsTrigger, workflowGetJobsHint, workflowGetJobsHelpText)
	workflowGetJobs.AddTextArgument("<workflowID>", workflowGetJobsHint, "")

	rerun := model.NewAutocompleteData(workflowRerunTrigger, workflowRerunHint, workflowRerunHelpText)
	rerun.AddTextArgument("<workflowID>", workflowRerunHint, "")

	cancel := model.NewAutocompleteData(workflowCancelTrigger, workflowCancelHint, workflowCancelHelpText)
	cancel.AddTextArgument("<workflowID>", workflowCancelHint, "")

	workflow.AddCommand(rerun)
	workflow.AddCommand(cancel)
	workflow.AddCommand(workflowGet)
	workflow.AddCommand(workflowGetJobs)
	return workflow
}

func (p *Plugin) executeWorkflow(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := commandHelpTrigger
	if len(split) > 0 {
		subcommand = split[0]
	}

	var workflow string
	if len(split) > 1 {
		workflow = split[1]
	} else if subcommand != commandHelpTrigger {
		return p.sendEphemeralResponse(args, "Please precise the ID of the workflow."), nil
	}

	switch subcommand {
	case workflowGetTrigger:
		return p.executeWorflowGet(args, circleciToken, workflow)
	case workflowGetJobsTrigger:
		return p.executeWorflowGetJobs(args, circleciToken, workflow)
	case workflowRerunTrigger:
		return p.executeRerunWorkflow(args, circleciToken, workflow)
	case workflowCancelTrigger:
		return p.executeCancelWorkflow(args, circleciToken, workflow)

	case commandHelpTrigger:
		return p.sendHelpResponse(args, workflowTrigger)
	default:
		return p.sendIncorrectSubcommandResponse(args, workflowTrigger)
	}
}

func (p *Plugin) executeWorflowGet(args *model.CommandArgs, token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	wf, err := circle.GetWorkflow(token, workflowID)
	if err != nil {
		p.API.LogError("Failed to fetch info for workflow", "workflowID", workflowID, "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Failed to fetch info for workflow `%s`", workflowID),
		), nil
	}

	uname := wf.StartedBy
	tmp, err := circle.GetNameByID(token, wf.StartedBy)
	if err == nil {
		uname = tmp
	}

	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				Fallback: fmt.Sprintf("Informations for workflow %s", wf.Name),
				Pretext:  fmt.Sprintf("Informations for worflow `%s`", wf.Id),
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Name",
						Value: wf.Name,
						Short: true,
					},
					{
						Title: "ID",
						Value: wf.Id,
						Short: true,
					},
					{
						Title: "Project slug",
						Value: wf.ProjectSlug,
						Short: true,
					},
					{
						Title: "Pipeline ID",
						Value: wf.PipelineId,
						Short: true,
					},
					{
						Title: "Pipeline Number",
						Value: string(wf.PipelineNumber),
						Short: true,
					},
					{
						Title: "Status",
						Value: wf.Status,
						Short: true,
					},
					{
						Title: "Created At",
						Value: wf.CreatedAt.String(),
						Short: true,
					},
					{
						Title: "Stopped At",
						Value: wf.StoppedAt.String(),
						Short: true,
					},
					{
						Title: "Started By",
						Value: uname,
						Short: true,
					},
				},
			},
		},
	)
	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeWorflowGetJobs(args *model.CommandArgs, token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	_, err := circle.GetWorkflow(token, workflowID)
	if err != nil {
		p.API.LogError("Failed to fetch info for workflow", "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Failed to fetch informations for workflow %s", workflowID),
		), nil
	}

	jobs, errr := circle.GetWorkflowJobs(token, workflowID)
	if errr != nil {
		p.API.LogError("Failed to fetch jobs informations for workflow", "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Failed to fetch jobs informations for workflow %s", workflowID),
		), nil
	}

	workflowJobsListString := "| Name | Type | Status | Project | ID |\n| :----- | :----- | :----- | :----- | :----- | \n"
	for _, job := range *jobs {
		workflowJobsListString += fmt.Sprintf(
			"| %s | %s | %s | %s | %s |\n",
			job.Name,
			job.Type_,
			*job.Status,
			job.ProjectSlug,
			job.Id,
		)
	}

	_ = p.sendEphemeralPost(
		args,
		fmt.Sprintf("Jobs for given workflow ID `%s`", workflowID),
		[]*model.SlackAttachment{
			{
				Fallback: "Workflow Jobs List",
				Text:     workflowJobsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeRerunWorkflow(args *model.CommandArgs, token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	_, err := circle.RerunWorkflow(token, workflowID)
	if err != nil {
		p.API.LogError("Failed to re run workflow", "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":information_source: Failed to re run workflow %s", workflowID),
		), nil
	}

	wf, err := circle.GetWorkflow(token, workflowID)
	var msg string
	if err != nil {
		msg = fmt.Sprintf(":information_source: Re running workflow: workflow ID: %s", workflowID)
	} else {
		msg = fmt.Sprintf(":information_source: Re running workflow: %s, workflow ID: %s", wf.Name, wf.Id)
	}

	return p.sendEphemeralResponse(args, msg), nil
}

func (p *Plugin) executeCancelWorkflow(args *model.CommandArgs, token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	_, err := circle.CancelWorkflow(token, workflowID)
	if err != nil {
		p.API.LogError("Failed to cancel workflow", "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":information_source: Failed to cancel workflow %s", workflowID),
		), nil
	}

	wf, err := circle.GetWorkflow(token, workflowID)
	var msg string
	if err != nil {
		msg = fmt.Sprintf(":information_source: Canceled workflow. Workflow ID: %s", workflowID)
	} else {
		msg = fmt.Sprintf(":information_source: Canceled workflow: %s, workflow ID: %s", wf.Name, wf.Id)
	}

	return p.sendEphemeralResponse(args, msg), nil
}
