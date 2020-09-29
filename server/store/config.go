package store

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

// Config configuration for the plugin
type Config struct {
	VcsType string
	Org     string
	Project string
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
		s.api.LogError("Unable to save config", err)
		return nil, errors.Wrap(err, "Unable to save config")
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
