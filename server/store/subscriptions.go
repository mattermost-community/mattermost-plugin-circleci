package store

import (
	"sort"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

const (
	subscriptionsKVKey = "subscriptions"
)

type Subscriptions struct {
	Repositories map[string][]*Subscription
}

// Add a new subscription in the struct
func (s *Subscriptions) AddSubscription(newSub *Subscription) {
	key := circle.GetFullNameFromOwnerAndRepo(newSub.Owner, newSub.Repository)

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

// Remove a subscription from the struct
// Return true if the subscription has been found and removed
func (s *Subscriptions) RemoveSubscription(channelID, owner, repository string) bool {
	key := circle.GetFullNameFromOwnerAndRepo(owner, repository)

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
	key := circle.GetFullNameFromOwnerAndRepo(owner, repository)

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

func (s *Subscriptions) GetFilteredChannelsForBuild(build *circle.WebhookInfo) []string {
	subs := s.GetSubscriptionsForRepository(build.Owner, build.Repository)
	if subs == nil {
		return nil
	}

	var channelIDs []string
	for _, sub := range subs {
		switch { // nolint:gocritic // It's expected that more flags get added.
		case build.IsFailed || !sub.Flags.OnlyFailedBuilds:
			channelIDs = append(channelIDs, sub.ChannelID)
		}
	}

	return channelIDs
}
