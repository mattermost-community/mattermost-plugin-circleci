package plugin

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

const (
	insightTrigger  = "insight"
	insightHint     = "<" + insightMetricsWorkflowTrigger + "|" + insightMetricsWorkflowJobsTrigger + ">"
	insightHelpText = "Get insights about"

	insightMetricsWorkflowTrigger  = "workflows"
	insightMetricsWorkflowHint     = "<vcs-slug/org-name/repo-name>"
	insightMetricsWorkflowHelpText = "Get summary metrics for a project's workflows"

	insightMetricsWorkflowJobsTrigger  = "jobs"
	insightMetricsWorkflowJobsHint     = "<vcs-slug/org-name/repo-name> <workflow name>"
	insightMetricsWorkflowJobsHelpText = "Get summary metrics for a project workflow's jobs"
)

func getInsightAutoCompeleteData() *model.AutocompleteData {
	insight := model.NewAutocompleteData(insightTrigger, insightHint, insightHelpText)
	wf := model.NewAutocompleteData(insightMetricsWorkflowTrigger, insightMetricsWorkflowHint, insightMetricsWorkflowHelpText)
	wf.AddTextArgument("<vcs-slug/org-name/repo-name>", "Project to get workflows metrics summary of. ex: gh/mattermost/mattermost-server", "")
	jb := model.NewAutocompleteData(insightMetricsWorkflowTrigger, insightMetricsWorkflowHint, insightMetricsWorkflowHelpText)
	jb.AddTextArgument("<vcs-slug/org-name/repo-name>", "Project to get metrics summary of. ex: gh/mattermost/mattermost-server", "")
	jb.AddTextArgument("<workflow name", "Name of workflow to get metrics. ex: worfkflow-test", "")
	insight.AddCommand(wf)
	insight.AddCommand(jb)
	return insight
}

func (p *Plugin) executeInsightTrigger(args *model.CommandArgs, circleciToken string,
	split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case insightMetricsWorkflowTrigger:
		return p.executeInsightWorkflowMetrics(args, circleciToken, split[1:])
	case insightMetricsWorkflowJobsTrigger:
		return p.executeInsightJobMetrics(args, circleciToken, split[1:])
	default:
		return p.sendIncorrectSubcommandResponse(args, pipelineTrigger)
	}
}

func (p *Plugin) executeInsightWorkflowMetrics(args *model.CommandArgs,
	token string, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 1 {
		return p.sendEphemeralResponse(args, "Please provide project slug to get workflow metrics"), nil
	}
	wfm, err := circle.GetWorkflowMetrics(token, split[0])
	if err != nil {
		return p.sendEphemeralResponse(args, fmt.Sprintf("Could not get workflow metrics for project %s", split[0])),
			&model.AppError{Message: "Failed to get workflow metrics for project " + split[0]}
	}
	wfMetricsString := "| Name | Sucess Rate | Failed Runs | Successful Runs | Throughput" +
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
		"Workflow metrics for project: "+split[0],
		[]*model.SlackAttachment{
			{
				Fallback: "Workflow Metrics",
				Text:     wfMetricsString,
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeInsightJobMetrics(args *model.CommandArgs,
	token string, split []string) (*model.CommandResponse, *model.AppError) {

	return &model.CommandResponse{}, nil
}
