package plugin

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	commandTrigger          = "circleci"
	commandAutocompleteHint = "<" + accountTrigger + "|" + projectTrigger + "|" + subscribeTrigger + "|" + workflowTrigger + ">"
	commandAutocompleteDesc = "Interact with CircleCI jobs and builds"

	notConnectedText    = "You are not connected to CircleCI. Please try `/" + commandTrigger + " " + accountTrigger + " " + accountConnectTrigger + "`"
	errorConnectionText = ":red_circle: Error when reaching to CircleCI. Please check that your token is still valid"

	namedArgProjectName     = "project"
	namedArgProjectHelpText = "The project identifier. VCS is either gh or bb."
	namedArgProjectHint     = "[vcs/org-name/project-name]"
	namedArgProjectPattern  = "" // TODO wait for https://github.com/mattermost/mattermost-server/pull/14781 to get merged

	// All the Triggers and HelpTexts for the subcommands are defined in the corresponding .go file
	commandHelpTrigger = "help"

	accountHelp = "#### Connect to your CircleCI account\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountViewTrigger + "` — " + AccountViewHelpText + "\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountConnectTrigger + " " + accountConnectHint + "` — " + accountConnectHelpText + "\n" +
		"* `/" + commandTrigger + " " + accountTrigger + " " + accountDisconnectTrigger + "` — " + accountDisconnectHelpText + "\n"

	projectHelp = "#### Manage CircleCI projects\n" +
		"* `/" + commandTrigger + " " + projectTrigger + " " + projectListTrigger + "` — " + projectListHelpText + "\n" +
		"* `/" + commandTrigger + " " + projectTrigger + " " + projectRecentBuildsTrigger + " " + projectRecentBuildsHint + "` — " + projectRecentBuildsHelpText + "\n" +
		"* `/" + commandTrigger + " " + projectTrigger + " " + projectEnvVarTrigger + " " + projectEnvVarHint + "` — " + projectEnvVarHelpText + "\n"

	subscriptionHelp = "#### Subscribe your channel to notifications\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeListTrigger + " " + subscribeListHint + "` — " + subscribeListHelpText + "\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeChannelTrigger + " " + subscribeChannelHint + "` — " + subscribeChannelHelpText + "\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeUnsubscribeChannelTrigger + " " + subscribeUnsubscribeChannelHint + "` — " + subscribeUnsubscribeChannelHelpText + "\n" +
		"* `/" + commandTrigger + " " + subscribeTrigger + " " + subscribeListAllChannelsTrigger + " " + subscribeListAllChannelsHint + "` — " + subscribeListAllChannelsHelpText + "\n"

	workflowHelp = "#### Manage worflows\n" +
		"* `/" + commandTrigger + " " + workflowTrigger + " " + workflowGetTrigger + " " + workflowGetHint + "` — " + workflowGetHelpText + "\n" +
		"* `/" + commandTrigger + " " + workflowTrigger + " " + workflowGetJobsTrigger + " " + workflowGetJobsHint + "` — " + workflowGetJobsHelpText + "\n" +
		"* `/" + commandTrigger + " " + workflowTrigger + " " + workflowRerunTrigger + " " + workflowRerunHint + "` — " + workflowRerunHelpText + "\n" +
		"* `/" + commandTrigger + " " + workflowTrigger + " " + workflowCancelTrigger + " " + workflowCancelHint + "` — " + workflowCancelHelpText + "\n"

	pipelineHelp = "#### Manage pipelines\n" +
		"* `/" + commandTrigger + " " + pipelineTrigger + " " + pipelineTriggerTrigger + " " + pipelineTriggerHint + "` — " + pipelineTriggerHelpText + "\n" +
		"* `/" + commandTrigger + " " + pipelineTrigger + " " + pipelineWorkflowTrigger + " " + pipelineWorkflowHint + "` — " + pipelineWorkflowHelpText + "\n" +
		"* `/" + commandTrigger + " " + pipelineTrigger + " " + pipelineGetRecentTrigger + " " + pipelineGetRecentHint + "` — " + pipelineGetRecentHelpText + "\n" +
		"* `/" + commandTrigger + " " + pipelineTrigger + " " + pipelineGetAllTrigger + " " + pipelineGetAllHint + "` — " + pipelineGetAllHelpText + "\n" +
		"* `/" + commandTrigger + " " + pipelineTrigger + " " + pipelineGetMineTrigger + " " + pipelineGetMineHint + "` — " + pipelineGetMineHelpText + "\n" +
		"* `/" + commandTrigger + " " + pipelineTrigger + " " + pipelineGetSingleTrigger + " " + pipelineGetSingleHint + "` — " + pipelineGetSingleHelpText + "\n"

	setDefaultHelp = "#### Set your default project\n" +

		"* `/" + commandTrigger + " " + setDefaultCommandTrigger + " " + setDefaultCommandHint + "` — " + setDefaultCommandHelpText + "\n"

	help = "## CircleCI plugin Help\n" + accountHelp + setDefaultHelp + subscriptionHelp + pipelineHelp + workflowHelp + projectHelp
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
	mainCommand.AddCommand(getSetDefaultAutoCompleteData())
	mainCommand.AddCommand(getWorkflowAutoCompeleteData())
	mainCommand.AddCommand(getSubscribeAutoCompleteData())
	mainCommand.AddCommand(getProjectAutoComplete())
	mainCommand.AddCommand(getPipelineAutoCompeleteData())
	mainCommand.AddCommand(getInsightAutoCompeleteData())

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

	case setDefaultCommandTrigger:
		message += setDefaultHelp

	case workflowTrigger:
		message += workflowHelp

	case pipelineTrigger:
		message += pipelineHelp

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

	if command == commandHelpTrigger {
		return p.sendHelpResponse(args, "")
	}

	token, shouldBeConnected := getTokenIfConnected(p, split, args.UserId)
	if shouldBeConnected {
		return p.sendEphemeralResponse(args, notConnectedText), nil
	}

	defaultProject, err := p.Store.GetDefaultProject(args.UserId)
	if err != nil {
		p.API.LogError("Could not get user default project", "error", err)
	}

	var splitWithoutProject []string
	if (defaultProject != nil && !strings.Contains(args.Command, "--project")) || command == accountTrigger || command == setDefaultCommandTrigger {
		splitWithoutProject = split
	} else {
		// Trying to get the default project from the commands, with the args `--project`
		slug := ""
		nextIsValue := false
		splitWithoutProject = []string{}

	scan:
		for _, arg := range split {
			switch {
			case nextIsValue:
				slug = arg
				break scan

			case arg == "--project":
				nextIsValue = true

			default:
				splitWithoutProject = append(splitWithoutProject, arg)
			}
		}

		if slug == "" {
			// The argument has not been found
			return p.sendEphemeralResponse(args,
				fmt.Sprintf("No CircleCI project set. Try `%s %s %s` to set a project to use", commandTrigger, setDefaultCommandTrigger, setDefaultCommandHint),
			), nil
		}

		confFromArg, userErr := store.CreateProjectIdentifierFromSlug(slug)
		if userErr != "" {
			return p.sendEphemeralResponse(args,
				fmt.Sprintf("Incorrect value for argument `--project`: `%s`. %s", slug, userErr),
			), nil
		}

		defaultProject = confFromArg
	}

	switch command {
	case accountTrigger:
		return p.executeAccount(args, token, splitWithoutProject[2:])

	case projectTrigger:
		return p.executeProject(args, token, defaultProject, splitWithoutProject[2:])

	case subscribeTrigger:
		return p.executeSubscribe(args, token, defaultProject, splitWithoutProject[2:])

	case setDefaultCommandTrigger:
		return p.executeSetDefault(args)

	case workflowTrigger:
		return p.executeWorkflow(args, token, splitWithoutProject[2:])

	case pipelineTrigger:
		return p.executePipelineTrigger(args, token, defaultProject, splitWithoutProject[2:])

	case insightTrigger:
		return p.executeInsightTrigger(args, token, defaultProject, splitWithoutProject[2:])

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

	circleToken, exists := p.Store.GetTokenForUser(userID, p.getConfiguration().EncryptionKey)
	if !exists {
		return "", true
	}

	return circleToken, false
}
