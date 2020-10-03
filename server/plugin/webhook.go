package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

// WebhookInfo from the webhook
type WebhookInfo struct {
	Organization           string `json:"Organization"`
	Repository             string `json:"Repository"`
	RepositoryURL          string `json:"RepositoryURL"`
	Username               string `json:"Username"`
	WorkflowID             string `json:"WorkflowID"`
	JobName                string `json:"JobName"`
	CircleBuildNumber      int    `json:"CircleBuildNumber"`
	CircleBuildURL         string `json:"CircleBuildURL"`
	Branch                 string `json:"Branch"`
	Tag                    string `json:"Tag"`
	Commit                 string `json:"Commit"`
	AssociatedPullRequests string `json:"AssociatedPullRequests"`
	IsFailed               bool   `json:"IsFailed"`
	IsWaitingApproval      bool   `json:"IsWaitingApproval"`
	Message                string `json:"Message"`
}

func (wi *WebhookInfo) toProjectIdentifier() *store.ProjectIdentifier {
	repoType := "gh"
	if strings.Contains(wi.RepositoryURL, "git@bitbucket.org") {
		repoType = "bb"
	}
	repoConf, _ := store.CreateProjectIdentifierFromSlug(fmt.Sprintf("%s/%s/%s", repoType, wi.Organization, wi.Repository))
	return repoConf
}

// Convert the build info into a Post
func (wi *WebhookInfo) ToPost(buildFailedIconURL, buildGreenIconURL string) *model.Post {
	if wi.AssociatedPullRequests == "" {
		wi.AssociatedPullRequests = ":grey_question: No PR"
	}

	repo := wi.toProjectIdentifier()

	attachment := &model.SlackAttachment{
		TitleLink: wi.CircleBuildURL,
		Fields: []*model.SlackAttachmentField{
			{
				Title: "Project",
				Short: true,
				Value: repo.ToMarkdown(),
			},
			{
				Title: "Branch",
				Short: true,
				Value: fmt.Sprintf("`%s`", wi.Branch),
			},
			{
				Title: "Commit",
				Short: true,
				Value: fmt.Sprintf("`%s`", wi.Commit),
			},
			{
				Title: "Job number",
				Short: true,
				Value: fmt.Sprintf("%d", wi.CircleBuildNumber),
			},
			{
				Title: "Job informations",
				Short: false,
				Value: fmt.Sprintf(
					"- Build triggered by: %s\n- Associated PRs: %s\n",
					wi.Username,
					wi.AssociatedPullRequests,
				),
			},
		},
	}

	switch {
	case wi.IsWaitingApproval:
		attachment.Title = "You have a CircleCI Workflow waiting for approval"
		attachment.Color = "#8267E4" // purple
		attachment.Actions = []*model.PostAction{
			{
				Id:   "approve-circleci-job",
				Name: "Approve Job",
				Type: model.POST_ACTION_TYPE_BUTTON,
				Integration: &model.PostActionIntegration{
				URL: "/plugins/job/approve/",
					Context: map[string]interface{}{
						"WorkflowID": wi.WorkflowID,
					},
				},
			},
		}

	case wi.IsFailed:
		attachment.ThumbURL = buildFailedIconURL
		attachment.Title = fmt.Sprintf("Your CircleCI Job has failed: %s", wi.JobName)
		attachment.Color = "#FF1919" // red

	default:
		// Not failed and not waiting approval = passed
		attachment.ThumbURL = buildGreenIconURL
		attachment.Title = fmt.Sprintf("Your CircleCI Job has passed: %s", wi.JobName)
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

func (p *Plugin) httpHandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Checking secret
	vars := mux.Vars(r)
	if vars["secret"] != p.getConfiguration().WebhooksSecret {
		http.NotFound(w, r)
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

	channelsToPost := allSubs.GetFilteredChannelsForJob(wi.toProjectIdentifier(), wi.IsFailed)
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
