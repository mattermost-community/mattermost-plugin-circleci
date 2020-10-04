package store

import (
	"fmt"
	"sort"

	"github.com/mattermost/mattermost-server/v5/model"

	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"
)

const (
	subscriptionsKVKey = "subscriptions"
)

type Subscription struct {
	ChannelID string            `json:"ChannelID"`
	CreatorID string            `json:"CreatorID"`
	Flags     SubscriptionFlags `json:"Flags"`
	// TODO Add bitbucket
	// TODO rename Owner to Organization
	Owner string `json:"Owner"`
	// TODO rename Repository to Project
	Repository string `json:"Repository"`
}

type Subscriptions struct {
	Repositories map[string][]*Subscription
}

// Transform the subscription to a well-formatted short slack attachment field
func (s *Subscription) ToSlackAttachmentField(username string) *model.SlackAttachmentField {
	return &model.SlackAttachmentField{
		Title: fmt.Sprintf("gh/%s/%s", s.Owner, s.Repository), // TODO add support for bitbucket
		Short: true,
		Value: fmt.Sprintf(
			"Subscribed by: @%s\nFlags: `%s`",
			username,
			s.Flags.String(),
		),
	}
}

// AddSubscription adds a new subscription in the struct
func (s *Subscriptions) AddSubscription(newSub *Subscription) {
	key := v1.GetFullNameFromOwnerAndRepo(newSub.Owner, newSub.Repository)

	repoSubs := s.Repositories[key]

	if repoSubs == nil {
		s.Repositories[key] = []*Subscription{newSub}
		return
	}

	exists := false
	for index, sub := range repoSubs {
		if sub.ChannelID == newSub.ChannelID {
			repoSubs[index] = newSub
			exists = true
			break
		}
	}

	if !exists {
		s.Repositories[key] = append(repoSubs, newSub)
	}
}

// RemoveSubscription removes a subscription from the struct
// Return true if the subscription has been found and removed
func (s *Subscriptions) RemoveSubscription(channelID, owner, repository string) bool {
	key := v1.GetFullNameFromOwnerAndRepo(owner, repository)

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
		return filteredSubs[i].Repository < filteredSubs[j].Repository
	})

	return filteredSubs
}

func (s *Subscriptions) GetSubscriptionsForRepository(owner, repository string) []*Subscription {
	key := v1.GetFullNameFromOwnerAndRepo(owner, repository)
	return s.Repositories[key]
}

// Return a list of subscribed channel IDs
func (s *Subscriptions) GetSubscribedChannelsForRepository(owner, repository string) []string {
	subs := s.GetSubscriptionsForRepository(owner, repository)
	if subs == nil {
		return nil
	}

	var channelIDs []string
	for _, sub := range subs {
		channelIDs = append(channelIDs, sub.ChannelID)
	}

	return channelIDs
}

func (s *Subscriptions) GetFilteredChannelsForBuild(owner, repository string, isFailed bool) []string {
	subs := s.GetSubscriptionsForRepository(owner, repository)
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
