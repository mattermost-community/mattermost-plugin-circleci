package store

import (
	"fmt"
	"strings"
)

// ProjectIdentifier structure to differentiate projects from differents VCS
type ProjectIdentifier struct {
	VCSType string // Currently supported : "gh" and "bb"
	Org     string
	Project string
}

// Return the slug in format "(gh|bb)/org-name/project-name)"
func (pi *ProjectIdentifier) ToSlug() string {
	return fmt.Sprintf("%s/%s/%s", pi.VCSType, pi.Org, pi.Project)
}

// Return a link to the repo formatted in Markdown
func (pi *ProjectIdentifier) ToMarkdown() string {
	VCSBaseURL := "https://github.com"
	if pi.VCSType == "bb" {
		VCSBaseURL = "https://bitbucket.org"
	}

	return fmt.Sprintf("[`%s/%s/%s`](%s/%s/%s)", pi.VCSType, pi.Org, pi.Project, VCSBaseURL, pi.Org, pi.Project)
}

// Return the link to CircleCI project page
func (pi *ProjectIdentifier) ToCircleURL() string {
	vcs := "github"
	if pi.VCSType == "bb" {
		vcs = "bitbucket"
	}

	return fmt.Sprintf("https://app.circleci.com/pipelines/%s/%s/%s", vcs, pi.Org, pi.Project)
}

// Create a ProjectIdentifier struct from a slug in format (gh|bb)/org-name/project-name)
// Return a string explaining the error is the format is not good
func CreateProjectIdentifierFromSlug(fullSlug string) (*ProjectIdentifier, string) {
	split := strings.Split(fullSlug, "/")

	if len(split) != 3 {
		return nil, ":red_circle: Project should be specified in the format `vcs/org-name/project-name`. Example: `gh/mattermost/mattermost-server`"
	}

	if split[0] != "gh" && split[0] != "bb" {
		return nil, ":red_circle: Invalid vcs value. VCS should be either `gh` or `bb`. Example: `gh/mattermost/mattermost-server`"
	}

	return &ProjectIdentifier{
		VCSType: split[0],
		Org:     split[1],
		Project: split[2],
	}, ""
}
