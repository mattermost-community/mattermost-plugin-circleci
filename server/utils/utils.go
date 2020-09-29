package utils

import (
	"fmt"

	v1 "github.com/jszwedko/go-circleci"
	"github.com/mattermost/mattermost-server/model"
)

// BuildStatusToMarkdown returns Markdown text with the formatted status, or a badge if we have it
func BuildStatusToMarkdown(build *v1.Build) string {
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

// BuildInfos ..
type BuildInfos struct {
	Owner          string `json:"Owner"`
	Repository     string `json:"Repository"`
	CircleBuildNum int    `json:"CircleBuildNum"`
	Failed         bool   `json:"Failed"`
	Message        string `json:"Message"`
}

// ToPostAttachments converts the build info into a post attachment
func (bi *BuildInfos) ToPostAttachments(buildFailedIconURL, buildGreenIconURL string) []*model.SlackAttachment {
	// TODO add link to build
	attachment := &model.SlackAttachment{
		Fields: []*model.SlackAttachmentField{
			{
				Title: "Repo",
				Short: true,
				Value: GetFullNameFromOwnerAndRepo(bi.Owner, bi.Repository),
			},
			{
				Title: "Job number",
				Short: true,
				Value: fmt.Sprintf("%d", bi.CircleBuildNum),
			},
		},
	}

	if bi.Message != "" {
		attachment.Fields = append(attachment.Fields,
			&model.SlackAttachmentField{
				Title: "Message",
				Short: false,
				Value: fmt.Sprintf("```\n%s\n```", bi.Message),
			},
		)
	}

	if bi.Failed {
		attachment.ThumbURL = buildFailedIconURL
		attachment.Title = "Job failed"
		attachment.Color = "#FF1919" // red
	} else {
		attachment.ThumbURL = buildGreenIconURL
		attachment.Title = "Job passed"
		attachment.Color = "#50F100" // green
	}

	attachment.Fallback = attachment.Title
	return []*model.SlackAttachment{
		attachment,
	}
}
