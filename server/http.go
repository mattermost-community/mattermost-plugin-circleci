package main

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	routeAutocompleteFollowedProjects = "/autocomplete/followedProjects"
)

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	// Check token
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.getTokenFromKVStore(userID)
	if !exists {
		http.NotFound(w, r)
	}

	// Call the handler
	switch r.URL.Path {
	case routeAutocompleteFollowedProjects:
		p.httpAutocompleteFollowedProject(w, r, circleciToken)

	default:
		http.NotFound(w, r)
	}
}
