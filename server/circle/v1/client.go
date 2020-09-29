package v1

import (
	"fmt"

	"github.com/jszwedko/go-circleci"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/utils"
)

var (
	badgeFailedURL string
	badgePassedURL string
)

// GetCircleUserInfo returns info about logged in user
func GetCircleUserInfo(circleToken string) (*circleci.User, error) {
	circleClient := &circleci.Client{
		Token: circleToken,
	}

	user, err := circleClient.Me()
	if err != nil {
		return nil, fmt.Errorf("Error when reaching CircleCI. CircleCI error: %s", err.Error())
	}

	return user, nil
}

// GetCircleciUserProjects returns projects for given user
func GetCircleciUserProjects(circleCiToken string) ([]*circleci.Project, error) {
	circleciClient := &circleci.Client{Token: circleCiToken}
	projects, err := circleciClient.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("Unable to get circleCI user projects. CircleCI API error: %s", err.Error())
	}

	return projects, nil
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
				Value: utils.GetFullNameFromOwnerAndRepo(bi.Owner, bi.Repository),
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
