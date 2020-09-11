package main

import "github.com/jszwedko/go-circleci"

// Return false if an error has occurred
func (p *Plugin) getCircleUserInfo(circleToken string) (*circleci.User, bool) {
	circleClient := &circleci.Client{
		Token: circleToken,
	}

	user, err := circleClient.Me()
	if err != nil {
		p.API.LogError("Error when reaching CircleCI", "CircleCI error:", err)
		return nil, false
	}

	return user, true
}

// Return false if an error has occurred
func (p *Plugin) getCircleciUserProjects(circleCiToken string) ([]*circleci.Project, bool) {
	circleciClient := &circleci.Client{Token: circleCiToken}
	projects, err := circleciClient.ListProjects()
	if err != nil {
		p.API.LogError("Unable to get circleCI user projects", "circleCI API error", err)
		return nil, false
	}

	return projects, true
}

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

// Return in the format "FullName (username)"
func circleciUserToString(user *circleci.User) string {
	if *user.Name != "" {
		return *user.Name + " (" + user.Login + ")"
	}
	return user.Login
}
