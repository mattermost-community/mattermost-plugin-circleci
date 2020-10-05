package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
)

const (
	setDefaultCommandTrigger  = "default"
	setDefaultCommandHint     = "[vcs/org-name/project-name]"
	setDefaultCommandHelpText = "View your current default project or pass in a new project (vcs/org-name/project-name)"
)

func getSetDefaultAutoCompleteData() *model.AutocompleteData {
	setDefaultCommand := model.NewAutocompleteData(setDefaultCommandTrigger, setDefaultCommandHint, setDefaultCommandHelpText)
	setDefaultCommand.AddTextArgument("project identifier. (vcs/org-name/project-name)", "<project identifier>", namedArgProjectPattern)
	return setDefaultCommand
}

// executeSetDefault executes the setDefault command
func (p *Plugin) executeSetDefault(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	commandArgs := strings.Fields(args.Command)
	projectSlug := ""
	if len(commandArgs) > 2 {
		projectSlug = commandArgs[2]
	}

	if projectSlug == "" {
		return p.sendEphemeralResponse(args, getDefaultProject(args.UserId, p.Store)), nil
	}

	defaultProject, userErr := store.CreateProjectIdentifierFromSlug(projectSlug)
	if userErr != "" {
		return p.sendEphemeralResponse(args, userErr), nil
	}

	result := setDefaultProject(args.UserId, *defaultProject, p.Store)
	return p.sendEphemeralResponse(args, result), nil
}

func getDefaultProject(userID string, db store.Store) string {
	savedDefaultProfile, _ := db.GetDefaultProject(userID)
	if savedDefaultProfile != nil {
		return fmt.Sprintf(":information_source: Current default project: %s", savedDefaultProfile.ToMarkdown())
	}

	return ":red_circle: No default project is set. Use `/circleci set-default <vcs/org-name/project-name>` to set the default project"
}

func setDefaultProject(userID string, newDefaultProject store.ProjectIdentifier, db store.Store) string {
	if err := db.StoreDefaultProject(userID, newDefaultProject); err != nil {
		return ":red_circle: An error has occurred while saving your default project"
	}

	return fmt.Sprintf(":white_check_mark: Successfully saved %s as your default project", newDefaultProject.ToMarkdown())
}
