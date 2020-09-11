package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	commandTrigger          = "circleci"
	commandAutocompleteHint = "<" + accountTrigger + "|" + projectTrigger + "|" + subscribeTrigger + ">"
	commandAutocompleteDesc = "Interact with CircleCI jobs and builds"

	notConnectedText    = "You are not connected to CircleCI. Please try `/" + commandTrigger + " " + accountTrigger + " " + accountConnectTrigger + "`"
	errorConnectionText = "Error when reaching to CircleCI. Please check that your token is still valid"

	// All the Triggers and HelpTexts for the subcommands are defined in the corresponding commands_*.go file
	commandHelpTrigger = "help"

	accountHelp = "#### Connect to your CircleCI account\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountViewTrigger + "` — " + AccountViewHelpText + "\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountConnectTrigger + " " + accountConnectHint + "` — " + accountConnectHelpText + "\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountDisconnectTrigger + "` — " + accountDisconnectHelpText + "\n"

	projectHelp = "#### Manage CircleCI projects\n" +
		"* `/" + commandTrigger + " " + projectTrigger + " " + projectListTrigger + "` — " + projectListHelpText + "\n" +
		"* `/" + commandTrigger + " " + projectTrigger + " " + projectRecentBuildsTrigger + " " + projectRecentBuildsHint + "` — " + projectRecentBuildsHelpText + "\n"

	help = "## CircleCI plugin Help\n" + accountHelp + projectHelp
)

func (p *Plugin) getCommand() *model.Command {
	return &model.Command{
		Trigger:              commandTrigger,
		AutoComplete:         true,
		AutoCompleteDesc:     commandAutocompleteDesc,
		AutoCompleteHint:     commandAutocompleteHint,
		AutocompleteData:     getAutocompleteData(),
		AutocompleteIconData: getAutocompleteIconData(p),
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
	projectRecentBuild.AddDynamicListArgument("", routeAutocompleteFollowedProjects, true)
	projectRecentBuild.AddTextArgument("Branch name", "[branch]", "")

	project.AddCommand(projectRecentBuild)
	project.AddCommand(projectList)

	// Subscriptions subcommands
	subscribe := model.NewAutocompleteData(subscribeTrigger, subscribeHint, subscribeHelpText)

	subscribeList := model.NewAutocompleteData(subscribeListTrigger, subscribeListHint, subscribeListHelpText)
	subscribeChannel := model.NewAutocompleteData(subscribeChannelTrigger, subscribeChannelHint, subscribeChannelHelpText)
	subscribeChannel.AddTextArgument("Owner of the project's repository", "[owner]", "")
	subscribeChannel.AddDynamicListArgument("", routeAutocompleteFollowedProjects, true)
	subscribeChannel.AddNamedTextArgument(flagOnlyFailedBuilds, "Only receive notifications for failed builds", "[write anything here]", "", false)
	unsubscribeChannel := model.NewAutocompleteData(subscribeUnsubscribeChannelTrigger, subscribeUnsubscribeChannelHint, subscribeUnsubscribeChannelHelpText)
	unsubscribeChannel.AddTextArgument("Owner of the project's repository", "[owner]", "") // TODO make dynamic autocomplete list
	unsubscribeChannel.AddTextArgument("Repository name", "[repository]", "")              // TODO make dynamic autocomplete list
	listAllSubscribedChannels := model.NewAutocompleteData(subscribeListAllChannelsTrigger, subscribeListAllChannelsHint, subscribeListAllChannelsHelpText)
	listAllSubscribedChannels.AddTextArgument("Owner of the project's repository", "[owner]", "") // TODO make dynamic autocomplete list
	listAllSubscribedChannels.AddTextArgument("Repository name", "[repository]", "")              // TODO make dynamic autocomplete list

	subscribe.AddCommand(subscribeList)
	subscribe.AddCommand(subscribeChannel)
	subscribe.AddCommand(unsubscribeChannel)
	subscribe.AddCommand(listAllSubscribedChannels)

	// Add all subcommands
	mainCommand.AddCommand(account)
	mainCommand.AddCommand(project)
	mainCommand.AddCommand(subscribe)
	return mainCommand
}

func getAutocompleteIconData(p *Plugin) string {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogError("Couldn't get bundle path", "error", err)
		return ""
	}

	icon, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "circleci.svg"))
	if err != nil {
		p.API.LogError("Failed to open icon", "error", err)
		return ""
	}

	return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(icon))
}

func (p *Plugin) sendHelpResponse(args *model.CommandArgs, currentCommand string) (*model.CommandResponse, *model.AppError) {
	message := ""

	switch currentCommand {
	case accountTrigger:
		message += accountHelp

	case projectTrigger:
		message += projectHelp

	default:
		message += help
	}

	return p.sendEphemeralResponse(args, message), nil
}

func (p *Plugin) sendIncorrectSubcommandResponse(args *model.CommandArgs, currentCommand string) (*model.CommandResponse, *model.AppError) {
	message := "Invalid subcommand given. Type `/" + commandTrigger

	if currentCommand != "" {
		message += " " + currentCommand
	}

	message += " help` to get a hint"

	return p.sendEphemeralResponse(args, message), nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)

	command := ""
	if 1 < len(split) {
		command = split[1]
	}

	token, shouldBeConnected := getTokenIfConnected(p, split, args.UserId)
	if shouldBeConnected {
		return p.sendEphemeralResponse(args, notConnectedText), nil
	}

	switch command {
	case accountTrigger:
		return p.executeAccount(args, token, split[2:])

	case projectTrigger:
		return p.executeProject(args, token, split[2:])

	case commandHelpTrigger:
		return p.sendHelpResponse(args, "")

	case subscribeTrigger:
		return p.executeSubscribe(args, token, split[2:])

	default:
		return p.sendIncorrectSubcommandResponse(args, "")
	}
}

// Return the token if it exists, or "", and true if the user should be connected to use this command
func getTokenIfConnected(p *Plugin, split []string, userID string) (string, bool) {
	if len(split) <= 1 {
		return "", false
	}

	// If command is /account connect, no need to be connected
	if split[1] == accountTrigger && 2 < len(split) && split[2] == accountConnectTrigger {
		return "", false
	}

	// If command is help, same thing
	if split[1] == commandHelpTrigger {
		return "", false
	}

	circleToken, exists := p.getTokenKV(userID)
	if !exists {
		return "", true
	}

	return circleToken, false
}
