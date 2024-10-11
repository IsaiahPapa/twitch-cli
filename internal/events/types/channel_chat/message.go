// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_chat

import (
	"encoding/json"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebhook:   true,
	models.TransportWebSocket: true,
}

// var triggerSupported = []string{"clear", "clear-user-messages", "message", "message-delete", "notification", "user-hold-message", "user-message-update"}
var triggerSupported = []string{"message"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebhook: {
		"message": "channel.chat.message",
		// "message-delete":      "channel.chat.message_delete",
		// "clear":               "channel.chat.clear",
		// "clear-user-messages": "channel.chat.clear_user_messages",
		// "notification":        "channel.chat.notification",
		// "user-hold-message":   "channel.chat.user_message_hold",
		// "user-message-update": "channel.chat.user_message_update",
	},
	models.TransportWebSocket: {
		"message": "channel.chat.message",
		// "message-delete":      "channel.chat.message_delete",
		// "clear":               "channel.chat.clear",
		// "clear-user-messages": "channel.chat.clear_user_messages",
		// "notification":        "channel.chat.notification",
		// "user-hold-message":   "channel.chat.user_message_hold",
		// "user-message-update": "channel.chat.user_message_update",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportWebhook, models.TransportWebSocket:
		body := models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.SubscriptionID,
				Status:  params.SubscriptionStatus,
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: e.SubscriptionVersion(),
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
					UserID:            params.FromUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      0,
				CreatedAt: params.Timestamp,
			},
			Event: models.ChannelChatMessageEventSubEvent{
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				ChatterUserID:        params.FromUserID,
				ChatterUserLogin:     params.FromUserName,
				ChatterUserName:      params.FromUserName,
				MessageID:            util.RandomGUID(),
				Message: models.Message{
					Text:      "Test chat message from CLI!",
					Fragments: []models.Fragment{},
				},
				Color:                       util.RandomColorInHex(),
				Badges:                      []models.Badge{},
				MessageType:                 "text",
				Cheer:                       nil,
				Reply:                       nil,
				ChannelPointsCustomRewardID: nil,
				ChannelPointsAnimationID:    nil,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

		// Delete event info if Subscription.Status is not set to "enabled"
		if !strings.EqualFold(params.SubscriptionStatus, "enabled") {
			var i interface{}
			if err := json.Unmarshal([]byte(event), &i); err != nil {
				return events.MockEventResponse{}, err
			}
			if m, ok := i.(map[string]interface{}); ok {
				delete(m, "event") // Matches JSON key defined in body variable above
			}

			event, err = json.Marshal(i)
			if err != nil {
				return events.MockEventResponse{}, err
			}
		}
	default:
		return events.MockEventResponse{}, nil
	}

	return events.MockEventResponse{
		ID:       params.EventMessageID,
		JSON:     event,
		FromUser: params.ToUserID,
		ToUser:   params.ToUserID,
	}, nil
}

func (e Event) ValidTransport(t string) bool {
	return transportsSupported[t]
}

func (e Event) ValidTrigger(t string) bool {
	for _, ts := range triggerSupported {
		if ts == t {
			return true
		}
	}
	return false
}
func (e Event) GetTopic(transport string, trigger string) string {
	return triggerMapping[transport][trigger]
}
func (e Event) GetAllTopicsByTransport(transport string) []string {
	allTopics := []string{}
	for _, topic := range triggerMapping[transport] {
		allTopics = append(allTopics, topic)
	}
	return allTopics
}
func (e Event) GetEventSubAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportWebhook] {
		if topic == t {
			return trigger
		}
	}
	return ""
}

func (e Event) SubscriptionVersion() string {
	return "1"
}
