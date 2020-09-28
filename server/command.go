package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/commands"
)

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	result := commands.ExecuteCommand(args, p.Store)
	return p.sendEphemeralResponse(args, result), nil
}
