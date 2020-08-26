package main

import (
	"github.com/jszwedko/go-circleci"
	"github.com/mattermost/mattermost-server/v5/model"
)

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

func (p *Plugin) getCircleCIUserInfo(circleCIToken string) (*circleci.User, bool) {
	circleciClient := &circleci.Client{
		Token: circleCIToken,
	}

	user, err := circleciClient.Me()
	if err != nil {
		p.API.LogError("Error when reaching CircleCI", "CircleCI error:", err)
		return nil, false
	}

	return user, true
}

func getFormattedNameAndLogin(user *circleci.User) string {
	if *user.Name != "" {
		return *user.Name + " (" + user.Login + ")"
	}
	return user.Login
}
