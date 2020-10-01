package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jszwedko/go-circleci"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"
)

const (
	routeAutocompleteFollowedProjects = "/autocomplete/followedProjects"
	routeWebhooksPrefix               = "/hooks"
)

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	// Check token
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.Store.GetTokenForUser(userID)
	if !exists {
		http.NotFound(w, r)
	}

	routeWebhooks := strings.Join([]string{
		routeWebhooksPrefix,
		p.getConfiguration().WebhooksSecret,
	},
		"/",
	)

	p.API.LogDebug("Receveid CircleCI http request", "URL", r.URL.Path, "route for CircleCI Webhooks", routeWebhooks)

	// Call the handler
	switch r.URL.Path {
	case routeAutocompleteFollowedProjects:
		httpAutocompleteFollowedProject(p, w, r, circleciToken)
		return

	case routeWebhooks:
		httpHandleWebhook(p, w, r)
		return

	default:
		http.NotFound(w, r)
	}
}

func httpAutocompleteFollowedProject(p *Plugin, w http.ResponseWriter, r *http.Request, circleciToken string) {
	if r.Method != http.MethodGet {
		p.respondAndLogErr(w, http.StatusMethodNotAllowed, errors.New("method"+r.Method+"is not allowed, must be GET"))
		return
	}

	circleciClient := &circleci.Client{Token: circleciToken}
	projects, err := circleciClient.ListProjects()
	if err != nil {
		p.respondAndLogErr(w, http.StatusInternalServerError, err)
		return
	}

	out := []model.AutocompleteListItem{
		{
			HelpText: "Manually type the project's VCS repository name",
			Item:     "[repository]",
		},
	}
	if len(projects) == 0 {
		p.respondJSON(w, out)
		return
	}

	for _, project := range projects {
		out = append(out, model.AutocompleteListItem{
			HelpText: project.VCSURL,
			Item:     project.Reponame,
		})
	}
	p.respondJSON(w, out)
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

func httpHandleWebhook(p *Plugin, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p.respondAndLogErr(w, http.StatusMethodNotAllowed, errors.New("method"+r.Method+"is not allowed, must be POST"))
		return
	}

	buildInfos := new(v1.BuildInfos)
	if err := json.NewDecoder(r.Body).Decode(&buildInfos); err != nil {
		p.API.LogError("Unable to decode JSON for received webkook.", "Error", err.Error())
		return
	}

	allSubs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return
	}

	channelsToPost := allSubs.GetSubscribedChannelsForRepository(buildInfos.Owner, buildInfos.Repository)
	if channelsToPost == nil {
		p.API.LogWarn("Received webhooks without any subscriptions", "webhook", buildInfos)
	}

	postWithoutChannel := &model.Post{
		UserId: p.botUserID,
	}
	postWithoutChannel.AddProp("attachments", buildInfos.ToPostAttachments(buildFailedIconURL, buildGreenIconURL))

	for _, channel := range channelsToPost {
		post := postWithoutChannel.Clone()
		post.ChannelId = channel

		_, appErr := p.API.CreatePost(post)
		if appErr != nil {
			p.API.LogError("Failed to create Post", "appError", appErr)
		}
	}
}
