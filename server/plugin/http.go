package plugin

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

const (
	routeWebhooks = "/hooks"

	routeAutocomplete        = "/autocomplete"
	subrouteFollowedProjects = "/followedProjects"
	routeApporveJob          = "/job/approve"
	routeEnvOverwrite        = "/env/overwrite"
)

func (p *Plugin) initializeRouter() {
	p.router = mux.NewRouter()

	p.router.HandleFunc(routeWebhooks+"/{secret}", p.httpHandleWebhook).Methods("POST")

	autocompleteRouter := p.router.PathPrefix(routeAutocomplete).Subrouter()
	autocompleteRouter.HandleFunc(subrouteFollowedProjects, p.autocompleteFollowedProject).Methods("GET")
	p.router.HandleFunc(routeApporveJob, p.httpHandleApprove).Methods("POST")
	p.router.HandleFunc(routeEnvOverwrite, p.httpHandleEnvOverwrite).Methods("POST")
}

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("Request received", "URL", r.URL)
	p.router.ServeHTTP(w, r)
}

// overwrite given env var after confirmation if already exist
func (p *Plugin) httpHandleEnvOverwrite(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, err := p.Store.GetTokenForUser(userID, p.getConfiguration().EncryptionKey)
	if err != nil {
		p.API.LogError("Error when getting token", err)
	}

	if circleciToken == "" {
		http.NotFound(w, r)
	}

	requestData := model.PostActionIntegrationRequestFromJson(r.Body)
	if requestData == nil {
		p.API.LogError("Empty request data", "request", r)
		return
	}

	responsePost := &model.Post{
		Id:        requestData.PostId,
		ChannelId: requestData.ChannelId,
		RootId:    requestData.PostId,
		UserId:    p.botUserID,
	}

	action := fmt.Sprintf("%s", requestData.Context["Action"])
	projectSlug := fmt.Sprintf("%v", requestData.Context["ProjectSlug"])
	name := fmt.Sprintf("%v", requestData.Context["EnvName"])
	val := fmt.Sprintf("%v", requestData.Context["EnvVal"])

	switch action {
	case "deny":
		responsePost.Message = fmt.Sprintf("Did not overwrite env variable %s for project %s", name, projectSlug)
		p.API.UpdateEphemeralPost(userID, responsePost)

	case "approve":
		if err := circle.AddEnvVar(circleciToken, projectSlug, name, val); err != nil {
			p.API.LogError("Error occurred while adding environment variable", err)
			responsePost.Message = fmt.Sprintf(":red_circle: Could not overwrite env var %s:%s from Mattermost.", name, val)
		} else {
			responsePost.Message = fmt.Sprintf(":white_check_mark: Successfully added environment variable `%s=%s` for project %s", name, val, projectSlug)
		}
		p.API.UpdateEphemeralPost(userID, responsePost)

	default:
		p.API.LogError("action %s is not valid", action)
	}
}
