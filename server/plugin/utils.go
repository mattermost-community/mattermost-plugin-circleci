package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
	"github.com/pkg/errors"
)

// Plugin utils

func (p *Plugin) sendEphemeralPost(args *model.CommandArgs, message string, attachments []*model.SlackAttachment) *model.Post {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   message,
	}

	if attachments != nil {
		post.AddProp("attachments", attachments)
	}

	return p.API.SendEphemeralPost(
		args.UserId,
		post,
	)
}

// Wrapper of p.sendEphemeralPost() to one-line the return statements in all executeCommand functions
func (p *Plugin) sendEphemeralResponse(args *model.CommandArgs, message string) *model.CommandResponse {
	p.sendEphemeralPost(args, message, nil)
	return &model.CommandResponse{}
}

func (p *Plugin) getWebhookURL() string {
	siteURL := *p.API.GetConfig().ServiceSettings.SiteURL
	siteURL = strings.TrimRight(siteURL, "/")
	webhookSecret := p.getConfiguration().WebhooksSecret
	return fmt.Sprintf("%s/plugins/%s%s/%s", siteURL, manifest.Id, routeWebhooks, webhookSecret)
}

// HTTP Utils below

func (p *Plugin) respondAndLogErr(w http.ResponseWriter, code int, err error) {
	http.Error(w, err.Error(), code)
	p.API.LogError(err.Error())
}

func (p *Plugin) respondJSON(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		p.respondAndLogErr(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to marshal response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		p.respondAndLogErr(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to write response"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// overwrite given env var after confirmation if already exist
func (p *Plugin) httpHandleEnvOverwrite(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.Store.GetTokenForUser(userID, p.getConfiguration().EncryptionKey)

	if !exists {
		http.NotFound(w, r)
	}
	requestData := model.PostActionIntegrationRequestFromJson(r.Body)
	if requestData == nil {
		p.API.LogError("Empty request data")
		return
	}
	originalPost, appErr := p.API.GetPost(requestData.PostId)
	if appErr != nil {
		p.API.LogError("Unable to get post", "postID", requestData.PostId)
	} else if _, appErr := p.API.UpdatePost(originalPost); appErr != nil {
		// TODO : remove the button
		p.API.LogError("Unable to update post", "postID", originalPost.Id)
	}

	responsePost := &model.Post{
		ChannelId: requestData.ChannelId,
		RootId:    requestData.PostId,
		UserId:    p.botUserID,
	}

	action := fmt.Sprintf("%s", requestData.Context["Action"])
	projectSlug := fmt.Sprintf("%v", requestData.Context["ProjectSlug"])
	name := fmt.Sprintf("%v", requestData.Context["EnvName"])
	val := fmt.Sprintf("%v", requestData.Context["EnvVal"])

	if action == "deny" {
		responsePost.Message = fmt.Sprintf("Did not overwrite env variable %s for project %s", name, projectSlug)
		p.API.SendEphemeralPost(userID, responsePost)
	} else if action == "approve" {
		err := circle.AddEnvVar(circleciToken, projectSlug, name, val)
		if err != nil {
			p.API.LogError("Error occurred while adding environment variable", err)
			responsePost.Message = fmt.Sprintf(":red_circle: Cannot overwrite env var %s:%s from mattermost.", name, val)
		} else {
			responsePost.Message = fmt.Sprintf(":white_check_mark: Successfully added environment variable %s:%s for project %s", name, val, projectSlug)
		}
		p.API.SendEphemeralPost(userID, responsePost)
	} else {
		p.API.LogError("action %s is not valid", action)
	}
}
