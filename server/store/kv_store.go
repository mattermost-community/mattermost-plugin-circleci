package store

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	KVStoreSuffix = "_circleci_token"
)

// Return false if no token is saved for this user
func (s *Store) GetTokenForUser(userID string) (string, bool) {
	raw, appErr := s.api.KVGet(userID + KVStoreSuffix)
	if appErr != nil {
		s.api.LogError("Unable to reach KVStore", "KVStore error", appErr)
		return "", false
	}

	if raw == nil {
		return "", false
	}

	userToken := string(raw)
	return userToken, true
}

// Return false if the token has not been saved
func (s *Store) StoreTokenForUser(userID string, circleciToken string) bool {
	appErr := s.api.KVSet(userID+KVStoreSuffix, []byte(circleciToken))
	if appErr != nil {
		s.api.LogError("Unable to write in KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

// Return false if the token has not been deleted
func (s *Store) DeleteTokenForUser(userID string) bool {
	if appErr := s.api.KVDelete(userID + KVStoreSuffix); appErr != nil {
		s.api.LogError("Unable to delete from KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

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
