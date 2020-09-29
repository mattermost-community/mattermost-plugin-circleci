package main

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"

	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"

	"github.com/pkg/errors"
)

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
