package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

type BuildInfos struct {
	Owner          string `json:"Owner"`
	Repository     string `json:"Repository"`
	CircleBuildNum int    `json:"CircleBuildNum"`
	Failed         bool   `json:"Failed"`
	Message        string `json:"Message"`
}

// Convert the build info into a post attachment
func (bi *BuildInfos) toPostAttachments() []*model.SlackAttachment {
	// TODO add link to build
	attachment := &model.SlackAttachment{
		Fields: []*model.SlackAttachmentField{
			{
				Title: "Repo",
				Short: true,
				Value: getFullNameFromOwnerAndRepo(bi.Owner, bi.Repository),
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

	buildInfos := new(BuildInfos)
	if err := json.NewDecoder(r.Body).Decode(&buildInfos); err != nil {
		p.API.LogError("Unable to decode JSON for received webkook.", "Error", err.Error())
		return
	}

	allSubs, err := p.getSubscriptionsKV()
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
	postWithoutChannel.AddProp("attachments", buildInfos.toPostAttachments())

	for _, channel := range channelsToPost {
		post := postWithoutChannel.Clone()
		post.ChannelId = channel

		_, appErr := p.API.CreatePost(post)
		if appErr != nil {
			p.API.LogError("Failed to create Post", "appError", appErr)
		}
	}
}
