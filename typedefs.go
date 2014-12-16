package goslack

import "golang.org/x/net/websocket"

type Client struct {
	msgId    int
	ws       *websocket.Conn
	messages []Message
	self     Self
}

type Self struct {
	Created        *int                    `json:"created,omitempty"`
	ManualPresence *string                 `json:"manual_presence,omitempty"`
	Name           *string                 `json:"name,omitempty"`
	Id             *string                 `json:"id,omitempty"`
	Prefs          *map[string]interface{} `json:"prefs,omitempty"`
}

type Message struct {
	Ok           *bool   `json:"ok,omitempty"`
	Url          *string `json:"url,omitempty"`
	Self         *Self   `json:"self,omitempty"`
	CacheVersion *string `json:"cache_version,omitempty"`
}

type Response struct {
	Ok      bool          `json:"ok"`
	Members []interface{} `json:"members"`
	Url     string        `json:"url"`
}

type MessageSend struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type MessageRecv struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Ts      string `json:"ts"`
	Text    string `json:"text"`
}
