package plugin

import (
	"fmt"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
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
	Store store.Store

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botUserID string

	router *mux.Router
}

// OnActivate is ran when the plugin is activated
func (p *Plugin) OnActivate() error {
	URLPluginStaticBase := fmt.Sprintf("%s/plugins/%s/public", *p.API.GetConfig().ServiceSettings.SiteURL, manifest.Id)
	badgeFailedURL = URLPluginStaticBase + "/circleci-failed.svg"
	badgePassedURL = URLPluginStaticBase + "/circleci-passed.svg"
	buildFailedIconURL = URLPluginStaticBase + "/circleci-build-fail.png"
	buildGreenIconURL = URLPluginStaticBase + "/circleci-build-green.png"

	st, err := store.NewStore(p.API)
	if err != nil {
		return errors.Wrap(err, "failed to create plugin store")
	}
	p.Store = st

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

	p.initializeRouter()

	// Register slash command
	if err := p.API.RegisterCommand(p.getCommand()); err != nil {
		return errors.Wrap(err, "failed to register new command")
	}

	return nil
}
