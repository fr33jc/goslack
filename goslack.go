package goslack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/websocket"
)

func (e *Event) String() string {
	return fmt.Sprintf("{id: %v, type:%v, channel:%v, user:%v, ts:%v, text:%v}", e.Id, e.Type, e.Channel, e.User, e.Ts, e.Text)
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

	var sr StartResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't decode json. ERR: %v", err)
		return nil, errors.New(thisError)
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

	ws, err := websocket.Dial(sr.Url, "", "http://localhost/")
	if err != nil {
		thisError := fmt.Sprintf("Couldn't dial websocket. ERR: %v", err)
		return nil, errors.New(thisError)
	}

	return ws, nil
}

func SendMessage(ws *websocket.Conn, msg Event) error {
	err := websocket.JSON.Send(ws, msg)
	if err != nil {
		thisError := fmt.Sprintf("Could not send the message. ERR: %v", err)
		return errors.New(thisError)
	}

	return nil
}

func ReadMessages(ws *websocket.Conn) (msg Event, err error) {

	if err := websocket.JSON.Receive(ws, &msg); err != nil {
		return Event{}, err
	}

	return msg, nil
}
