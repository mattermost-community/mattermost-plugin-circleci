package main

import (
	"fmt"

	"github.com/jszwedko/go-circleci"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	accountTrigger  = "account"
	accountHint     = "<" + accountViewTrigger + "|" + accountConnectTrigger + "|" + accountDisconnectTrigger + ">"
	accountHelpText = "Manage the connection to your CircleCI acccount"

	accountViewTrigger  = "view"
	AccountViewHelpText = "Get informations about yourself"

	accountConnectTrigger  = "connect"
	accountConnectHint     = "<API token>"
	accountConnectHelpText = "Connect your Mattermost account to CircleCI"

	accountDisconnectTrigger  = "disconnect"
	accountDisconnectHelpText = "Disconnect your Mattermost account from CircleCI"
)

func (p *Plugin) executeAccount(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case accountViewTrigger:
		return p.executeAccountView(args, circleciToken)

	case accountConnectTrigger:
		return p.executeAccountConnect(args, split[1:])

	case accountDisconnectTrigger:
		return p.executeAccountDisconnect(args)

	case commandHelpTrigger:
		return p.sendHelpResponse(args, accountTrigger)

	default:
		return p.sendIncorrectSubcommandResponse(args, accountTrigger)
	}
}

func (p *Plugin) executeAccountView(args *model.CommandArgs, token string) (*model.CommandResponse, *model.AppError) {
	user, ok := p.getCircleUserInfo(token)
	if !ok {
		p.API.LogInfo("Unable to get CircleCI info", "MM UserID", args.UserId)
		return p.sendEphemeralResponse(args, errorConnectionText), nil
	}

	projects, _ := p.getCircleciUserProjects(token)
	projectsListString := ""
	for _, project := range projects {
		// TODO : add circleCI url
		projectsListString += fmt.Sprintf("- [%s](%s) owned by %s\n", project.Reponame, project.VCSURL, project.Username)
	}

	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				ThumbURL: user.AvatarURL,
				Fallback: "User:" + circleciUserToString(user) + ". Email:" + *user.SelectedEmail,
				Pretext:  "Information for CircleCI user " + circleciUserToString(user),
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Name",
						Value: user.Name,
						Short: true,
					},
					{
						Title: "Email",
						Value: user.SelectedEmail,
						Short: true,
					},
					{
						Title: "Followed projects",
						Value: projectsListString,
						Short: false,
					},
				},
			},
		},
	)

	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeAccountConnect(args *model.CommandArgs, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 1 {
		return p.sendEphemeralResponse(args, "Please tell me your token. If you don't have a CircleCI Personal API Token, you can get one from your [Account Dashboard](https://circleci.com/account/api)"), nil
	}

	if token, exists := p.getTokenKV(args.UserId); exists {
		user, ok := p.getCircleUserInfo(token)
		if !ok {
			return p.sendEphemeralResponse(args, "Internal error when reaching CircleCI"), nil
		}

		return p.sendEphemeralResponse(args, "You are already connected as "+circleciUserToString(user)), nil
	}

	circleciToken := split[0]
	circleciClient := &circleci.Client{
		Token: circleciToken,
	}

	user, err := circleciClient.Me()
	if err != nil {
		p.API.LogError("Error when reaching CircleCI", "CircleCI error:", err)
		return p.sendEphemeralResponse(args, "Can't connect to CircleCI. Please check that your user API token is valid"), nil
	}

	if ok := p.storeTokenKV(args.UserId, circleciToken); !ok {
		return p.sendEphemeralResponse(args, "Internal error when storing your token"), nil
	}

	return p.sendEphemeralResponse(args, "Successfully connected to CircleCI as "+circleciUserToString(user)), nil
}

func (p *Plugin) executeAccountDisconnect(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if ok := p.deleteTokenKV(args.UserId); !ok {
		return p.sendEphemeralResponse(args, errorConnectionText), nil
	}

	return p.sendEphemeralResponse(args, "Your CircleCI account has been successfully disconnected from Mattermost"), nil
}
