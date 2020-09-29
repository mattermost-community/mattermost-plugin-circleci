package v1

import (
	"fmt"

	"github.com/jszwedko/go-circleci"
)

var (
	badgeFailedURL string
	badgePassedURL string
)

// GetCircleUserInfo returns info about logged in user
func GetCircleUserInfo(circleToken string) (*circleci.User, error) {
	circleClient := &circleci.Client{
		Token: circleToken,
	}

	user, err := circleClient.Me()
	if err != nil {
		return nil, fmt.Errorf("Error when reaching CircleCI. CircleCI error: %s", err.Error())
	}

	return user, nil
}

// GetCircleciUserProjects returns projects for given user
func GetCircleciUserProjects(circleCiToken string) ([]*circleci.Project, error) {
	circleciClient := &circleci.Client{Token: circleCiToken}
	projects, err := circleciClient.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("Unable to get circleCI user projects. CircleCI API error: %s", err.Error())
	}

	return projects, nil
}
