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

type SlackResponse struct {
	Ok      bool          `json:"ok"`
	Members []interface{} `json:"members"`
	Url     string        `json:"url"`
}

type SlackMessage struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func ConnectSlack(token string, ws *websocket.Conn) error {
	resp, err := http.PostForm("https://slack.com/api/rtm.start", url.Values{"token": {token}})
	if err != nil {
		thisError := fmt.Sprintf("Could't start real time slack api. ERR: %v", err)
		return errors.New(thisError)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't read response body. ERR: %v", err)
		return errors.New(thisError)
	}

	var sr SlackResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		thisError := fmt.Sprintf("Couldn't decode json. ERR: %v", err)
		return errors.New(thisError)
	}

	splitUrl := strings.Split(sr.Url, "/")
	splitUrl[2] = splitUrl[2] + ":443"
	sr.Url = strings.Join(splitUrl, "/")

	ws, err = websocket.Dial(sr.Url, "", "http://localhost/")
	if err != nil {
		thisError := fmt.Sprintf("Couldn't dial websocket. ERR: %v", err)
		return errors.New(thisError)
	}
	return nil
}
