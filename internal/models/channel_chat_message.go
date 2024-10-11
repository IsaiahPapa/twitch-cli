package models

type ChannelChatMessageEventSubEvent struct {
	BroadcasterUserID           string  `json:"broadcaster_user_id"`
	BroadcasterUserLogin        string  `json:"broadcaster_user_login"`
	BroadcasterUserName         string  `json:"broadcaster_user_name"`
	ChatterUserID               string  `json:"chatter_user_id"`
	ChatterUserLogin            string  `json:"chatter_user_login"`
	ChatterUserName             string  `json:"chatter_user_name"`
	MessageID                   string  `json:"message_id"`
	Message                     Message `json:"message"`
	Color                       string  `json:"color"`
	Badges                      []Badge `json:"badges"`
	MessageType                 string  `json:"message_type"`
	Cheer                       *Cheer  `json:"cheer"`
	Reply                       *Reply  `json:"reply"`
	ChannelPointsCustomRewardID *string `json:"channel_points_custom_reward_id"`
	ChannelPointsAnimationID    *string `json:"channel_points_animation_id"`
}

type ChannelChatMessageEventSubResponse struct {
	Subscription EventsubSubscription      `json:"subscription"`
	Event        AdBreakBeginEventSubEvent `json:"event"`
}

type Message struct {
	Text      string     `json:"text"`
	Fragments []Fragment `json:"fragments"`
}

type Fragment struct {
	Type      string     `json:"type"`
	Text      string     `json:"text"`
	Cheermote *Cheermote `json:"cheermote"`
	Emote     *Emote     `json:"emote"`
	Mention   *Mention   `json:"mention"`
}

type Cheermote struct {
	Prefix string `json:"prefix"`
	Bits   int    `json:"bits"`
	Tier   int    `json:"tier"`
}

type Emote struct {
	Id         string   `json:"id"`
	EmoteSetId string   `json:"emote_set_id"`
	OwnerId    string   `json:"owner_id"`
	Format     []string `json:"format"`
}

type Mention struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type Badge struct {
	SetID string `json:"set_id"`
	ID    string `json:"id"`
	Info  string `json:"info"`
}

type Cheer struct {
	Bits int `json:"bits"`
}

type Reply struct {
	ParentMessageID   string `json:"parent_message_id"`
	ParentMessageBody string `json:"parent_message_body"`
	ParentUserID      string `json:"parent_user_id"`
	ParentUserName    string `json:"parent_user_name"`
	ParentUserLogin   string `json:"parent_user_login"`
	ThreadMessageId   string `json:"thread_message_id"`
	ThreadUserId      string `json:"thread_user_id"`
	ThreadUserName    string `json:"thread_user_name"`
	ThreadUserLogin   string `json:"thread_user_login"`
}
