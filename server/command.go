package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	commandTrigger          = "circleci"
	commandAutocompleteHint = "[command]"
	commandAutocompleteDesc = "Available commands: `connect`, `disconnect`, `me`"

	meTrigger          = "me"
	meHelpText         = "Get informations about yourself"
	connectTrigger     = "connect"
	connectHint        = "[API token]"
	connectHelpText    = "Connect your Mattermost account to CircleCI"
	disconnectTrigger  = "disconnect"
	disconnectHelpText = "Disconnect your Mattermost account from CircleCI"

	notConnectedText    = "You are not connected to CircleCI. Please try `/" + commandTrigger + " " + connectTrigger + "`"
	errorConnectionText = "Error when reaching to CircleCI. Please check that your token is still valid"

	helpText = "## CircleCI plugin\n" +
		"* `/" + commandTrigger + " " + meTrigger + "` - " + meHelpText + "\n" +
		"* `/" + commandTrigger + " " + connectTrigger + " " + connectHint + "` - " + connectHelpText + "\n" +
		"* `/" + commandTrigger + " " + disconnectTrigger + "` - " + disconnectHelpText + "\n"
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

	me := model.NewAutocompleteData(meTrigger, "", meHelpText)
	mainCommand.AddCommand(me)

	connect := model.NewAutocompleteData(connectTrigger, connectHint, connectHelpText)
	connect.AddTextArgument("Generate a Personal API Token from your CircleCI user settings", connectHint, "")
	mainCommand.AddCommand(connect)

	disconnect := model.NewAutocompleteData(disconnectTrigger, "", disconnectHelpText)
	mainCommand.AddCommand(disconnect)

	return mainCommand
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)

	command := "help"
	if 1 < len(split) {
		command = split[1]
	}

	switch command {
	case meTrigger:
		return p.executeMe(args)

	case connectTrigger:
		return p.executeConnect(args, split[2:])

	case disconnectTrigger:
		return p.executeDisconnect(args)

	default:
		return p.sendEphemeralResponse(args, helpText), nil
	}
}
