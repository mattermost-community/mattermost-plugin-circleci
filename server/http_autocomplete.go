package main

import (
	"net/http"

	"github.com/jszwedko/go-circleci"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
)

func httpAutocompleteFollowedProject(p *Plugin, w http.ResponseWriter, r *http.Request, circleciToken string) {
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
