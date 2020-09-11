package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	commandTrigger          = "circleci"
	commandAutocompleteHint = "[command]"
	commandAutocompleteDesc = "Available commands: " + accountTrigger + ", " + projectTrigger

	notConnectedText    = "You are not connected to CircleCI. Please try `/" + commandTrigger + " " + accountTrigger + " " + accountConnectTrigger + "`"
	errorConnectionText = "Error when reaching to CircleCI. Please check that your token is still valid"

	// All the Triggers and HelpTexts for the subcommands are defined in the corresponding commands_*.go file
	commandHelpTrigger = "help"

	accountHelp = "#### Connect to your CircleCI account\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountViewTrigger + "` - " + AccountViewHelpText + "\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountConnectTrigger + " " + accountConnectHint + "` - " + accountConnectHelpText + "\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountDisconnectTrigger + "` - " + accountDisconnectHelpText + "\n"

	projectHelp = "#### Manage CircleCI projects\n" +
		"* `/" + commandTrigger + " " + projectTrigger + " " + projectListTrigger + "` - " + projectListHelpText + "\n" +
		"* `/" + commandTrigger + " " + projectTrigger + " " + projectRecentBuildsTrigger + " " + projectRecentBuildsHint + "` - " + projectRecentBuildsHelpText + "\n"

	help = "## CircleCI plugin\n" + accountHelp + projectHelp
)

func (p *Plugin) getCommand() *model.Command {
	return &model.Command{
		Trigger:          commandTrigger,
		AutoComplete:     true,
		AutoCompleteDesc: commandAutocompleteDesc,
		AutoCompleteHint: commandAutocompleteHint,
		AutocompleteData: getAutocompleteData(),
	}
}

func getAutocompleteData() *model.AutocompleteData {
	mainCommand := model.NewAutocompleteData(commandTrigger, commandAutocompleteHint, commandAutocompleteDesc)

	// Account subcommands
	account := model.NewAutocompleteData(accountTrigger, accountHint, accountHelpText)
	view := model.NewAutocompleteData(accountViewTrigger, "", AccountViewHelpText)
	connect := model.NewAutocompleteData(accountConnectTrigger, accountConnectHint, accountConnectHelpText)
	connect.AddTextArgument("Generate a Personal API Token from your CircleCI user settings", accountConnectHint, "")
	disconnect := model.NewAutocompleteData(accountDisconnectTrigger, "", accountDisconnectHelpText)

	account.AddCommand(view)
	account.AddCommand(connect)
	account.AddCommand(disconnect)

	// Project management subcommands
	project := model.NewAutocompleteData(projectTrigger, projectHint, projectHelpText)
	projectList := model.NewAutocompleteData(projectListTrigger, projectListHint, projectListHelpText)
	projectRecentBuild := model.NewAutocompleteData(projectRecentBuildsTrigger, projectRecentBuildsHint, projectRecentBuildsHelpText)
	projectRecentBuild.AddTextArgument("Owner of the project's repository", "[username]", "")
	projectRecentBuild.AddDynamicListArgument("dzeiufuazifhzauefio", routeAutocompleteFollowedProjects, true)
	projectRecentBuild.AddTextArgument("Branch name", "[branch]", "")

	project.AddCommand(projectList)
	project.AddCommand(projectRecentBuild)

	// Add all subcommands
	mainCommand.AddCommand(account)
	mainCommand.AddCommand(project)
	return mainCommand
}

// Send an ephemeralCommandResponse (instead of a post from the bot)
// to be able to override user image to a green one
func (p *Plugin) sendHelpResponse(currentCommand string) (*model.CommandResponse, *model.AppError) {
	message := "CircleCI Plugin help\n"

	switch currentCommand {
	case accountTrigger:
		message += accountHelp

	case projectTrigger:
		message += projectHelp

	default:
		message += help
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Username:     botUserName,
		IconURL:      p.iconBuildGreenURL,
		Text:         message,
	}, nil
}

// Send an ephemeralCommandResponse (instead of a post from the bot)
// to be able to override user image to a red one
func (p *Plugin) sendIncorrectSubcommandResponse(currentCommand string) (*model.CommandResponse, *model.AppError) {
	message := "Invalid subcommand given. Type `/" + commandTrigger

	if currentCommand != "" {
		message += " " + currentCommand
	}

	message += " help` to get a hint"

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Username:     botUserName,
		IconURL:      p.iconBuildFailURL,
		Text:         message,
	}, nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)

	command := ""
	if 1 < len(split) {
		command = split[1]
	}

	// if command is not "connect" or "help", check that the user is connected
	token := ""
	if command != accountTrigger && command != "help" {
		circlecitoken, exists := p.getTokenFromKVStore(args.UserId)
		if !exists {
			return p.sendEphemeralResponse(args, notConnectedText), nil
		}
		token = circlecitoken
	}

	switch command {
	case accountTrigger:
		return p.executeAccount(args, token, split[2:])

	case projectTrigger:
		return p.executeProject(args, token, split[2:])

	case commandHelpTrigger:
		return p.sendHelpResponse("")

	default:
		return p.sendIncorrectSubcommandResponse("")
	}
}
