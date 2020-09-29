package utils

import (
	"fmt"

	v1 "github.com/jszwedko/go-circleci"
)

// BuildStatusToMarkdown returns Markdown text with the formatted status, or a badge if we have it
func BuildStatusToMarkdown(build *v1.Build, badgeFailedURL string, badgePassedURL string) string {
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

// BuildStartTimeToString returns the formatted start time
func BuildStartTimeToString(build *v1.Build) string {
	buildStartTime := "Not Run"
	if build.StartTime != nil {
		buildStartTime = build.StartTime.Format("15:04:05 on 2006-01-02") // TODO Improve format
	}

	return buildStartTime
}

// CircleciUserToString return in the format "FullName (username)"
func CircleciUserToString(user *v1.User) string {
	if *user.Name != "" {
		return *user.Name + " (" + user.Login + ")"
	}
	return user.Login
}

// GetFullNameFromOwnerAndRepo get full name
func GetFullNameFromOwnerAndRepo(owner string, repository string) string {
	return fmt.Sprintf("%s/%s", owner, repository)
}
