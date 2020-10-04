package store

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	storeTokenSuffix          = "_circleci_token"  // Full key format is userID_circleci_token
	defaultProjectStoreSuffix = "_default_project" // Full key format is userID_default_project
	subscriptionsKVKey        = "subscriptions"
)

// Return false if no token is saved for this user
func (s *Store) GetTokenForUser(userID, encryptionKey string) (string, bool) {
	raw, appErr := s.api.KVGet(userID + storeTokenSuffix)
	if appErr != nil {
		s.api.LogError("Unable to reach KVStore", "KVStore error", appErr)
		return "", false
	}

	if raw == nil {
		return "", false
	}

	userToken, err := decrypt([]byte(encryptionKey), string(raw))
	if err != nil {
		s.api.LogWarn("Failed to decrypt access token", "error", err)
		return "", false
	}

	return userToken, true
}

// Return false if the token has not been saved
func (s *Store) StoreTokenForUser(userID, circleciToken, encryptionKey string) bool {
	encryptedToken, err := encrypt([]byte(encryptionKey), circleciToken)
	if err != nil {
		s.api.LogError("Error occurred while encrypting access token", "error", err)
		return false
	}

	appErr := s.api.KVSet(userID+storeTokenSuffix, []byte(encryptedToken))
	if appErr != nil {
		s.api.LogError("Unable to write in KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

// Return false if the token has not been deleted
func (s *Store) DeleteTokenForUser(userID string) bool {
	if appErr := s.api.KVDelete(userID + storeTokenSuffix); appErr != nil {
		s.api.LogError("Unable to delete from KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

// Return all the subscriptions from the KVStore
func (s *Store) GetSubscriptions() (*Subscriptions, error) {
	var subscriptions *Subscriptions

	value, appErr := s.api.KVGet(subscriptionsKVKey)
	if appErr != nil {
		return nil, errors.Wrap(appErr, "could not get subscriptions from KVStore")
	}

	if value == nil {
		return &Subscriptions{Repositories: map[string][]*Subscription{}}, nil
	}

	err := json.NewDecoder(bytes.NewReader(value)).Decode(&subscriptions)
	if err != nil {
		return nil, errors.Wrap(err, "could not properly decode subscriptions key")
	}

	return subscriptions, nil
}

// Store the subscriptions in the KVStore
func (s *Store) StoreSubscriptions(subs *Subscriptions) error {
	b, err := json.Marshal(subs)
	if err != nil {
		return errors.Wrap(err, "error while converting subscriptions map to json")
	}

	if appErr := s.api.KVSet(subscriptionsKVKey, b); appErr != nil {
		return errors.Wrap(appErr, "could not store subscriptions in KV store")
	}

	return nil
}

// GetDefaultProjectConfig retrieves the saved config for the user. Returns nil incase no config exists for the user
func (s *Store) GetDefaultProjectConfig(userID string) (*ProjectIdentifier, error) {
	var pi *ProjectIdentifier

	savedConfig, err := s.api.KVGet(userID + defaultProjectStoreSuffix)
	if err != nil {
		s.api.LogError("Unable to get config", err)
		return nil, errors.Wrap(err, "Unable to get config")
	}

	if savedConfig == nil {
		return nil, nil
	}
	appError := json.NewDecoder(bytes.NewReader(savedConfig)).Decode(&pi)
	if appError != nil {
		return nil, errors.Wrap(appError, "could not properly decode saved config")
	}

	return pi, nil
}

// StoreDefaultProjectConfig saves the passed in config
func (s *Store) StoreDefaultProjectConfig(userID string, config ProjectIdentifier) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "error while converting config to json")
	}

	if err := s.api.KVSet(userID+defaultProjectStoreSuffix, configBytes); err != nil {
		s.api.LogError("Unable to save config", err)
		return errors.Wrap(err, "Unable to save config")
	}
	return nil
}
