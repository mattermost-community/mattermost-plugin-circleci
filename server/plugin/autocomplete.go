package plugin

import (
	"net/http"

	"github.com/jszwedko/go-circleci"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) autocompleteFollowedProject(w http.ResponseWriter, r *http.Request) {
	// Check token
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.Store.GetTokenForUser(userID, p.getConfiguration().EncryptionKey)
	if !exists {
		http.NotFound(w, r)
	}

	if r.Method != http.MethodGet {
		p.respondAndLogErr(w, http.StatusMethodNotAllowed, errors.New("method"+r.Method+"is not allowed, must be GET"))
		return
	}

	circleciClient := &circleci.Client{Token: circleciToken}
	projects, err := circleciClient.ListProjects()
	if err != nil {
		p.respondAndLogErr(w, http.StatusInternalServerError, err)
		return
	}

	out := []model.AutocompleteListItem{
		{
			HelpText: "Manually type the project's VCS repository name",
			Item:     "[repository]",
		},
	}
	if len(projects) == 0 {
		p.respondJSON(w, out)
		return
	}

	for _, project := range projects {
		out = append(out, model.AutocompleteListItem{
			HelpText: project.VCSURL,
			Item:     project.Reponame,
		})
	}
	p.respondJSON(w, out)
}
