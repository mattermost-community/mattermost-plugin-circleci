package plugin

import (
	"github.com/mattermost/mattermost-server/v5/model"
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

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeInsightJobMetrics(args *model.CommandArgs,
	token string, split []string) (*model.CommandResponse, *model.AppError) {

	return &model.CommandResponse{}, nil
}
