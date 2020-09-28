package commands

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	accountTrigger  = "account"
	accountHint     = "<" + accountViewTrigger + "|" + accountConnectTrigger + "|" + accountDisconnectTrigger + ">"
	accountHelpText = "Manage the connection to your CircleCI acccount"

	accountConnectTrigger  = "connect"
	accountConnectHint     = "<API token>"
	accountConnectHelpText = "Connect your Mattermost account to CircleCI"

	accountDisconnectTrigger  = "disconnect"
	accountDisconnectHelpText = "Disconnect your Mattermost account from CircleCI"

	accountViewTrigger  = "view"
	accountViewHelpText = "Get informations about yourself"

	// Help message for the command
	accountHelpMessage = "#### Connect to your CircleCI account\n" +
		"* `/" + mainTrigger + " " + accountTrigger + " " + accountViewTrigger + "` — " + accountViewHelpText + "\n" +
		"* `/" + mainTrigger + " " + accountTrigger + " " + accountConnectTrigger + " " + accountConnectHint + "` — " + accountConnectHelpText + "\n" +
		"* `/" + mainTrigger + " " + accountTrigger + " " + accountDisconnectTrigger + "` — " + accountDisconnectHelpText + "\n"
)

func getAccountAutocompleteData() *model.AutocompleteData {
	account := model.NewAutocompleteData(accountTrigger, accountHint, accountHelpText)
	view := model.NewAutocompleteData(accountViewTrigger, "", accountViewHelpText)
	connect := model.NewAutocompleteData(accountConnectTrigger, accountConnectHint, accountConnectHelpText)
	connect.AddTextArgument("Generate a Personal API Token from your CircleCI user settings", accountConnectHint, "")
	disconnect := model.NewAutocompleteData(accountDisconnectTrigger, "", accountDisconnectHelpText)

	account.AddCommand(view)
	account.AddCommand(connect)
	account.AddCommand(disconnect)

	return account
}

// executeAccount the config command
// Return a string which is meant to be shown to the user as an ephemeral post
func executeAccount(args *model.CommandArgs, db store.Store) string {
	split := strings.Split(args.Command, " ")

	subcommand := "help"
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case accountViewTrigger:
		return executeAccountView(args, db)

	case accountConnectTrigger:
		return executeAccountConnect(args, db)

	case accountDisconnectTrigger:
		return executeAccountDisconnect(args, db)

	case helpTrigger:
		return formatHelpMessage(args, accountTrigger)

	default:
		return formatIncorrectSubcommand(args, accountTrigger)
	}
}

func executeAccountConnect(args *model.CommandArgs, db store.Store) string {
	split := strings.Split(args.Command, " ")
	if len(split) < 2 || split[0] != mainTrigger || split[1] != accountTrigger || split[2] != accountConnectTrigger {
		panic("This function should not been called")
	}

	if len(split) < 3 {
		return "Please tell me your token. If you don't have a CircleCI Personal API Token, you can get one from your [Account Dashboard](https://circleci.com/account/api)"
	}

	if token, exists := db.GetTokenForUser(args.UserId); exists {
		user, err := circle.GetCurrentUser(token)
		if err != nil {
			// TODO Should log error
			return "Internal error when reaching CircleCI"
		}

		return "You are already connected as " + circle.CircleciUserToString(user)
	}

	circleciToken := split[3]

	user, err := circle.GetCurrentUser(circleciToken)
	if err != nil {
		// TODO : should log error p.API.LogError("Error when reaching CircleCI", "CircleCI error:", err)
		return "Can't connect to CircleCI. Please check that your user API token is valid"
	}

	if ok := db.StoreTokenForUser(args.UserId, circleciToken); !ok {
		return "Internal error when storing your token"
	}

	return fmt.Sprintf("Successfully connected to CircleCI as %s (%s)", user.Name, user.Login)
}

func executeAccountDisconnect(args *model.CommandArgs, db store.Store) string {
	if ok := db.DeleteTokenForUser(args.UserId); !ok {
		return errorConnectionText
	}

	return "Your CircleCI account has been successfully disconnected from Mattermost"
}

func executeAccountView(args *model.CommandArgs, db store.Store) string {
	return "This command is not supported yet, sorry"
	// TODO this needs to return a SlackAttachment
	// token, ok := db.GetTokenForUser(args.UserId)
	// if !ok {
	// 	return ErrorConnectionText
	// }

	// user, err := circle.GetCurrentUser(token)
	// if err != nil {
	// 	// TODO Should log error p.API.LogInfo("Unable to get CircleCI info", "MM UserID", args.UserId)
	// 	return ErrorConnectionText
	// }

	// projects, _ := p.getCircleciUserProjects(token)
	// projectsListString := ""
	// for _, project := range projects {
	// 	// TODO : add circleCI url
	// 	projectsListString += fmt.Sprintf("- [%s](%s) owned by %s\n", project.Reponame, project.VCSURL, project.Username)
	// }

	// _ = p.sendEphemeralPost(
	// 	args,
	// 	"",
	// 	[]*model.SlackAttachment{
	// 		{
	// 			ThumbURL: user.AvatarURL,
	// 			Fallback: "User:" + circleciUserToString(user) + ". Email:" + *user.SelectedEmail,
	// 			Pretext:  "Information for CircleCI user " + circleciUserToString(user),
	// 			Fields: []*model.SlackAttachmentField{
	// 				{
	// 					Title: "Name",
	// 					Value: user.Name,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Email",
	// 					Value: user.SelectedEmail,
	// 					Short: true,
	// 				},
	// 				{
	// 					Title: "Followed projects",
	// 					Value: projectsListString,
	// 					Short: false,
	// 				},
	// 			},
	// 		},
	// 	},
	// )

	// return &model.CommandResponse{}, nil
}
