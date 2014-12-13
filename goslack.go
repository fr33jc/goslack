package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

type SlackConfig struct {
	token string
}

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

func NewSlackConfig() *SlackConfig {
	return &Gopher{token: "xoxb-3211150999-AlHt3inJ1QAvTrzIYuoD2W1B"}
}

func Connect(config SlackConfig) ws websocket, err error {
	resp, err := http.PostForm("https://slack.com/api/rtm.start", url.Values{"token": {g.token}})
	if err != nil {
		fmt.Printf("Couldn't start realt time slack api. ERR: %v", err)
		os.Exit(-1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Couldn't read response body. ERR: %v", err)
		os.Exit(-1)
	}

	var sr SlackResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		fmt.Printf("Couldn't decode json. ERR: %v", err)
		os.Exit(-1)
	}

	splitUrl := strings.Split(sr.Url, "/")
	splitUrl[2] = splitUrl[2] + ":443"
	sr.Url = strings.Join(splitUrl, "/")

	ws, err := websocket.Dial(sr.Url, "", "http://localhost/")
	if err != nil {
		fmt.Printf("Couldn't dial websocket. ERR: %v", err)
		os.Exit(-1)
	}
	defer ws.Close()
}

func main() {

	g := NewGopher()
	fmt.Println("Connected to slack")
	//	msg_scanner := bufio.NewScanner(ws)
	//var msg bytes.Buffer
	//for msg_scanner.Scan() {
	//		msg.WriteString(msg_scanner.Text())
	//	}
	//	fmt.Println(msg.String())
	var message SlackMessage
	message.Id += 1
	message.Type = "message"
	message.Channel = "C0369QNBH"
	message.Text = "hello slack"
	b_message, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Could not marshal the message. ERR: %v", err)
		os.Exit(-1)
	}
	if _, err := ws.Write(b_message); err != nil {
		fmt.Println("Couldn't write message. ERR: %v", err)
		os.Exit(-1)
	}

}
