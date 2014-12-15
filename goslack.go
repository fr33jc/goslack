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

	"golang.org/x/net/websocket"
)

type SlackResponse struct {
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

func (m *MessageSend) String() string {
	return fmt.Sprintf("{id:%v, type:%v. channel:%v, text:%v}", m.Id, m.Type, m.Channel, m.Text)
}

func (m *MessageRecv) String() string {
	return fmt.Sprintf("{type:%v. channel:%v, user:%v, ts:%v, text:%v}", m.Type, m.Channel, m.User, m.Ts, m.Text)
}

func Connect(token string) (*websocket.Conn, error) {
	resp, err := http.PostForm("https://slack.com/api/rtm.start", url.Values{"token": {token}})
	if err != nil {
		thisError := fmt.Sprintf("Could't start real time slack api. ERR: %v", err)
		return nil, errors.New(thisError)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't read response body. ERR: %v", err)
		return nil, errors.New(thisError)
	}

	var sr SlackResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't decode json. ERR: %v", err)
		return nil, errors.New(thisError)
	}

	splitUrl := strings.Split(sr.Url, "/")
	splitUrl[2] = splitUrl[2] + ":443"
	sr.Url = strings.Join(splitUrl, "/")

	ws, err := websocket.Dial(sr.Url, "", "http://localhost/")
	if err != nil {
		thisError := fmt.Sprintf("Couldn't dial websocket. ERR: %v", err)
		return nil, errors.New(thisError)
	}

	return ws, nil
}

func SendMessage(ws *websocket.Conn, msg MessageSend) error {
	err := websocket.JSON.Send(ws, msg)
	if err != nil {
		thisError := fmt.Sprintln("Could not send the message. ERR: %v", err)
		return errors.New(thisError)
	}

	return nil
}

func ReadMessages(ws *websocket.Conn, ch chan MessageRecv) error {
	var msg MessageRecv
	for {
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			thisError := fmt.Sprintln("Could not receive the message. ERR: %v", err)
			return errors.New(thisError)
		}
		time.Sleep(1)
		ch <- msg
		msg = MessageRecv{}
	}
}
