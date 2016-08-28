package api

import (
	"encoding/json"
	"github.com/elcct/taillachat/chat"
	"github.com/elcct/taillachat/helpers"
	"github.com/golang/glog"
	"github.com/vincent-petithory/dataurl"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"html/template"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

var Template *template.Template = nil
var MediaContent string = ""
var ChatSessions = chat.NewMap()

func GetRandomKey(data map[string]bool) string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	n := rand.Intn(len(keys))
	return keys[n]
}

func ChatFindMatch(sessionId string, region string) (chatSession *chat.Session) {
	var sessions map[string]bool

	if region == "All UK" {
		sessions = ChatSessions.GetReadyIds()
	} else {
		sessions = ChatSessions.GetReadyIdsByRegion(region)
		if len(sessions) < 2 {
			sessions = ChatSessions.GetReadyIds()
		}
	}

	delete(sessions, sessionId)

	if len(sessions) > 0 {
		otherSessionId := GetRandomKey(sessions)

		return ChatSessions.Get(otherSessionId)
	}

	return
}

// Runs when user is ready to chat
func ChatReady(sessionId string, region string) {
	cs := ChatSessions.Get(sessionId)
	if cs != nil {
		ChatSessions.Action(func() {
			cs.Region = region
			cs.IsReady = true
		})
	}

	// Let's find match
	partner := ChatFindMatch(sessionId, region)

	if partner != nil {
		// Let's create a chat!
		ChatSessions.Action(func() {
			glog.Info("Chat session started")
			cs.IsReady = false
			partner.IsReady = false

			room := &chat.Room{
				Sessions: []*chat.Session{cs, partner},
			}

			cs.Room = room
			partner.Room = room

			room.Broadcast("join", "")
			room.BroadcastOthers(cs.ID, "message", helpers.Parse(Template, "chat/join", cs))
			room.BroadcastOthers(partner.ID, "message", helpers.Parse(Template, "chat/join", partner))
		})

	}
}

func Chat(session sockjs.Session) {
	glog.Info("Session started")
	sessionID := session.ID()

	chatSession := &chat.Session{
		ID:      sessionID,
		IsReady: false,
		Session: &session,
	}

	ChatSessions.Set(sessionID, chatSession)

	acceptedTypes := &map[string]string{
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}

	ticker := time.NewTicker(time.Second)

	// Not too nice, will be refactored later...
	online := func() {
		ready, chatting := ChatSessions.GetNumberOfReadyAndChatting()

		msg := &chat.Message{
			"event": "online",
			"r":     ready,
			"c":     chatting,
		}
		out, _ := json.Marshal(msg)
		body := string(out)
		session.Send(body)
	}

	go func() {
		for _ = range ticker.C {
			online()
		}
	}()
	online()

	for {
		if msg, err := session.Recv(); err == nil {
			var data chat.Message
			json.Unmarshal([]byte(msg), &data)

			switch data["event"] {
			case "ready":
				ChatReady(sessionID, data["region"].(string))
			case "typing":
				ChatSessions.Action(func() {
					chatSession.Room.BroadcastOthers(sessionID, "typing", strconv.FormatBool(data["typing"].(bool)))
				})
			case "send":
				ChatSessions.Action(func() {
					chatSession.Room.BroadcastOthers(sessionID, "message", data["message"].(string))
				})
			case "exit":
				ChatSessions.Action(func() {
					glog.Info("Chat session ended")
					chatSession.Room.BroadcastOthers(sessionID, "exit", "")

					for i, _ := range chatSession.Room.Sessions {
						s := chatSession.Room.Sessions[i]
						if s != chatSession {
							s.Room = nil
						}
						//s.IsReady = true
					}
					chatSession.Room = nil
				})
			case "picture":
				glog.Info("Picture received")
				ChatSessions.Action(func() {
					chatSession.Room.BroadcastOthers(sessionID, "picturebefore", "true")
					dataURL, err := dataurl.DecodeString(data["data"].(string))
					if err != nil {
						glog.Error("Problem decoding file: ", err)
					}
					filename := helpers.GetRandomString(8)

					mt := dataURL.ContentType()

					if ext, ok := (*acceptedTypes)[mt]; ok {
						err = ioutil.WriteFile(MediaContent+filename+ext, dataURL.Data, 0644)
						if err != nil {
							glog.Error("Error saving file: ", err)
						}
						chatSession.Room.BroadcastOthers(sessionID, "picture", data["data"].(string))
					}
				})
			}

			continue
		}
		break
	}

	ticker.Stop()

	ChatSessions.Close(sessionID)

	glog.Info("Session closed")
}
