package main

import (
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	botUserName    = "circleci"
	botDisplayName = "CircleCI"
	botDescription = "Created by the CircleCI Plugin"
)

var (
	badgeFailedURL     string
	badgePassedURL     string
	buildFailedIconURL string
	buildGreenIconURL  string
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botUserID string
}

func (p *Plugin) OnActivate() error {
	URLPluginStaticBase := "/plugins/" + manifest.Id + "/public/" // TODO add siteURL ?
	badgeFailedURL = URLPluginStaticBase + "circleci-failed.svg"
	badgePassedURL = URLPluginStaticBase + "circleci-passed.svg"
	buildFailedIconURL = URLPluginStaticBase + "circleci-build-fail.png"
	buildGreenIconURL = URLPluginStaticBase + "circleci-build-green.png"

	// Create bot user
	botUserID, err := p.Helpers.EnsureBot(
		&model.Bot{
			Username:    botUserName,
			DisplayName: botDisplayName,
			Description: botDescription,
		},
		plugin.ProfileImagePath("/assets/circleci.png"),
	)
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot")
	}
	p.botUserID = botUserID

	// Register slash command
	if err := p.API.RegisterCommand(p.getCommand()); err != nil {
		return errors.Wrap(err, "failed to register new command")
	}

	return nil
}
