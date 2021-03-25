package store

import (
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Store is an interface to interact with the KV Store
type Store struct {
	api plugin.API
}

// NewStore returns a fresh store
func NewStore(api plugin.API) (Store, error) {
	store := Store{
		api: api,
	}
	return store, nil
}
