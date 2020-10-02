package plugin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	routeWebhooks = "/hooks"

	routeAutocomplete        = "/autocomplete"
	subrouteFollowedProjects = "/followedProjects"
)

func (p *Plugin) initializeRouter() {
	p.router = mux.NewRouter()

	p.router.HandleFunc(routeWebhooks+"/{secret}", p.httpHandleWebhook).Methods("POST")

	autocompleteRouter := p.router.PathPrefix(routeAutocomplete).Subrouter()
	autocompleteRouter.HandleFunc(subrouteFollowedProjects, p.autocompleteFollowedProject).Methods("GET")
}

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("Request received", "URL", r.URL)
	p.router.ServeHTTP(w, r)
}
