package circle

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

// TODO : rename with a more meaningful name
type BuildInfos struct {
	Owner                  string `json:"Owner"`
	Repository             string `json:"Repository"`
	RepositoryURL          string `json:"RepositoryURL"`
	CircleBuildNum         int    `json:"CircleBuildNum"`
	CircleBuildURL         string `json:"CircleBuildURL"`
	Username               string `json:"Username"`
	AssociatedPullRequests string `json:"AssociatedPullRequests"`
	IsFailed               bool   `json:"IsFailed"`
	IsWaitingApproval      bool   `json:"IsWaitingApproval"`
	Message                string `json:"Message"`
}

// Convert the build info into a post attachment
func (bi *BuildInfos) ToPost(buildFailedIconURL, buildGreenIconURL string) *model.Post {
	attachment := &model.SlackAttachment{
		TitleLink: bi.CircleBuildURL,
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
			{
				Title: "Build informations",
				Short: false,
				Value: fmt.Sprintf(
					"- Build triggered by: %s\n- Associated PRs: %s",
					bi.Username,
					bi.AssociatedPullRequests,
				),
			},
		},
	}

	switch {
	case bi.IsFailed:
		attachment.ThumbURL = buildFailedIconURL
		attachment.Title = "CircleCI Job failed"
		attachment.Color = "#FF1919" // red

	case bi.IsWaitingApproval:
		attachment.Title = "CircleCI Job waiting approval"
		attachment.Color = "#DBAB09" // yellow

		// TODO : add button to approve / refuse the job
		// attachment.Actions = []*model.PostAction{
		// 	{
		// 		Id:   "approve-circleci-job",
		// 		Name: "Approve Job",
		// 		Integration: &model.PostActionIntegration{
		// 			URL: "",
		// 			Context: map[string]interface{}{
		// 				"a": "b",
		// 			},
		// 		},
		// 	},
		// }

	default:
		// Not failed and not waiting approval = passed
		attachment.ThumbURL = buildGreenIconURL
		attachment.Title = "CircleCI Job passed"
		attachment.Color = "#50F100" // green
	}

	attachment.Fallback = attachment.Title

	post := model.Post{
		Message: bi.Message,
	}

	post.AddProp("attachments", []*model.SlackAttachment{
		attachment,
	})

	return &post
}
