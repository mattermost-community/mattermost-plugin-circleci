package main

import (
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	routeAutocompleteFollowedProjects = "/autocomplete/followedProjects"
	routeWebhooksPrefix               = "/hooks"
)

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	// Check token
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.getTokenKV(userID)
	if !exists {
		http.NotFound(w, r)
	}

	routeWebhooks := strings.Join([]string{
		routeWebhooksPrefix,
		p.getConfiguration().WebhooksSecret,
	},
		"/",
	)

	p.API.LogDebug("Receveid CircleCI http request", "URL", r.URL.Path, "route for CircleCI Webhooks", routeWebhooks)

	// Call the handler
	switch r.URL.Path {
	case routeAutocompleteFollowedProjects:
		httpAutocompleteFollowedProject(p, w, r, circleciToken)
		return

	case routeWebhooks:
		httpHandleWebhook(p, w, r)
		return

	default:
		http.NotFound(w, r)
	}
}
