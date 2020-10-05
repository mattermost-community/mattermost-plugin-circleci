package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

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
	subscribeChannelHint     = "[--flags]"
	subscribeChannelHelpText = "Subscribe the current channel to CircleCI notifications for a project"

	subscribeUnsubscribeChannelTrigger  = "remove"
	subscribeUnsubscribeChannelHint     = "[--flags]"
	subscribeUnsubscribeChannelHelpText = "Unsubscribe the current channel to CircleCI notifications for a project"

	subscribeListAllChannelsTrigger  = "list-channels"
	subscribeListAllChannelsHint     = ""
	subscribeListAllChannelsHelpText = "List all channels in the current team subscribed to a project"
)

func getSubscribeAutoCompleteData() *model.AutocompleteData {
	subscribe := model.NewAutocompleteData(subscribeTrigger, subscribeHint, subscribeHelpText)

	subscribeList := model.NewAutocompleteData(subscribeListTrigger, subscribeListHint, subscribeListHelpText)

	subscribeChannel := model.NewAutocompleteData(subscribeChannelTrigger, subscribeChannelHint, subscribeChannelHelpText)
	subscribeChannel.AddNamedTextArgument(store.FlagOnlyFailedBuilds, "Only receive notifications for failed builds", "[write anything here]", "", false) // TODO use boolean flag when then are available. See https://github.com/mattermost/mattermost-server/pull/14810
	subscribeChannel.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	unsubscribeChannel := model.NewAutocompleteData(subscribeUnsubscribeChannelTrigger, subscribeUnsubscribeChannelHint, subscribeUnsubscribeChannelHelpText)
	unsubscribeChannel.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	listAllSubscribedChannels := model.NewAutocompleteData(subscribeListAllChannelsTrigger, subscribeListAllChannelsHint, subscribeListAllChannelsHelpText)
	listAllSubscribedChannels.AddNamedTextArgument(namedArgProjectName, namedArgProjectHelpText, namedArgProjectHint, namedArgProjectPattern, false)

	subscribe.AddCommand(subscribeChannel)
	subscribe.AddCommand(unsubscribeChannel)
	subscribe.AddCommand(subscribeList)
	subscribe.AddCommand(listAllSubscribedChannels)

	return subscribe
}

func (p *Plugin) executeSubscribe(context *model.CommandArgs, circleciToken string, project *store.ProjectIdentifier, split []string) (*model.CommandResponse, *model.AppError) {
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
		var rawFlags []string
		if len(split) > 1 {
			rawFlags = split[1:]
		}
		return executeSubscribeChannel(p, context, project, rawFlags)

	case subscribeUnsubscribeChannelTrigger:

		return executeUnsubscribeChannel(p, context, project)

	case subscribeListAllChannelsTrigger:
		return executeSubscribeListAllChannels(p, context, project)

	default:
		return p.sendIncorrectSubcommandResponse(context, subscribeTrigger)
	}
}

