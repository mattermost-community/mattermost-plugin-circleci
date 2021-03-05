package store

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	storeTokenPrefix          = "circleci_token_"  // Full key format is circleci_token_userID
	defaultProjectStorePrefix = "default_project_" // Full key format is default_project_userID
	subscriptionsKVKey        = "subscriptions"
)

// GetTokenForUser returns false if no token is saved for this user
func (s *Store) GetTokenForUser(userID, encryptionKey string) (string, bool) {
	raw, appErr := s.api.KVGet(storeTokenPrefix + userID)
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

// StoreTokenForUser returns false if the token has not been saved
func (s *Store) StoreTokenForUser(userID, circleciToken, encryptionKey string) bool {
	encryptedToken, err := encrypt([]byte(encryptionKey), circleciToken)
	if err != nil {
		s.api.LogError("Error occurred while encrypting access token", "error", err)
		return false
	}

	appErr := s.api.KVSet(storeTokenPrefix+userID, []byte(encryptedToken))
	if appErr != nil {
		s.api.LogError("Unable to write in KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

// DeleteTokenForUser return sfalse if the token has not been deleted
func (s *Store) DeleteTokenForUser(userID string) bool {
	if appErr := s.api.KVDelete(storeTokenPrefix + userID); appErr != nil {
		s.api.LogError("Unable to delete from KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

// GetSubscriptions returns all the subscriptions from the KVStore
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

// StoreSubscriptions stores the subscriptions in the KVStore
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

// GetDefaultProject retrieves the saved default project for the user. Returns nil incase no default project exists for the user
func (s *Store) GetDefaultProject(userID string) (*ProjectIdentifier, error) {
	var pi *ProjectIdentifier

	savedDefaultProject, err := s.api.KVGet(defaultProjectStorePrefix + userID)
	if err != nil {
		s.api.LogError("Unable to get default project", err)
		return nil, errors.Wrap(err, "Unable to get default project")
	}

	if savedDefaultProject == nil {
		return nil, nil
	}
	appError := json.NewDecoder(bytes.NewReader(savedDefaultProject)).Decode(&pi)
	if appError != nil {
		return nil, errors.Wrap(appError, "could not properly decode saved default project")
	}

	return pi, nil
}

// StoreDefaultProject saves the passed in default project
func (s *Store) StoreDefaultProject(userID string, project ProjectIdentifier) error {
	projectBytes, err := json.Marshal(project)
	if err != nil {
		return errors.Wrap(err, "error while converting project to json")
	}

	if err := s.api.KVSet(defaultProjectStorePrefix+userID, projectBytes); err != nil {
		s.api.LogError("Unable to save default project", err)
		return errors.Wrap(err, "Unable to save default project")
	}
	return nil
}
