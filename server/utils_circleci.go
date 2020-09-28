package main

import "github.com/jszwedko/go-circleci"

// Return false if an error has occurred
func (p *Plugin) getCircleUserInfo(circleToken string) (*circleci.User, bool) {
	circleClient := &circleci.Client{
		Token: circleToken,
	}

	user, err := circleClient.Me()
	if err != nil {
		p.API.LogError("Error when reaching CircleCI", "CircleCI error:", err)
		return nil, false
	}

	return user, true
}
