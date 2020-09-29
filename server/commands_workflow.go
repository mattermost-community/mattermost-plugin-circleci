package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

const (
	workflowTrigger  = "workflow"
	workflowHint     = "<" + workflowGetTrigger + ">"
	workflowHelpText = "Manage the connection to your CircleCI acccount"

	workflowGetTrigger  = "get"
	workflowGetHint     = "< workflow_id >"
	workflowGetHelpText = "Get informations about workflow"

	workflowGetJobsTrigger         = "jobs"
	workflowGetJobsHint            = "< workflow_id >"
	workflowGetJobsTriggerHelpText = "Get jobs list of workflow"
)

func (p *Plugin) executeWorkflowTrigger(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case workflowGetTrigger:
		if len(split) > 1 {
			if split[1] == workflowGetJobsTrigger {
				return p.executeWorflowGetJobs(args, circleciToken, split[1])
			}
			return p.executeWorflowGet(args, circleciToken, split[1])
		}
		return p.sendIncorrectSubcommandResponse(args, workflowTrigger)
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
			"Failed to fetch jobs info for workflow", workflowID, err.Error())}
	}

	_ = jobs
	// TODO: return formatted jobs list

	// _ = p.sendEphemeralPost(
	// 	args,
	// 	"",
	// 	[]*model.SlackAttachment{
	// 		{
	// 			Fallback: "Workflow Name: " + wf.Name,
	// 			Pretext:  "Jobs information for worflow Id " + wf.Id,
	// 			Fields: []*model.SlackAttachmentField{
	// 				{
	// 					Title: "Name",
	// 					Value: wf.Name,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "ID",
	// 					Value: wf.Id,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Project slug",
	// 					Value: wf.ProjectSlug,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Pipeline ID",
	// 					Value: wf.PipelineId,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Pipeline Number",
	// 					Value: wf.PipelineNumber,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Status",
	// 					Value: wf.Status,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Created At",
	// 					Value: wf.CreatedAt,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Stopped At",
	// 					Value: wf.StoppedAt,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Started By",
	// 					Value: wf.StartedBy,
	// 					Short: true,
	// 				},
	// 			},
	// 		},
	// 	},
	// )
	return &model.CommandResponse{}, nil
}
