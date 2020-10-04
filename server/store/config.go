package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Config configuration for the plugin
type Config struct {
	VCSType string
	Org     string
	Project string
}

// Return the slug in format "gh/mattermost/mattermost-server"
func (c *Config) ToSlug() string {
	return fmt.Sprintf("%s/%s/%s", c.VCSType, c.Org, c.Project)
}

// Return a link to the repo formatted in Markdown
func (c *Config) ToMarkdown() string {
	VCSBaseURL := "https://github.com"
	if c.VCSType == "bb" {
		VCSBaseURL = "https://bitbucket.org"
	}

	return fmt.Sprintf("[`%s/%s/%s`](%s/%s/%s)", c.VCSType, c.Org, c.Project, VCSBaseURL, c.Org, c.Project)
}

// Create a Config struct from a string in format (gh|bb)/org-name/project-name
func CreateConfigFromSlug(fullSlug string) (*Config, string) {
	split := strings.Split(fullSlug, "/")

	if len(split) != 3 {
		return nil, ":red_circle: Project should be specified in the format `vcs/org-name/project-name`. Example: `gh/mattermost/mattermost-server`"
	}

	if split[0] != "gh" && split[0] != "bb" {
		return nil, ":red_circle: Invalid vcs value. VCS should be either `gh` or `bb`. Example: `gh/mattermost/mattermost-server`"
	}

	return &Config{
		VCSType: split[0],
		Org:     split[1],
		Project: split[2],
	}, ""
}

const (
	configStoreSuffix = "_circleci_config"
)

// SaveConfig saves the passed in config
func (s *Store) SaveConfig(userID string, config Config) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "error while converting config to json")
	}

	if err := s.api.KVSet(userID+configStoreSuffix, configBytes); err != nil {
		s.api.LogError("Unable to save config", err)
		return errors.Wrap(err, "Unable to save config")
	}
	return nil
}

// GetConfig retrieves the saved config for the user. Returns nil incase no config exists for the user
func (s *Store) GetConfig(userID string) (*Config, error) {
	var config *Config

	savedConfig, err := s.api.KVGet(userID + configStoreSuffix)
	if err != nil {
		s.api.LogError("Unable to get config", err)
		return nil, errors.Wrap(err, "Unable to get config")
	}

	if savedConfig == nil {
		return nil, nil
	}
	appError := json.NewDecoder(bytes.NewReader(savedConfig)).Decode(&config)
	if appError != nil {
		return nil, errors.Wrap(appError, "could not properly decode saved config")
	}

	return config, nil
}
