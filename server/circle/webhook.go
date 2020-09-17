package circle

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

type BuildInfos struct {
	Owner          string `json:"Owner"`
	Repository     string `json:"Repository"`
	CircleBuildNum int    `json:"CircleBuildNum"`
	Failed         bool   `json:"Failed"`
	Message        string `json:"Message"`
}

// Convert the build info into a post attachment
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
