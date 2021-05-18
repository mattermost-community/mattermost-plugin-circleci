package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// Plugin utils

// sendEphemeralPost send an ephemeral post as a response to a slash command
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

// sendEphemeralResponse is a wrapper of p.sendEphemeralPost() to one-line the return statements in all executeCommand functions
func (p *Plugin) sendEphemeralResponse(args *model.CommandArgs, message string) *model.CommandResponse {
	p.sendEphemeralPost(args, message, nil)
	return &model.CommandResponse{}
}

// createPost creates a post or logs an error
func (p *Plugin) createPost(r *model.Post) *model.Post {
	post, appErr := p.API.CreatePost(r)
	if appErr != nil {
		p.API.LogError("Error when creating post", "appError", appErr)
		return nil
	}

	return post
}

// getWebhookURL generates the webhook URL from the SiteURL and the webhook secret
func (p *Plugin) getWebhookURL() string {
	siteURL := *p.API.GetConfig().ServiceSettings.SiteURL
	siteURL = strings.TrimRight(siteURL, "/")
	webhookSecret := p.getConfiguration().WebhooksSecret
	return fmt.Sprintf("%s/plugins/%s%s/%s", siteURL, manifest.Id, routeWebhooks, webhookSecret)
}

// getUsername returns the username of the given user, or "unknown user" if not found
func (p *Plugin) getUsername(userID string) string {
	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		p.API.LogError("Unable to get user informations", "appError", appErr, "userID", userID)
		return "unknown user"
	}

	return user.Username
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
