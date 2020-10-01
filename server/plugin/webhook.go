package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"
)

// TODO Add information to the notification: which branCIRCLE_BRANCHche is concerned / which commit? (may need modification to the orb)

// WebhookInfo from the webhookCIRCLE_BRANCH
type WebhookInfo struct {
	Owner                  string `json:"Owner"`
	Repository             string `json:"Repository"`
	RepositoryURL          string `json:"RepositoryURL"`
	Branch                 string `json:"Branch"`
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
					v1.GetFullNameFromOwnerAndRepo(wi.Owner, wi.Repository),
					wi.RepositoryURL,
				),
			},
			{
				Title: "Branch",
				Short: true,
				Value: fmt.Sprintf("`%s`", wi.Branch),
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

func httpHandleWebhook(p *Plugin, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p.respondAndLogErr(w, http.StatusMethodNotAllowed, errors.New("method "+r.Method+" is not allowed, must be POST"))
		return
	}

	wi := new(WebhookInfo)
	if err := json.NewDecoder(r.Body).Decode(&wi); err != nil {
		p.API.LogError("Unable to decode JSON for received webhook.", "Error", err.Error())
		return
	}

	allSubs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return
	}

	channelsToPost := allSubs.GetFilteredChannelsForBuild(wi.Owner, wi.Repository, wi.IsFailed)
	if channelsToPost == nil {
		p.API.LogWarn("The received webhook doesn't match any subscriptions (or flags)", "webhook", wi)
	}

	postWithoutChannel := wi.ToPost(buildFailedIconURL, buildGreenIconURL)
	postWithoutChannel.UserId = p.botUserID

	for _, channel := range channelsToPost {
		post := postWithoutChannel.Clone()
		post.ChannelId = channel

		_, appErr := p.API.CreatePost(post)
		if appErr != nil {
			p.API.LogError("Failed to create Post", "appError", appErr)
		}
	}
}
