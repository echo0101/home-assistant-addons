package main

import (
	"flag"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/gorilla/websocket"
	"net/http"
)

// Messages

type HaEnvelope struct {
	Type string `json:"type"`
}

const (
	// RX
	TYPE_AUTH_REQUIRED = "auth_required"
	TYPE_AUTH_OK = "auth_ok"
	TYPE_AUTH_INVALID = "auth_invalid"
	TYPE_RESULT = "result"
	TYPE_EVENT = "event"

	//TX
	TYPE_AUTH = "auth"
	TYPE_SUB_EVENTS = "subscribe_events"

	// EVENT TYPES
	EVENT_NOTIFY_KEYBASE = "NOTIFY_KEYBASE"
)

type AuthMessage struct {
	Type string `json:"type"`
	AccessToken string `json:"access_token"`
}

type SubscribeMessage struct {
	Id int `json:"id"`
	Type string `json:"type"`
	EventType string `json:"event_type"`
}

type ErrorMessage struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

type ResultMessage struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Success bool `json:"success"`
	Error ErrorMessage `json:"error"`
}

type EventData struct {
	Message string `json:"message"`
}

type EventInfo struct {
	EventType string `json:"event_type"`
	Data EventData `json:"data"`
}

type EventMessage struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Event EventInfo `json:"event"`
}

// RX auth_required

// TX auth
// required: access_token (str) or api_password (str)

// RX auth_ok
// RX auth_invalid


// TX subscribe_events
// required: id (int)
// optional: event_type (str)

// RX result
// id, success, result

// RX event
// event.data
// - entity_id

type KbConfigData struct {
	Username string `json:"username"`
	PaperKey string `json:"paperKey"`
	TeamName string `json:"teamName"`
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(4)
}

func main() {
	var kbLoc string
	var kbc *kbchat.API
	var err error

	configFile, err := os.Open("/data/options.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	var configData KbConfigData
	err = json.Unmarshal(byteValue, &configData)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	oneShot := kbchat.OneshotOptions{
		Username: configData.Username,
		PaperKey: configData.PaperKey,
	}
	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc, Oneshot: &oneShot}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	/*if _, err = kbc.SendMessageByTeamName(configData.TeamName, "hello!", nil); err != nil {
		fail("Error sending message; %s", err.Error())
	}*/

	headers := http.Header {}
	conn, _, err := websocket.DefaultDialer.Dial("ws://hassio/homeassistant/websocket", headers)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	var authSuccess = false
	for authSuccess == false {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			os.Exit(6)
		}
		if msgType != websocket.TextMessage {
			fmt.Println("Unexpected message type")
			os.Exit(7)
		}
		var envelope HaEnvelope
		err = json.Unmarshal(data, &envelope)
		if err != nil {
			fmt.Println(err)
			os.Exit(8)
		}
		switch envelope.Type {
		case TYPE_AUTH_REQUIRED:
			err := conn.WriteJSON(AuthMessage{TYPE_AUTH,os.Getenv("HASSIO_TOKEN")})
			if err != nil {
				fmt.Println(err)
				os.Exit(9)
			}
		case TYPE_AUTH_OK:
			authSuccess = true
		case TYPE_AUTH_INVALID:
			fmt.Println("INVALID AUTH")
			os.Exit(10)
		default:
			fmt.Println("UNKNOWN MESSAGE TYPE")
			os.Exit(11)
		}
	}
	var msgId=1
	err = conn.WriteJSON(SubscribeMessage{msgId, TYPE_SUB_EVENTS, EVENT_NOTIFY_KEYBASE})
	msgId++
	if err != nil {
		fmt.Println(err)
		os.Exit(12)
	}
	for {
		msgType, data, rErr := conn.ReadMessage()
		if rErr != nil {
			fmt.Println(rErr)
			os.Exit(13)
		}
		if msgType != websocket.TextMessage {
			fmt.Println("Unexpected message type")
			os.Exit(14)
		}
		var envelope HaEnvelope
		err := json.Unmarshal(data, &envelope)
		if err != nil {
			fmt.Println(err)
			os.Exit(15)
		}
		switch envelope.Type {
		case TYPE_RESULT:
			// TODO verify ID matches and success==true
			var resultMessage ResultMessage
			err = json.Unmarshal(data, &resultMessage)
			if err != nil {
				fmt.Println(err)
				os.Exit(16)
			}
			if !resultMessage.Success{
				fmt.Printf("Subscription FAILED: %s\n", data)
				os.Exit(17)
			}
			fmt.Printf("Confirmed Subscription: %+v\n", resultMessage)
		case TYPE_EVENT:
			fmt.Printf("Got Event: %s\n", data)
			var eventMessage EventMessage
			err = json.Unmarshal(data, &eventMessage)
			if err != nil {
				fmt.Println(err)
				os.Exit(18)
			}

			if _, err = kbc.SendMessageByTeamName(configData.TeamName, eventMessage.Event.Data.Message, nil); err != nil {
				fail("Error sending message; %s", err.Error())
			}
		default:
			fmt.Println("UNEXPECTED MESSAGE")
		}
	}
}
