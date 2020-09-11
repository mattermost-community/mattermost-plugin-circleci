package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	flagOnlyFailedBuilds = "only-failed"
)

type SubscriptionFlags struct {
	OnlyFailedBuilds bool `json:"OnlyFailedBuilds"`
}

func (s *SubscriptionFlags) AddFlag(flag string) error {
	switch flag { // nolint:gocritic // It's expected that more flags get added.
	case flagOnlyFailedBuilds:
		s.OnlyFailedBuilds = true

	default:
		return errors.New("Unknown flag " + flag)
	}

	return nil
}

func (s SubscriptionFlags) String() string {
	flags := []string{}

	if s.OnlyFailedBuilds {
		flag := "--" + flagOnlyFailedBuilds
		flags = append(flags, flag)
	}

	return strings.Join(flags, ",")
}

type Subscription struct {
	ChannelID  string            `json:"ChannelID"`
	CreatorID  string            `json:"CreatorID"`
	Flags      SubscriptionFlags `json:"Flags"`
	Owner      string            `json:"Owner"`
	Repository string            `json:"Repository"`
}

func (s *Subscription) ToSlackAttachmentField(p *Plugin) *model.SlackAttachmentField {
	username := "Unknown user"
	if user, appErr := p.API.GetUser(s.CreatorID); appErr != nil {
		p.API.LogError("Unable to get username", "userID", s.CreatorID)
	} else {
		username = user.Username
	}

	return &model.SlackAttachmentField{
		Title: getFullNameFromOwnerAndRepo(s.Owner, s.Repository),
		Short: true,
		Value: fmt.Sprintf(
			"Subscribed by: @%s\nFlags: ` %s`",
			username,
			s.Flags.String(),
		),
	}
}
