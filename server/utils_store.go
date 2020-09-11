package main

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	KVStoreSuffix = "_circleci_token"
)

// Return false if no token is saved for this user
func (p *Plugin) getTokenKV(userID string) (string, bool) {
	raw, appErr := p.API.KVGet(userID + KVStoreSuffix)
	if appErr != nil {
		p.API.LogError("Unable to reach KVStore", "KVStore error", appErr)
		return "", false
	}

	if raw == nil {
		return "", false
	}

	userToken := string(raw)
	return userToken, true
}

// Return false if the token has not been saved
func (p *Plugin) storeTokenKV(userID string, circleciToken string) bool {
	appErr := p.API.KVSet(userID+KVStoreSuffix, []byte(circleciToken))
	if appErr != nil {
		p.API.LogError("Unable to write in KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

// Return false if the token has not been deleted
func (p *Plugin) deleteTokenKV(userID string) bool {
	if appErr := p.API.KVDelete(userID + KVStoreSuffix); appErr != nil {
		p.API.LogError("Unable to delete from KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

func (p *Plugin) getSubscriptionsKV() (*Subscriptions, error) {
	var subscriptions *Subscriptions

	value, appErr := p.API.KVGet(subscriptionsKVKey)
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

func (p *Plugin) storeSubscriptionsKV(s *Subscriptions) error {
	b, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "error while converting subscriptions map to json")
	}

	if appErr := p.API.KVSet(subscriptionsKVKey, b); appErr != nil {
		return errors.Wrap(appErr, "could not store subscriptions in KV store")
	}

	return nil
}
