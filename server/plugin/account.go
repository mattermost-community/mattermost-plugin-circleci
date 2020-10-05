package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
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

func getAccountAutoCompleteData() *model.AutocompleteData {
	account := model.NewAutocompleteData(accountTrigger, accountHint, accountHelpText)

	view := model.NewAutocompleteData(accountViewTrigger, "", AccountViewHelpText)
	connect := model.NewAutocompleteData(accountConnectTrigger, accountConnectHint, accountConnectHelpText)
	connect.AddTextArgument("Generate a Personal API Token from your CircleCI user settings", accountConnectHint, "")
	disconnect := model.NewAutocompleteData(accountDisconnectTrigger, "", accountDisconnectHelpText)

	account.AddCommand(view)
	account.AddCommand(connect)
	account.AddCommand(disconnect)

	return account
}

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
	user, err := v1.GetCircleUserInfo(token)
	if err != nil {
		p.API.LogInfo("Unable to get CircleCI info", "MM UserID", args.UserId)
		return p.sendEphemeralResponse(args, errorConnectionText), nil
	}

	projects, _ := v1.GetCircleciUserProjects(token)
	projectsListString := ""
	for _, project := range projects {
		projectIdentifier := &store.ProjectIdentifier{
			VCSType: "gh",
			Org:     project.Username,
			Project: project.Reponame,
		}
		if strings.Contains(project.VCSURL, "https://bitbucket.org") {
			projectIdentifier.VCSType = "bb"
		}

		projectsListString += fmt.Sprintf("%s - [CircleCI page](%s)\n", projectIdentifier.ToMarkdown(), projectIdentifier.ToCircleURL())
	}

	_ = p.sendEphemeralPost(
		args,
		"",
		[]*model.SlackAttachment{
			{
				ThumbURL: user.AvatarURL,
				Fallback: "User:" + v1.CircleciUserToString(user) + ". Email:" + *user.SelectedEmail,
				Pretext:  "Information for CircleCI user " + v1.CircleciUserToString(user),
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

	if token, exists := p.Store.GetTokenForUser(args.UserId, p.getConfiguration().EncryptionKey); exists {
		user, err := v1.GetCircleUserInfo(token)
		if err != nil {
			p.API.LogWarn("Internal error when reaching CircleCI", "error", err)
			return p.sendEphemeralResponse(args, ":red_circle: Internal error when reaching CircleCI"), nil
		}

		return p.sendEphemeralResponse(args,
			fmt.Sprintf("You are already connected as %s. Please disconnect first.", v1.CircleciUserToString(user)),
		), nil
	}

	circleciToken := split[0]

	user, err := circle.GetCurrentUser(circleciToken)
	if err != nil {
		p.API.LogError("Error when reaching CircleCI", "CircleCI error:", err)
		return p.sendEphemeralResponse(args, ":red_circle: Can't connect to CircleCI. Please check that your user API token is valid"), nil
	}

	if ok := p.Store.StoreTokenForUser(args.UserId, circleciToken, p.getConfiguration().EncryptionKey); !ok {
		return p.sendEphemeralResponse(args, ":red_circle: Internal error when storing your token"), nil
	}

	return p.sendEphemeralResponse(args, fmt.Sprintf(":white_check_mark: Successfully connected to CircleCI as %s (%s)", user.Name, user.Login)), nil
}

func (p *Plugin) executeAccountDisconnect(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if ok := p.Store.DeleteTokenForUser(args.UserId); !ok {
		return p.sendEphemeralResponse(args, errorConnectionText), nil
	}

	return p.sendEphemeralResponse(args, ":white_check_mark: Your CircleCI account has been successfully disconnected from Mattermost"), nil
}
