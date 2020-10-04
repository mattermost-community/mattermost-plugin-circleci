package store

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	FlagOnlyFailedBuilds = "only-failed"
)

// Contain the options for subscriptions
type SubscriptionFlags struct {
	OnlyFailedBuilds bool `json:"OnlyFailedBuilds"`
}

// Add a flag to the structure
func (s *SubscriptionFlags) AddFlag(flag string) error {
	switch flag { // nolint:gocritic // It's expected that more flags get added.
	case FlagOnlyFailedBuilds:
		s.OnlyFailedBuilds = true

	default:
		return errors.New("Unknown flag " + flag)
	}

	return nil
}

// Output the flags in a well-formatted string
func (s SubscriptionFlags) String() string {
	flags := []string{}

	if s.OnlyFailedBuilds {
		flag := "--" + FlagOnlyFailedBuilds
		flags = append(flags, flag)
	}

	if len(flags) == 0 {
		return "No flags set"
	}

	return strings.Join(flags, ",")
}

type Subscription struct {
	ChannelID          string            `json:"ChannelID"`
	CreatorID          string            `json:"CreatorID"`
	Flags              SubscriptionFlags `json:"Flags"`
	ProjectInformation ProjectIdentifier `json:"ProjectInformation"`
}

// Store the subscriptions.
// Keys of the map are projects slugs, in format (gh|bb)/org-name/project-name
// Values of the map are arrays of subscriptions, with differents channels, flags and creator ids
type Subscriptions struct {
	Repositories map[string][]*Subscription
}

// Transform the subscription to a well-formatted short slack attachment field
func (s *Subscription) ToSlackAttachmentField(username string) *model.SlackAttachmentField {
	if username == "" {
		username = s.CreatorID
	}

	return &model.SlackAttachmentField{
		Title: s.ProjectInformation.ToSlug(),
		Short: true,
		Value: fmt.Sprintf(
			"Subscribed by: @%s\nFlags: `%s`",
			username,
			s.Flags.String(),
		),
	}
}

// AddSubscription adds a new subscription in the struct
// Return true if the subscription was already existing and has been updated
func (s *Subscriptions) AddSubscription(newSub *Subscription) bool {
	key := newSub.ProjectInformation.ToSlug()

	repoSubs := s.Repositories[key]

	if repoSubs == nil {
		s.Repositories[key] = []*Subscription{newSub}
		return false
	}

	exists := false
	for index, sub := range repoSubs {
		if sub.ChannelID == newSub.ChannelID {
			// Replace the existing subscriptions
			repoSubs[index] = newSub
			exists = true
			break
		}
	}

	if !exists {
		s.Repositories[key] = append(repoSubs, newSub)
	}

	return exists
}

// RemoveSubscription removes a subscription from the struct
// Return true if the subscription has been found and removed
func (s *Subscriptions) RemoveSubscription(channelID string, conf *ProjectIdentifier) bool {
	key := conf.ToSlug()

	repoSubs := s.Repositories[key]
	if repoSubs == nil {
		return false
	}

	removed := false
	for index, sub := range repoSubs {
		if sub.ChannelID == channelID {
			repoSubs = append(repoSubs[:index], repoSubs[index+1:]...)
			removed = true
			break
		}
	}

	if removed {
		if len(repoSubs) == 0 {
			delete(s.Repositories, key)
		} else {
			s.Repositories[key] = repoSubs
		}

		return true
	}

	return false
}

// Get the subscriptions for a given channel
func (s *Subscriptions) GetSubscriptionsByChannel(channelID string) []*Subscription {
	var filteredSubs []*Subscription

	for _, v := range s.Repositories {
		for _, sub := range v {
			if sub.ChannelID == channelID {
				filteredSubs = append(filteredSubs, sub)
			}
		}
	}

	sort.Slice(filteredSubs, func(i, j int) bool {
		return filteredSubs[i].ProjectInformation.Project < filteredSubs[j].ProjectInformation.Project
	})

	return filteredSubs
}

// Get all the subscriptions for a given project
func (s *Subscriptions) GetSubscriptionsForProject(conf *ProjectIdentifier) []*Subscription {
	key := conf.ToSlug()
	return s.Repositories[key]
}

// Return a list of subscribed channel IDs for a project
func (s *Subscriptions) GetSubscribedChannelsForProject(conf *ProjectIdentifier) []string {
	subs := s.GetSubscriptionsForProject(conf)
	if subs == nil {
		return nil
	}

	var channelIDs []string
	for _, sub := range subs {
		channelIDs = append(channelIDs, sub.ChannelID)
	}

	return channelIDs
}

// Get all the channels concerned by a job for a project, filtered with subscription flags.
func (s *Subscriptions) GetFilteredChannelsForJob(conf *ProjectIdentifier, isFailed bool) []string {
	subs := s.GetSubscriptionsForProject(conf)
	if subs == nil {
		return nil
	}

	var channelIDs []string
	for _, sub := range subs {
		switch { // nolint:gocritic // It's expected that more flags get added.
		case isFailed || !sub.Flags.OnlyFailedBuilds:
			channelIDs = append(channelIDs, sub.ChannelID)
		}
	}

	return channelIDs
}
