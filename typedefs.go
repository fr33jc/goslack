package goslack

import "golang.org/x/net/websocket"

type Client struct {
	msgId    int
	ws       *websocket.Conn
	messages []Event
	self     Self
}

type Self struct {
	Created        int                    `json:"created,omitempty"`
	ManualPresence string                 `json:"manual_presence,omitempty"`
	Name           string                 `json:"name,omitempty"`
	Id             string                 `json:"id,omitempty"`
	Prefs          map[string]interface{} `json:"prefs,omitempty"`
}

type Channel struct {
	Id          string                 `json:"id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Created     int                    `json:"created,omitempty"`
	Creator     string                 `json:"creator,omitempty"`
	IsArchived  bool                   `json:"is_archived"`
	IsGeneral   bool                   `json:"is_general"`
	IsMember    bool                   `json:"is_member"`
	LastRead    string                 `json:"last_read,omitempty"`
	Latest      Event                  `json:"latest,omitempty"`
	Members     []string               `json:"members,omitempty"`
	Purpose     interface{}            `json:"purpose,omitempty"`
	Topic       map[string]interface{} `json:"topic,omitempty"`
	UnreadCount int                    `json:"unread_count,omitempty"`
}

type StartResponse struct {
	Ok           bool      `json:"ok,omitempty"`
	Url          string    `json:"url,omitempty"`
	Self         Self      `json:"self,omitempty"`
	CacheVersion string    `json:"cache_version,omitempty"`
	Channels     []Channel `json:"channels,omitempty"`
}

type Event struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
	Ts      string `json:"ts,omitempty"`
}
