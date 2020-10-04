package plugin

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	insightTrigger  = "insight"
	insightHint     = "<" + insightMetricsWorkflowTrigger + "|" + insightMetricsWorkflowJobsTrigger + ">"
	insightHelpText = "Get insights about"

	insightMetricsWorkflowTrigger  = "workflows"
	insightMetricsWorkflowHint     = ""
	insightMetricsWorkflowHelpText = "Get summary metrics for a project's workflows"

	insightMetricsWorkflowJobsTrigger  = "jobs"
	insightMetricsWorkflowJobsHint     = "<workflow name>"
	insightMetricsWorkflowJobsHelpText = "Get summary metrics for a project workflow's jobs"
)

func getInsightAutoCompeleteData() *model.AutocompleteData {
	insight := model.NewAutocompleteData(insightTrigger, insightHint, insightHelpText)

	wf := model.NewAutocompleteData(insightMetricsWorkflowTrigger, insightMetricsWorkflowHint, insightMetricsWorkflowHelpText)
	wf.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	jb := model.NewAutocompleteData(insightMetricsWorkflowJobsTrigger, insightMetricsWorkflowJobsHint, insightMetricsWorkflowJobsHelpText)
	jb.AddTextArgument("<workflow name>", "Name of workflow to get metrics from", "")
	jb.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	insight.AddCommand(wf)
	insight.AddCommand(jb)

	return insight
}

func (p *Plugin) executeInsightTrigger(args *model.CommandArgs, circleciToken string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case insightMetricsWorkflowTrigger:
		return p.executeInsightWorkflowMetrics(args, circleciToken, project)
	case insightMetricsWorkflowJobsTrigger:
		return p.executeInsightJobMetrics(args, circleciToken, project, split[1:])
	default:
		return p.sendIncorrectSubcommandResponse(args, pipelineTrigger)
	}
}

func (p *Plugin) executeInsightWorkflowMetrics(args *model.CommandArgs, token string, project *store.ProjectIdentifier) (*model.CommandResponse, *model.AppError) {
	wfm, err := circle.GetWorkflowMetrics(token, project.ToSlug())
	if err != nil {
		p.API.LogError("Failed to get workflow metrics", "project", project.ToSlug(), "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Could not get workflow metrics for project %s", project.ToMarkdown()),
		), nil
	}

	wfMetricsString := "| Name | Success Rate | Failed Runs | Successful Runs | Throughput" +
		"| MTTR | Credits Used | Mean | Median | Min | Max | Time Widnow |\n| :---- | :----- | :---- |\n"
	for _, wf := range wfm {
		mean := float32(wf.Metrics.DurationMetrics.Mean / 3600)
		median := float32(wf.Metrics.DurationMetrics.Median / 3600)
		min := float32(wf.Metrics.DurationMetrics.Min / 3600)
		max := float32(wf.Metrics.DurationMetrics.Max / 3600)
		mttr := wf.Metrics.Mttr / 3600
		wfMetricsString += fmt.Sprintf(
			"| %s | %f | %d | %d | %f | %d | %d | %f | %f | %f | %f | %s |\n",
			wf.Name, wf.Metrics.SuccessRate*100, wf.Metrics.FailedRuns,
			wf.Metrics.SuccessfulRuns, wf.Metrics.Throughput, mttr,
			wf.Metrics.TotalCreditsUsed, mean, median, min, max,
			fmt.Sprintf("%s to %s", wf.WindowStart.Format("2006-01-02"), wf.WindowEnd.Format("2006-01-02")),
		)
	}

	_ = p.sendEphemeralPost(
		args,
		fmt.Sprintf("Workflow metrics for project %s ", project.ToMarkdown()),
		[]*model.SlackAttachment{
			{
				Fallback: "Workflow Metrics",
				Text:     wfMetricsString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeInsightJobMetrics(args *model.CommandArgs, token string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 1 {
		return p.sendEphemeralResponse(args, "Please provide the workflow name to get jobs metrics"), nil
	}

	workflowName := split[0]

	wfm, err := circle.GetWorkflowJobsMetrics(token, project.ToSlug(), workflowName)
	if err != nil {
		p.API.LogError("Failed to get jobs metrics", "project", project.ToSlug(), "workflow", workflowName, "error", err)
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: Could not get job metrics for project %s, workflow `%s`", project.ToMarkdown(), workflowName),
		), nil
	}

	if len(wfm) == 0 {
		return p.sendEphemeralResponse(args,
			fmt.Sprintf("There is no metric for project: %s — workflow: `%s` ", project.ToMarkdown(), workflowName),
		), nil
	}

	wfMetricsString := "| Name | Success Rate | Failed Runs | Successful Runs | Throughput" +
		"| Credits Used | Mean | Median | Min | Max | Time Widnow |\n| :---- | :----- | :---- |\n"
	for _, wf := range wfm {
		mean := float32(wf.Metrics.DurationMetrics.Mean / 3600)
		median := float32(wf.Metrics.DurationMetrics.Median / 3600)
		min := float32(wf.Metrics.DurationMetrics.Min / 3600)
		max := float32(wf.Metrics.DurationMetrics.Max / 3600)
		wfMetricsString += fmt.Sprintf(
			"| %s | %f | %d | %d | %f | %d | %f | %f | %f | %f | %s |\n",
			wf.Name, wf.Metrics.SuccessRate*100, wf.Metrics.FailedRuns,
			wf.Metrics.SuccessfulRuns, wf.Metrics.Throughput,
			wf.Metrics.TotalCreditsUsed, mean, median, min, max,
			fmt.Sprintf("%s to %s", wf.WindowStart.Format("2006-01-02"), wf.WindowEnd.Format("2006-01-02")),
		)
	}

	_ = p.sendEphemeralPost(
		args,
		fmt.Sprintf("Job metrics for project: %s — workflow: `%s`", project.ToMarkdown(), workflowName),
		[]*model.SlackAttachment{
			{
				Fallback: "Job Metrics",
				Text:     wfMetricsString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}
