package main

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	botUserName    = "circleci"
	botDisplayName = "CircleCI"
	botDescription = "Created by the CircleCI Plugin"

	botIconFile       = "circleci.png"
	botIconBuildFail  = "circleci-build-fail.png"
	botIconBuildGreen = "circleci-build-green.png"
	badgeFailedFile   = "circleci-failed.svg"
	badgePassedFile   = "circleci-passed.svg"
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

	iconCircleciURL   string
	iconBuildFailURL  string
	iconBuildGreenURL string
	badgeFailedURL    string
	badgePassedURL    string
}

func (p *Plugin) OnActivate() error {
	URLPluginStaticBase := "/plugins/" + manifest.Id + "/public/" // TODO add siteURL ?
	p.iconCircleciURL = URLPluginStaticBase + botIconFile
	p.iconBuildFailURL = URLPluginStaticBase + botIconBuildFail
	p.iconBuildGreenURL = URLPluginStaticBase + botIconBuildGreen
	p.badgeFailedURL = URLPluginStaticBase + badgeFailedFile
	p.badgePassedURL = URLPluginStaticBase + badgePassedFile

	// Create bot user
	botUserID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
		Description: botDescription,
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot")
	}
	p.botUserID = botUserID
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "failed to get bundle path")
	}

	profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "public", botIconFile))
	if err != nil {
		return errors.Wrap(err, "failed to read profile image")
	}

	if appErr := p.API.SetProfileImage(botUserID, profileImage); appErr != nil {
		return errors.Wrap(errors.New(appErr.Error()), "failed to set profile image")
	}

	// Register slash command
	if err := p.API.RegisterCommand(p.getCommand()); err != nil {
		return errors.Wrap(err, "failed to register new command")
	}

	return nil
}
