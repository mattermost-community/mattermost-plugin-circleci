package main

import (
	"sort"
)

const (
	subscriptionsKVKey = "subscriptions"
)

type Subscriptions struct {
	Repositories map[string][]*Subscription
}

// Add a new subscription in the struct
func (s *Subscriptions) AddSubscription(newSub *Subscription) {
	key := getFullNameFromOwnerAndRepo(newSub.Owner, newSub.Repository)

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
	key := getFullNameFromOwnerAndRepo(owner, repository)

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

// Return a list of subscribed channel IDs
func (s *Subscriptions) GetSubscribedChannelsForRepository(owner, repository string) []string {
	key := getFullNameFromOwnerAndRepo(owner, repository)

	subs := s.Repositories[key]
	if subs == nil {
		return nil
	}

	var channelIDs []string
	for _, sub := range subs {
		channelIDs = append(channelIDs, sub.ChannelID)
	}

	return channelIDs
}
