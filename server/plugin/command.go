package plugin

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
	commandAutocompleteHint = "<" + accountTrigger + "|" + projectTrigger + "|" + subscribeTrigger + "|" + workflowTrigger + ">"
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

	subscriptionHelp = "#### Subscribe to notifications projects\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeListTrigger + " " + subscribeListHint + "` — " + subscribeListHelpText + "\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeChannelTrigger + " " + subscribeChannelHint + "` — " + subscribeChannelHelpText + "\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeUnsubscribeChannelTrigger + " " + subscribeUnsubscribeChannelHint + "` — " + subscribeUnsubscribeChannelHelpText + "\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeListAllChannelsTrigger + " " + subscribeListAllChannelsHint + "` — " + subscribeListAllChannelsHelpText + "\n"

	help = "## CircleCI plugin Help\n" + accountHelp + projectHelp + subscriptionHelp
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

	// Add all subcommands
	mainCommand.AddCommand(getAccountAutoCompleteData())
	mainCommand.AddCommand(getProjectAutoComplete())
	mainCommand.AddCommand(getSubscribeAutoCompleteData())
	mainCommand.AddCommand(getConfigAutoCompleteData())
	mainCommand.AddCommand(getWorkflowAutoCompeleteData())
	mainCommand.AddCommand(getPipelineAutoCompeleteData())

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

	case subscribeTrigger:
		message += subscriptionHelp

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

	case subscribeTrigger:
		return p.executeSubscribe(args, token, split[2:])

	case configCommandTrigger:
		return p.executeConfig(args)

	case workflowTrigger:
		return p.executeWorkflowTrigger(args, token, split[2:])

	case pipelineTrigger:
		return p.executePipelineTrigger(args, token, split[2:])

	case commandHelpTrigger:
		return p.sendHelpResponse(args, "")

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

	circleToken, exists := p.Store.GetTokenForUser(userID)
	if !exists {
		return "", true
	}

	return circleToken, false
}
