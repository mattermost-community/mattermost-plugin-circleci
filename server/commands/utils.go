package commands

import (
	"github.com/TomTucka/go-circleci"
)

// Return Markdown text with the formatted status, or a badge if we have it
func buildStatusToMarkdown(build *circleci.Build, p *Plugin) string {
	switch build.Status {
	case "running":
		return "Running"

	case "not_run":
		return "Not Run"

	case "canceled":
		return "Canceled"

	case "failing":
		return "Failing"

	case "failed":
		return "![Status image](" + badgeFailedURL + ")"

	case "success":
		return "![Status image](" + badgePassedURL + ")"

	case "on_hold":
		return "On Hold"

	case "needs_setup":
		return "Need Setup"

	default:
		return build.Status
	}
}

// Return the formatted start time
func buildStartTimeToString(build *circleci.Build) string {
	buildStartTime := "Not Run"
	if build.StartTime != nil {
		buildStartTime = build.StartTime.Format("15:04:05 on 2006-01-02") // TODO Improve format
	}

	return buildStartTime
}
