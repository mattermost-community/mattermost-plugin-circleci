package store

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

type Subscription struct {
	ChannelID  string            `json:"ChannelID"`
	CreatorID  string            `json:"CreatorID"`
	Flags      SubscriptionFlags `json:"Flags"`
	Owner      string            `json:"Owner"`
	Repository string            `json:"Repository"`
}

func (s *Subscription) ToSlackAttachmentField(username string) *model.SlackAttachmentField {
	return &model.SlackAttachmentField{
		Title: circle.GetFullNameFromOwnerAndRepo(s.Owner, s.Repository),
		Short: true,
		Value: fmt.Sprintf(
			"Subscribed by: @%s\nFlags: ` %s`",
			username,
			s.Flags.String(),
		),
	}
}
