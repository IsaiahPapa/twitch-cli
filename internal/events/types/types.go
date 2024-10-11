// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package types

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/types/ad_break"
	"github.com/twitchdev/twitch-cli/internal/events/types/authorization_grant"
	"github.com/twitchdev/twitch-cli/internal/events/types/authorization_revoke"
	"github.com/twitchdev/twitch-cli/internal/events/types/ban"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_chat"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_redemption"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_reward"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_update_v1"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_update_v2"
	"github.com/twitchdev/twitch-cli/internal/events/types/charity"
	"github.com/twitchdev/twitch-cli/internal/events/types/cheer"
	"github.com/twitchdev/twitch-cli/internal/events/types/drop"
	"github.com/twitchdev/twitch-cli/internal/events/types/extension_transaction"
	"github.com/twitchdev/twitch-cli/internal/events/types/follow"
	"github.com/twitchdev/twitch-cli/internal/events/types/gift"
	"github.com/twitchdev/twitch-cli/internal/events/types/goal"
	"github.com/twitchdev/twitch-cli/internal/events/types/hype_train"
	"github.com/twitchdev/twitch-cli/internal/events/types/moderator_change"
	"github.com/twitchdev/twitch-cli/internal/events/types/poll"
	"github.com/twitchdev/twitch-cli/internal/events/types/prediction"
	"github.com/twitchdev/twitch-cli/internal/events/types/raid"
	"github.com/twitchdev/twitch-cli/internal/events/types/shield_mode"
	"github.com/twitchdev/twitch-cli/internal/events/types/shoutout"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamdown"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamup"
	"github.com/twitchdev/twitch-cli/internal/events/types/subscribe"
	"github.com/twitchdev/twitch-cli/internal/events/types/subscription_message"
	"github.com/twitchdev/twitch-cli/internal/events/types/unban"
	"github.com/twitchdev/twitch-cli/internal/events/types/unban_requests"
	user_update "github.com/twitchdev/twitch-cli/internal/events/types/user"
	"github.com/twitchdev/twitch-cli/internal/models"
)

func AllEvents() []events.MockEvent {
	return []events.MockEvent{
		ad_break.Event{},
		authorization_grant.Event{},
		authorization_revoke.Event{},
		ban.Event{},
		channel_chat.Event{},
		channel_points_redemption.Event{},
		channel_points_reward.Event{},
		charity.Event{},
		cheer.Event{},
		drop.Event{},
		extension_transaction.Event{},
		follow.Event{},
		gift.Event{},
		goal.Event{},
		hype_train.Event{},
		moderator_change.Event{},
		poll.Event{},
		prediction.Event{},
		raid.Event{},
		shield_mode.Event{},
		shoutout.Event{},
		channel_update_v1.Event{},
		channel_update_v2.Event{},
		streamup.Event{},
		streamdown.Event{},
		subscribe.Event{},
		subscription_message.Event{},
		unban.Event{},
		unban_requests.Event{},
		user_update.Event{},
	}
}

func AllWebhookTopics() []string {
	allEvents := []string{}
	allEventsMap := make(map[string]int)

	for _, e := range AllEvents() {
		for _, topic := range e.GetAllTopicsByTransport(models.TransportWebhook) {
			_, duplicate := allEventsMap[topic]
			if !duplicate {
				allEvents = append(allEvents, topic)
				allEventsMap[topic] = 1
			}
		}
	}

	// Sort the topics alphabetically
	sort.Strings(allEvents)

	return allEvents
}

func WebSocketCommandTopics() []string {
	allEvents := []string{}

	for _, e := range AllEvents() {
		for _, topic := range e.GetAllTopicsByTransport(models.TransportWebSocket) {
			if strings.HasPrefix(topic, "websocket") {
				allEvents = append(allEvents, topic)
			}
		}
	}

	// Sort the topics alphabetically
	sort.Strings(allEvents)

	return allEvents
}

func GetByTriggerAndTransportAndVersion(trigger string, transport string, version string) (events.MockEvent, error) {
	validEventBadVersions := []string{}
	var latestEventSeen events.MockEvent

	for _, e := range AllEvents() {
		if transport == models.TransportWebhook || transport == models.TransportWebSocket {
			newTrigger := e.GetEventSubAlias(trigger)
			if newTrigger != "" {
				trigger = newTrigger
			}
		}
		if e.ValidTrigger(trigger) == true && e.ValidTransport(transport) == true {
			if e.SubscriptionVersion() == version {
				return e, nil
			} else {
				validEventBadVersions = append(validEventBadVersions, e.SubscriptionVersion())
				latestEventSeen = e
			}
		}
	}

	// When no version is given, and there's only one version available, use the default version.
	if version == "" && len(validEventBadVersions) == 1 {
		return latestEventSeen, nil
	}

	// Error for events with non-existent version used
	if len(validEventBadVersions) != 0 {
		errStr := fmt.Sprintf("Invalid version given. Valid version(s): %v", strings.Join(validEventBadVersions, ", "))
		if version == "" {
			errStr += "\nUse --version to specify"
		}
		return nil, errors.New(errStr)
	}

	// Error for websocket transport
	if strings.EqualFold(transport, "websocket") {
		return nil, errors.New("Invalid event, or this event is not available via WebSockets.")
	}

	// Default error
	return nil, errors.New("Invalid event")
}

// These events were removed from production
// This does not include any "beta" events, just old production versions
func RemovedEvents() map[string]string {
	return map[string]string{
		"channel.follow": "1",
	}
}
