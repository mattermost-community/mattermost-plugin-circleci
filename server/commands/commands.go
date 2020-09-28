package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	mainTrigger  = "circleci"
	mainHint     = "<" + accountTrigger + ">"
	mainHelpText = "Interact with CircleCI jobs and builds"

	helpTrigger = "help"
	helpMessage = "## CircleCI plugin Help\n"

	notConnectedText    = "You are not connected to CircleCI. Please try `/" + mainTrigger + " " + accountTrigger + " " + accountConnectTrigger + "`"
	errorConnectionText = "Error when reaching to CircleCI. Please check that your token is still valid"
)

// bundlePath is the absolute path where the plugin's bundle was unpacked by the Mattermost server
func GetCommand(bundlePath string) (*model.Command, error) {
	autocompleteIconData, err := getAutocompleteIconData(bundlePath)

	return &model.Command{
		Trigger:              mainTrigger,
		AutoComplete:         true,
		AutoCompleteDesc:     mainHelpText,
		AutoCompleteHint:     mainHint,
		AutocompleteData:     getAutocompleteData(),
		AutocompleteIconData: autocompleteIconData,
	}, err
}

func getAutocompleteIconData(bundlePath string) (string, error) {
	icon, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "circleci.svg"))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(icon)), nil
}

func getAutocompleteData() *model.AutocompleteData {
	mainCommand := model.NewAutocompleteData(mainTrigger, mainHint, mainHelpText)

	mainCommand.AddCommand(getAccountAutocompleteData())
	mainCommand.AddCommand(getProjectAutocompleteData())
	mainCommand.AddCommand(getSubscribeAutocompleteData())
	mainCommand.AddCommand(getConfigAutoCompeleteData())

	return mainCommand
}

func ExecuteCommand(args *model.CommandArgs, db store.Store) string {
	split := strings.Fields(args.Command)

	command := ""
	if 1 < len(split) {
		command = split[1]
	}

	switch command {
	case accountTrigger:
		return executeAccount(args, token, split[2:])

	case configCommandTrigger:
		return executeConfigCommand(args, db)

	case projectTrigger:
		return executeProject(args, token, split[2:])

	case subscribeTrigger:
		return executeSubscribe(args, token, split[2:])

	case helpTrigger:
		return formatHelpMessage(args, "")

	default:
		return formatIncorrectSubcommand(args, "")
	}
}

func formatIncorrectSubcommand(args *model.CommandArgs, currentCommand string) string {
	subcommandText := ""
	if currentCommand != "" {
		subcommandText = " " + currentCommand
	}

	return fmt.Sprint("Invalid subcommand given. Type `/%s%s help` to get a hint", mainTrigger, subcommandText)
}

func formatHelpMessage(args *model.CommandArgs, currentCommand string) string {
	message := ""

	switch currentCommand {
	case accountTrigger:
		message += accountHelpMessage

	case projectTrigger:
		message += projectHelpMessage

	case subscribeTrigger:
		// No message at the moment

	default:
		message += helpMessage + accountHelpMessage + projectHelpMessage
	}

	return message
}
