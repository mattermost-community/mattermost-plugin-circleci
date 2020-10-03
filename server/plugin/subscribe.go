package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	v1 "github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle/v1"
	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/store"
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

	subscribeChannelTrigger  = "add"
	subscribeChannelHint     = "<owner> <repository> [--flags]"
	subscribeChannelHelpText = "Subscribe the current channel to CircleCI notifications for a repository"

	subscribeUnsubscribeChannelTrigger  = "remove"
	subscribeUnsubscribeChannelHint     = "<owner> <repository> [--flags]"
	subscribeUnsubscribeChannelHelpText = "Unsubscribe the current channel to CircleCI notifications for a repository"

	subscribeListAllChannelsTrigger  = "list-channels"
	subscribeListAllChannelsHint     = "<owner> <repository>"
	subscribeListAllChannelsHelpText = "List all channels subscribed to this repository in the current team"
)

func getSubscribeAutoCompleteData() *model.AutocompleteData {
	subscribe := model.NewAutocompleteData(subscribeTrigger, subscribeHint, subscribeHelpText)

	subscribeList := model.NewAutocompleteData(subscribeListTrigger, subscribeListHint, subscribeListHelpText)
	subscribeChannel := model.NewAutocompleteData(subscribeChannelTrigger, subscribeChannelHint, subscribeChannelHelpText)
	subscribeChannel.AddTextArgument("Owner of the project's repository", "[owner]", "")
	subscribeChannel.AddDynamicListArgument("", routeAutocomplete+subrouteFollowedProjects, true)
	subscribeChannel.AddNamedTextArgument(store.FlagOnlyFailedBuilds, "Only receive notifications for failed builds", "[write anything here]", "", false)
	unsubscribeChannel := model.NewAutocompleteData(subscribeUnsubscribeChannelTrigger, subscribeUnsubscribeChannelHint, subscribeUnsubscribeChannelHelpText)
	unsubscribeChannel.AddTextArgument("Owner of the project's repository", "[owner]", "") // TODO make dynamic autocomplete list
	unsubscribeChannel.AddTextArgument("Repository name", "[repository]", "")              // TODO make dynamic autocomplete list
	listAllSubscribedChannels := model.NewAutocompleteData(subscribeListAllChannelsTrigger, subscribeListAllChannelsHint, subscribeListAllChannelsHelpText)
	listAllSubscribedChannels.AddTextArgument("Owner of the project's repository", "[owner]", "") // TODO make dynamic autocomplete list
	listAllSubscribedChannels.AddTextArgument("Repository name", "[repository]", "")              // TODO make dynamic autocomplete list

	subscribe.AddCommand(subscribeList)
	subscribe.AddCommand(subscribeChannel)
	subscribe.AddCommand(unsubscribeChannel)
	subscribe.AddCommand(listAllSubscribedChannels)

	return subscribe
}

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
	allSubs, err := p.Store.GetSubscriptions()
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

		username := "Unknown user"
		if user, appErr := p.API.GetUser(sub.CreatorID); appErr != nil {
			p.API.LogError("Unable to get username", "userID", sub.CreatorID)
		} else {
			username = user.Username
		}

		attachment.Fields = append(attachment.Fields, sub.ToSlackAttachmentField(username))
	}

	p.sendEphemeralPost(context, "", []*model.SlackAttachment{&attachment})
	return &model.CommandResponse{}, nil
}

func executeSubscribeChannel(p *Plugin, context *model.CommandArgs, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(context, "Please provide the project owner and repository names"), nil
	}

	owner, repo := split[0], split[1]

	// ? TODO check that project exists

	subs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, "Internal error when getting subscriptions"), nil
	}

	newSub := &store.Subscription{
		ChannelID:  context.ChannelId,
		CreatorID:  context.UserId,
		Owner:      owner,
		Repository: repo,
		Flags:      store.SubscriptionFlags{},
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

	if err := p.Store.StoreSubscriptions(subs); err != nil {
		p.API.LogError("Unable to store subscriptions", "error", err)
		return p.sendEphemeralResponse(context, "Internal error when storing new subscription."), nil
	}

	// TODO add message "add the orb, here is the docs for doing it"
	return p.sendEphemeralResponse(context, fmt.Sprintf(
		"Successfully subscribed this channel to notifications from **%s**\nSend webhooks to `%s`",
		v1.GetFullNameFromOwnerAndRepo(owner, repo),
		p.getWebhookURL(),
	)), nil
}

func executeUnsubscribeChannel(p *Plugin, context *model.CommandArgs, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(context, "Please provide the project owner and repository names"), nil
	}

	owner, repo := split[0], split[1]

	subs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, "Internal error when getting subscriptions"), nil
	}

	if removed := subs.RemoveSubscription(context.ChannelId, owner, repo); !removed {
		return p.sendEphemeralResponse(context, fmt.Sprintf("This channel is not subscribed to **%s**",
			v1.GetFullNameFromOwnerAndRepo(owner, repo))), nil
	}

	if err := p.Store.StoreSubscriptions(subs); err != nil {
		p.API.LogError("Unable to store subscriptions", "error", err)
		return p.sendEphemeralResponse(context, "Internal error when storing new subscription."), nil
	}

	return p.sendEphemeralResponse(context, fmt.Sprintf(
		"Successfully unsubscribed this channel to notifications from **%s**",
		v1.GetFullNameFromOwnerAndRepo(owner, repo),
	)), nil
}

func executeSubscribeListAllChannels(p *Plugin, context *model.CommandArgs, split []string) (*model.CommandResponse, *model.AppError) {
	if len(split) < 2 {
		return p.sendEphemeralResponse(context, "Please provide the project owner and repository names"), nil
	}

	owner, repo := split[0], split[1]

	allSubs, err := p.Store.GetSubscriptions()
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

	message := "Channels of this team subscribed to **" + v1.GetFullNameFromOwnerAndRepo(owner, repo) + "**\n"
	for _, channelID := range channelIDs {
		channel, appErr := p.API.GetChannel(channelID)
		if appErr != nil {
			p.API.LogError("Unable to get channel", "channelID", channelID)
		}

		message += fmt.Sprintf("- ~%s\n", channel.Name)
	}

	return p.sendEphemeralResponse(context, message), nil
}
