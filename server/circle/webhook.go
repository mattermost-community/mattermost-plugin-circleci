package circle

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

type WebhookInfo struct {
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
func (wi *WebhookInfo) ToPost(buildFailedIconURL, buildGreenIconURL string) *model.Post {
	attachment := &model.SlackAttachment{
		TitleLink: wi.CircleBuildURL,
		Fields: []*model.SlackAttachmentField{
			{
				Title: "Repo",
				Short: true,
				Value: fmt.Sprintf(
					"[%s](%s)",
					GetFullNameFromOwnerAndRepo(wi.Owner, wi.Repository),
					wi.RepositoryURL,
				),
			},
			{
				Title: "Job number",
				Short: true,
				Value: fmt.Sprintf("%d", wi.CircleBuildNum),
			},
			{
				Title: "Build informations",
				Short: false,
				Value: fmt.Sprintf(
					"- Build triggered by: %s\n- Associated PRs: %s",
					wi.Username,
					wi.AssociatedPullRequests,
				),
			},
		},
	}

	switch {
	case wi.IsFailed:
		attachment.ThumbURL = buildFailedIconURL
		attachment.Title = "CircleCI Job failed"
		attachment.Color = "#FF1919" // red

	case wi.IsWaitingApproval:
		attachment.Title = "CircleCI Job waiting approval"
		attachment.Color = "#8267E4" // purple

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
		Message: wi.Message,
	}

	post.AddProp("attachments", []*model.SlackAttachment{
		attachment,
	})

	return &post
}
