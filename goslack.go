package goslack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func (e *Event) String() string {
	return fmt.Sprintf("{id: %v, type:%v, channel:%v, user:%v, ts:%v, text:%v}", e.Id, e.Type, e.Channel, e.User, e.Ts, e.Text)
}

func NewClient(token string) (Client, error) {
	resp, err := http.PostForm("https://slack.com/api/rtm.start", url.Values{"token": {token}})
	if err != nil {
		thisError := fmt.Sprintf("Could't start real time slack api. ERR: %v", err)
		return Client{}, errors.New(thisError)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't read response body. ERR: %v", err)
		return Client{}, errors.New(thisError)
	}

	var sr StartResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't decode json. ERR: %v", err)
		return Client{}, errors.New(thisError)
	}

	/*
		Hack because rtm.start returns a
		websocket URL without the port number
		on it and websocket.Dial barfs if that
		port number isn't there.
	*/
	splitUrl := strings.Split(sr.Url, "/")
	splitUrl[2] = splitUrl[2] + ":443"
	sr.Url = strings.Join(splitUrl, "/")

	var Dialer websocket.Dialer
	header := make(http.Header)
	header.Add("Origin", "http://localhost/")
	ws, resp, err := Dialer.Dial(sr.Url, header)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't dial websocket. ERR: %v", err)
		return Client{}, errors.New(thisError)
	}

	return Client{1, ws, sr.Self, token, make(chan Event), make(chan Event)}, nil
}

func (c *Client) PushMessage(channel, message string) {
	c.MsgOut <- Event{c.MsgId, "message", channel, message, "", ""}
	c.MsgId++
}

func (c *Client) SendMessages() {
	for {
		select {
		case msg := <-c.MsgOut:
			if msgb, _ := json.Marshal(msg); len(msgb) >= 16000 {
				msg = Event{msg.Id, msg.Type, msg.Channel, fmt.Sprintf("ERROR! Response too large. %v Bytes!", len(msgb)), "", ""}
			}

			c.Ws.WriteJSON(msg)
			time.Sleep(time.Second * 1)
		}
	}
}

func (c *Client) ReadMessages() {
	msg := Event{}
	for {
		err := c.Ws.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		if (msg != Event{}) {
			c.MsgIn <- msg
			msg = Event{}
		}
	}
}