func executeSubscribeList(p *Plugin, context *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	allSubs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, ":red_circle: Internal error when getting subscriptions"), nil
	}

	subs := allSubs.GetSubscriptionsByChannel(context.ChannelId)
	if subs == nil {
		return p.sendEphemeralResponse(
			context,
			fmt.Sprintf(
				":information_source: This channel is not subscribed to any repository. Try `/%s %s %s`",
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

func executeSubscribeChannel(p *Plugin, context *model.CommandArgs, project *store.ProjectIdentifier, rawFlags []string) (*model.CommandResponse, *model.AppError) {
	// ? TODO check that project exists

	subs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, ":red_circle: Internal error when getting subscriptions"), nil
	}

	newSub := &store.Subscription{
		ChannelID:          context.ChannelId,
		CreatorID:          context.UserId,
		Flags:              store.SubscriptionFlags{},
		ProjectInformation: *project,
	}

	for _, arg := range rawFlags {
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
	wasUpdated := subs.AddSubscription(newSub)

	if err := p.Store.StoreSubscriptions(subs); err != nil {
		p.API.LogError("Unable to store subscriptions", "error", err)
		return p.sendEphemeralResponse(context, ":red_circle: Internal error when storing new subscription."), nil
	}

	var msg string
	if wasUpdated {
		msg = fmt.Sprintf(
			"This channel was already subscribed to notifications from %s. It has been updated with flags `%s`\n"+
				"The [Mattermost Plugin Notify Orb](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify) should already be configured, but you can check it to be sure. See the full guide [here](%s/blob/master/docs/HOW_TO.md#subscribe-to-webhooks-notifications)\n"+
				"**Webhook URL: `%s`**",
			project.ToMarkdown(),
			newSub.Flags.String(),
			manifest.HomepageURL,
			p.getWebhookURL(),
		)
	} else {
		msg = fmt.Sprintf(
			"This channel has been subscribed to notifications from %s with flags: `%s`\n"+
				"#### How to finish setup:\n"+
				"(See the full guide [here](%s#subscribe-to-webhooks-notifications))\n"+
				"1. Setup the [Mattermost Plugin Notify Orb](https://circleci.com/developer/orbs/orb/nathanaelhoun/mattermost-plugin-notify) for your CircleCI project\n"+
				"2. Add the `MM_WEBHOOK` environment variable to your project using the [CircleCI UI](https://circleci.com/docs/2.0/env-vars/#setting-an-environment-variable-in-a-project) or with \n```\n/%s %s %s %s MM_WEBHOOK %s\n```\n"+
				"**Webhook URL: `%s`**",
			project.ToMarkdown(),
			newSub.Flags,
			manifest.HomepageURL,
			commandTrigger,
			projectTrigger,
			projectEnvVarTrigger,
			projectEnvVarAddTrigger,
			p.getWebhookURL(),
			p.getWebhookURL(),
		)
	}

	return p.sendEphemeralResponse(context, msg), nil
}

func executeUnsubscribeChannel(p *Plugin, args *model.CommandArgs, project *store.ProjectIdentifier) (*model.CommandResponse, *model.AppError) {
	subs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(args, ":red_circle: Internal error when getting subscriptions"), nil
	}

	if removed := subs.RemoveSubscription(args.ChannelId, project); !removed {
		return p.sendEphemeralResponse(args,
			fmt.Sprintf(":red_circle: This channel was not subscribed to %s", project.ToMarkdown()),
		), nil
	}

	if err := p.Store.StoreSubscriptions(subs); err != nil {
		p.API.LogError("Unable to store subscriptions", "error", err)
		return p.sendEphemeralResponse(args, ":red_circle: Internal error when storing new subscription."), nil
	}

	return p.sendEphemeralResponse(args,
		fmt.Sprintf(":white_check_mark: Successfully unsubscribed this channel from %s", project.ToMarkdown()),
	), nil
}

func executeSubscribeListAllChannels(p *Plugin, context *model.CommandArgs, project *store.ProjectIdentifier) (*model.CommandResponse, *model.AppError) {
	allSubs, err := p.Store.GetSubscriptions()
	if err != nil {
		p.API.LogError("Unable to get subscriptions", "err", err)
		return p.sendEphemeralResponse(context, ":red_circle: Internal error when getting subscriptions"), nil
	}

	channelIDs := allSubs.GetSubscribedChannelsForProject(project)
	if channelIDs == nil {
		return p.sendEphemeralResponse(
			context,
			fmt.Sprintf(
				":information_source: No channel is subscribed to the project %s. Try `/%s %s %s`",
				project.ToMarkdown(),
				commandTrigger,
				subscribeTrigger,
				subscribeChannelTrigger,
			),
		), nil
	}

	message := fmt.Sprintf(":information_source: Channels of this team subscribed to %s\n", project.ToMarkdown())
	for _, channelID := range channelIDs {
		channel, appErr := p.API.GetChannel(channelID)
		if appErr != nil {
			p.API.LogError("Unable to get channel", "channelID", channelID)
		}

		message += fmt.Sprintf("- ~%s\n", channel.Name)
	}

	return p.sendEphemeralResponse(context, message), nil
}
