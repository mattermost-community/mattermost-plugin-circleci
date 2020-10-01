package main

import (
	"github.com/mattermost/mattermost-server/v5/plugin"

	local "github.com/nathanaelhoun/mattermost-plugin-circleci/server/plugin"
)

func main() {
	plugin.ClientMain(&local.Plugin{})
}
