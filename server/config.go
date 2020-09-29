package commands

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	// ConfigCommandTrigger trigger for the command
	ConfigCommandTrigger = "config"
	hint                 = "[vcs/org-name/project-name]"
	helpText             = "View the config. Pass in the project (vcs/org/projectname) to set the default config"
)

// GetConfigAutoCompeleteData returns the auto complete info
func getConfigAutoCompeleteData() *model.AutocompleteData {
	configCommand := model.NewAutocompleteData(ConfigCommandTrigger, hint, helpText)
	configCommand.AddTextArgument("project identifier. (vcs/org/projectname)", "[project identifier]", "")
	return configCommand
}

// ExecuteConfigCommand executes the config command
func (p *Plugin) executeConfigCommand(args *model.CommandArgs, db store.Store) string {
	commandArgs := strings.Fields(args.Command)
	projectSlug := ""
	if len(commandArgs) > 2 {
		projectSlug = commandArgs[2]
	}

	if projectSlug == "" {
		return getConfig(args.UserId, db)
	}

	slug := strings.Split(projectSlug, "/")

	if len(slug) != 3 {
		return ":red_circle: Project should be specified in the format `vcs/orgname/projectname`. ex: `gh/mattermost/mattermost-server`"
	}
	if slug[0] != "gh" && slug[0] != "bb" {
		return ":red_circle: Invalid vcs value. Vcs should be either `gh` or `bb`. Example `gh/mattermost/mattermost-server`"
	}

	defaultConfig := &store.Config{
		VcsType: slug[0],
		Org:     slug[1],
		Project: slug[2],
	}

	return setConfig(args.UserId, *defaultConfig, db)
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

	return fmt.Sprintf(":white_check_mark: Successfully saved config. Vcs %s Org %s, Project %s as your default", config.VcsType, config.Org, config.Project)
}
