package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	configCommandTrigger  = "config"
	configCommandHint     = "[vcs/org-name/project-name]"
	configCommandHelpText = "View the config. Pass in the project (vcs/org/projectname) to set the default config"
)

func getConfigAutoCompleteData() *model.AutocompleteData {
	configCommand := model.NewAutocompleteData(configCommandTrigger, configCommandHint, configCommandHelpText)
	configCommand.AddTextArgument("project identifier. (vcs/org-name/project-name)", "[project identifier]", "")
	return configCommand
}

// ExecuteConfigCommand executes the config command
func (p *Plugin) executeConfig(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	commandArgs := strings.Fields(args.Command)
	projectSlug := ""
	if len(commandArgs) > 2 {
		projectSlug = commandArgs[2]
	}

	if projectSlug == "" {
		return p.sendEphemeralResponse(args, getConfig(args.UserId, p.Store)), nil
	}

	defaultConfig, userErr := store.CreateConfigFromSlug(projectSlug)
	if userErr != "" {
		return p.sendEphemeralResponse(args, userErr), nil
	}

	result := setConfig(args.UserId, *defaultConfig, p.Store)
	return p.sendEphemeralResponse(args, result), nil
}

func getConfig(userID string, db store.Store) string {
	savedConfig, _ := db.GetConfig(userID)
	if savedConfig != nil {
		return fmt.Sprintf(":information_source: Organization: %s, Project: %s", savedConfig.Org, savedConfig.Project)
	}
	return ":red_circle: No config exists. use `/circleci config vcs/orgname/projectname` to set the default project"
}

func setConfig(userID string, config store.Config, db store.Store) string {
	if err := db.SaveConfig(userID, config); err != nil {
		fmt.Println("error occurred while saving")
	}

	return fmt.Sprintf(":white_check_mark: Successfully saved config. Vcs %s Org %s, Project %s as your default", config.VCSType, config.Org, config.Project)
}
