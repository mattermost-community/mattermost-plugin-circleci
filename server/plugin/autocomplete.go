package plugin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jszwedko/go-circleci"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) autocompleteFollowedProject(w http.ResponseWriter, r *http.Request) {
	// Check token
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.Store.GetTokenForUser(userID, p.getConfiguration().EncryptionKey)
	if !exists {
		http.NotFound(w, r)
	}

	circleciClient := &circleci.Client{Token: circleciToken}
	projects, err := circleciClient.ListProjects()
	if err != nil {
		p.respondAndLogErr(w, http.StatusInternalServerError, err)
		return
	}

	out := []model.AutocompleteListItem{
		{
			HelpText: "Manually type the project identifier",
			Item:     "<vcs/org-name/project-name>",
		},
	}
	if len(projects) == 0 {
		p.respondJSON(w, out)
		return
	}

	for _, project := range projects {
		vcs := "gh"
		if strings.Contains(project.VCSURL, "https://bitbucket.org") {
			vcs = "bb"
		}
		out = append(out, model.AutocompleteListItem{
			HelpText: project.VCSURL,
			Item:     fmt.Sprintf("%s/%s/%s", vcs, project.Username, project.Reponame),
		})
	}
	p.respondJSON(w, out)
}
