package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	subscribeTrigger = "subscription"
	subscribeHint    = "<" + subscribeListTrigger + "|" +
		subscribeChannelTrigger + "|" +
		subscribeUnsubscribeChannelTrigger + "|" +
		subscribeListAllChannelsTrigger + ">"
	subscribeHelpText = "Manage your subscriptions"

	subscribeListTrigger  = "list"
	subscribeListHint     = ""
	subscribeListHelpText = "List the CircleCI subscriptions for the current channel"

	subscribeChannelTrigger  = "subscribe"
	subscribeChannelHint     = "<username> <repository> [--flags]"
	subscribeChannelHelpText = "Subscribe the current channel to CircleCI notifications for a repository"

	subscribeUnsubscribeChannelTrigger  = "unsubscribe"
	subscribeUnsubscribeChannelHint     = "<username> <repository> [--flags]"
	subscribeUnsubscribeChannelHelpText = "Unsubscribe the current channel to CircleCI notifications for a repository"

	subscribeListAllChannelsTrigger  = "list-channels"
	subscribeListAllChannelsHint     = "<username> <repository>"
	subscribeListAllChannelsHelpText = "List all channels subscribed to this repository in the current team"
)

func (p *Plugin) executeSubscribe(context *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := commandHelpTrigger
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	case commandHelpTrigger:
		return p.sendHelpResponse(context, subscribeTrigger)

	case subscribeListTrigger:
		return executeSubscribeList(p, context)

	case subscribeChannelTrigger:
		return executeSubscribeChannel(p, context, split[1:])

	case subscribeUnsubscribeChannelTrigger:
		return executeUnsubscribeChannel(p, context, split[1:])

	case subscribeListAllChannelsTrigger:
		return executeSubscribeListAllChannels(p, context, split[1:])

	default:
		return p.sendIncorrectSubcommandResponse(context, subscribeTrigger)
	}
}

func executeSubscribeList(p *Plugin, context *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	allSubs, err := p.getSubscriptionsKV()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, "Internal error when getting subscriptions"), nil
	}

	subs := allSubs.GetSubscriptionsByChannel(context.ChannelId)
	if subs == nil {
		return p.sendEphemeralResponse(
			context,
			fmt.Sprintf(
				"This channel is not subscribed to any repository. Try `/%s %s %s`",
				commandTrigger,
				subscribeTrigger,
				subscribeChannelTrigger,
			),
		), nil
	}

	attachment := model.SlackAttachment{
		Title:    "Repositories this channel is subscribed to :",
		Fallback: "List of repositories this channel is subscribed to",
	}

	for _, sub := range subs {
		p.API.LogDebug("Parsing CircleCI subscription", "sub", sub)
		attachment.Fields = append(attachment.Fields, sub.ToSlackAttachmentField(p))
	}

	p.sendEphemeralPost(context, "", []*model.SlackAttachment{&attachment})
	return &model.CommandResponse{}, nil
}

func executeSubscribeChannel(p *Plugin, context *model.CommandArgs, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(context, "Please provide the project owner and repository names)"), nil
	}

	owner, repo := split[0], split[1]

	// ? TODO check that project exists

	subs, err := p.getSubscriptionsKV()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, "Internal error when getting subscriptions"), nil
	}

	newSub := &Subscription{
		ChannelID:  context.ChannelId,
		CreatorID:  context.UserId,
		Owner:      owner,
		Repository: repo,
		Flags:      SubscriptionFlags{},
	}

	for _, arg := range split[2:] {
		if strings.HasPrefix(arg, "--") {
			flag := arg[2:]
			err := newSub.Flags.AddFlag(flag)
			if err != nil {
				return p.sendEphemeralResponse(context, fmt.Sprintf(
					"Unknown subscription flag `%s`. Try `/%s %s %s`",
					arg,
					commandTrigger,
					subscribeTrigger,
					commandHelpTrigger,
				)), nil
			}
		}
	}

	p.API.LogDebug("Adding a new subscription", "subscription", newSub)
	subs.AddSubscription(newSub)

	if err := p.storeSubscriptionsKV(subs); err != nil {
		p.API.LogError("Unable to store subscriptions", "error", err)
		return p.sendEphemeralResponse(context, "Internal error when storing new subscription."), nil
	}

	// TODO add message "add the orb, here is the docs for doing it"
	return p.sendEphemeralResponse(context, fmt.Sprintf(
		"Successfully subscribed this channel to notifications from **%s**",
		getFullNameFromOwnerAndRepo(owner, repo),
	)), nil
}

func executeUnsubscribeChannel(p *Plugin, context *model.CommandArgs, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(context, "Please provide the project owner and repository names)"), nil
	}

	owner, repo := split[0], split[1]

	subs, err := p.getSubscriptionsKV()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, "Internal error when getting subscriptions"), nil
	}

	if removed := subs.RemoveSubscription(context.ChannelId, owner, repo); !removed {
		return p.sendEphemeralResponse(context, fmt.Sprintf("This channel is not subscribed to **%s**", getFullNameFromOwnerAndRepo(owner, repo))), nil
	}

	if err := p.storeSubscriptionsKV(subs); err != nil {
		p.API.LogError("Unable to store subscriptions", "error", err)
		return p.sendEphemeralResponse(context, "Internal error when storing new subscription."), nil
	}

	return p.sendEphemeralResponse(context, fmt.Sprintf(
		"Successfully unsubscribed this channel to notifications from **%s**",
		getFullNameFromOwnerAndRepo(owner, repo),
	)), nil
}

func executeSubscribeListAllChannels(p *Plugin, context *model.CommandArgs, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(context, "Please provide the project owner and repository names)"), nil
	}

	owner, repo := split[0], split[1]

	allSubs, err := p.getSubscriptionsKV()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, "Internal error when getting subscriptions"), nil
	}

	channelIDs := allSubs.GetSubscribedChannelsForRepository(owner, repo)
	if channelIDs == nil {
		return p.sendEphemeralResponse(
			context,
			fmt.Sprintf(
				"No channel is subscribed to this repository. Try `/%s %s %s`",
				commandTrigger,
				subscribeTrigger,
				subscribeChannelTrigger,
			),
		), nil
	}

	message := "Channels of this team subscribed to **" + getFullNameFromOwnerAndRepo(owner, repo) + "**\n"
	for _, channelID := range channelIDs {
		channel, appErr := p.API.GetChannel(channelID)
		if appErr != nil {
			p.API.LogError("Unable to get channel", "channelID", channelID)
		}

		message += fmt.Sprintf("- ~%s\n", channel.Name)
	}

	return p.sendEphemeralResponse(context, message), nil
}
