package commands

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	configCommandTrigger = "config"
	configHint           = "<org-name/project-name>"
	configHelpText       = "View the config. Pass in the project (org-name/project-name) to set the default project"
)

// getConfigAutoCompeleteData returns the auto complete info
func getConfigAutoCompeleteData() *model.AutocompleteData {
	configCommand := model.NewAutocompleteData(configCommandTrigger, configHint, configHelpText)
	configCommand.AddTextArgument("The project identifier (org-name/project-name)", configHint, ".*/.*")
	return configCommand
}

// executeConfigCommand executes the config command
func executeConfigCommand(args *model.CommandArgs, db store.Store) string {
	commandArgs := strings.Fields(args.Command)
	projectSlug := ""
	if len(commandArgs) > 2 {
		projectSlug = commandArgs[2]
	}

	if projectSlug == "" {
		return getConfig(args.UserId, db)
	}

	slug := strings.Split(projectSlug, "/")

	if len(slug) != 2 {
		return ":red_circle: Project should be specified in the format orgname/projectname. ex: mattermost/mattermost-server"
	}
	defaultConfig := &store.Config{
		Org:     slug[0],
		Project: slug[1],
	}

	return setConfig(args.UserId, *defaultConfig, db)
}

func getConfig(userID string, db store.Store) string {
	savedConfig, _ := db.GetConfig(userID)
	if savedConfig != nil {
		return fmt.Sprintf(":information_source: Organization: %s, Project: %s", savedConfig.Org, savedConfig.Project)
	}
	return ":red_circle: No config exists. use `/circleci config <org-name/project-name>` to set the default project"
}

func setConfig(userID string, config store.Config, db store.Store) string {
	if err := db.SaveConfig(userID, config); err != nil {
		fmt.Println("error occurred while saving")
	}

	return fmt.Sprintf(":white_check_mark: Successfully saved config. Org %s, Project %s as your default", config.Org, config.Project)
}
