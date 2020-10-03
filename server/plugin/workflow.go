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
	workflowGetHint     = "< workflow_id >"
	workflowGetHelpText = "Get informations about workflow"

	workflowGetJobsTrigger         = "jobs"
	workflowGetJobsHint            = "< workflow_id >"
	workflowGetJobsTriggerHelpText = "Get jobs list of workflow"

	workflowRerunTrigger  = "rerun"
	workflowRerunHint     = "< workflow ID >"
	workflowRerunHelpText = "Rerun a workflow"

	workflowCancelTrigger  = "cancel"
	workflowCancelHint     = "< workflow ID >"
	workflowCancelHelpText = "Cancel a workflow"
)

func getWorkflowAutoCompeleteData() *model.AutocompleteData {
	workflow := model.NewAutocompleteData(workflowTrigger, workflowHint, workflowHelpText)
	workflowGet := model.NewAutocompleteData(workflowGetTrigger, workflowGetHint, workflowGetHelpText)
	workflowGet.AddTextArgument("workflow id", workflowGetHint, "")
	workflowGetJobs := model.NewAutocompleteData(workflowGetJobsTrigger, workflowGetJobsHint, workflowGetJobsTriggerHelpText)
	workflowGetJobs.AddTextArgument("workflow id", workflowGetJobsHint, "")
	rerun := model.NewAutocompleteData(workflowRerunTrigger, workflowRerunHint, workflowRerunHelpText)
	rerun.AddTextArgument("workflow id", workflowRerunHint, "")
	cancel := model.NewAutocompleteData(workflowCancelTrigger, workflowCancelHint, workflowCancelHelpText)
	cancel.AddTextArgument("workflow id", workflowCancelHint, "")
	workflow.AddCommand(workflowGet)
	workflow.AddCommand(workflowGetJobs)
	workflow.AddCommand(rerun)
	workflow.AddCommand(cancel)
	return workflow
}

func (p *Plugin) executeWorkflowTrigger(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	var workflow string
	if len(split) > 1 {
		workflow = split[1]
	} else {
		return p.sendIncorrectSubcommandResponse(args, workflowTrigger)
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
	default:
		return p.sendIncorrectSubcommandResponse(args, workflowTrigger)
	}
}

func (p *Plugin) executeWorflowGet(args *model.CommandArgs,
	token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	wf, err := circle.GetWorkflow(token, workflowID)
	if err != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Failed to fetch info for workflow", workflowID, err.Error())}
	}
	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				Fallback: "Workflow Name: " + wf.Name,
				Pretext:  "Information for worflow Id " + wf.Id,
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
						Value: wf.PipelineNumber,
						Short: true,
					},
					{
						Title: "Status",
						Value: wf.Status,
						Short: true,
					},
					{
						Title: "Created At",
						Value: wf.CreatedAt,
						Short: true,
					},
					{
						Title: "Stopped At",
						Value: wf.StoppedAt,
						Short: true,
					},
					{
						Title: "Started By",
						Value: wf.StartedBy,
						Short: true,
					},
				},
			},
		},
	)
	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeWorflowGetJobs(args *model.CommandArgs,
	token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	_, err := circle.GetWorkflow(token, workflowID)
	if err != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Failed to fetch info for workflow", workflowID, err.Error())}
	}
	jobs, errr := circle.GetWorkflowJobs(token, workflowID)
	if errr != nil {
		return nil, &model.AppError{Message: fmt.Sprintf("%s%s. err %s",
			"Failed to fetch jobs info for workflow", workflowID, errr.Error())}
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
		"Jobs for given workflow ID: "+workflowID,
		[]*model.SlackAttachment{
			{
				Fallback: "Workflow Jobs List",
				Text:     workflowJobsListString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeRerunWorkflow(args *model.CommandArgs,
	token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	_, err := circle.RerunWorkflow(token, workflowID)
	if err != nil {
		return p.sendEphemeralResponse(args, fmt.Sprintf("Could not re-run workflow. workflow ID: %s", workflowID)),
			&model.AppError{Message: fmt.Sprintf("%s%s. err %s", "Failed to re run workflow", workflowID, err.Error())}
	}
	wf, err := circle.GetWorkflow(token, workflowID)
	var errstr string
	if err != nil {
		errstr = fmt.Sprintf("Re running workflow: workflow ID: %s", workflowID)
	} else {
		errstr = fmt.Sprintf("Re running workflow: %s, workflow ID: %s", wf.Name, wf.Id)
	}
	return p.sendEphemeralResponse(args, errstr), nil
}

func (p *Plugin) executeCancelWorkflow(args *model.CommandArgs,
	token string, workflowID string) (*model.CommandResponse, *model.AppError) {
	_, err := circle.CancelWorkflow(token, workflowID)
	if err != nil {
		return p.sendEphemeralResponse(args, fmt.Sprintf("Could not cancel workflow. workflow ID: %s", workflowID)),
			&model.AppError{Message: fmt.Sprintf("%s%s. err %s", "Failed to cancel workflow", workflowID, err.Error())}
	}
	wf, err := circle.GetWorkflow(token, workflowID)
	var errstr string
	if err != nil {
		errstr = "Canceled workflow. workflow ID: " + workflowID
	} else {
		errstr = fmt.Sprintf("Canceled workflow: %s, workflow ID: %s", wf.Name, wf.Id)
	}
	return p.sendEphemeralResponse(args, errstr), nil
}
