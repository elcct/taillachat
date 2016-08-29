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
	"strconv"
	"time"
)

var Template *template.Template = nil
var MediaContent string = ""

var chatMap = chat.NewMap()

func findMatch(sessionID string, region string) (chatSession *chat.Session) {
	var sessions map[string]bool

	if region == "All UK" {
		sessions = chatMap.GetReadyIds()
	} else {
		sessions = chatMap.GetReadyIdsByRegion(region)
		if len(sessions) < 2 {
			sessions = chatMap.GetReadyIds()
		}
	}

	delete(sessions, sessionID)

	if len(sessions) > 0 {
		otherSessionID := helpers.GetRandomKey(sessions)
		return chatMap.Get(otherSessionID)
	}

	return
}

// chatReady runs when user is ready to chat
func chatReady(sessionID string, region string) {
	cs := chatMap.Get(sessionID)
	if cs != nil {
		chatMap.Action(func() {
			cs.Region = region
			cs.IsReady = true
		})
	}

	// Let's find match
	partner := findMatch(sessionID, region)

	if partner != nil {
		// Let's create a chat!
		chatMap.Action(func() {
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

// Chat handles new sockjs Session
func Chat(session sockjs.Session) {
	glog.Info("Session started")
	sessionID := session.ID()

	chatSession := &chat.Session{
		ID:      sessionID,
		IsReady: false,
		Session: session,
	}

	chatMap.Set(sessionID, chatSession)

	acceptedTypes := &map[string]string{
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}

	ticker := time.NewTicker(time.Second)

	// Not too nice, will be refactored later...
	online := func() {
		ready, chatting := chatMap.GetNumberOfReadyAndChatting()

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
				chatReady(sessionID, data["region"].(string))
			case "typing":
				chatMap.Action(func() {
					chatSession.Room.BroadcastOthers(sessionID, "typing", strconv.FormatBool(data["typing"].(bool)))
				})
			case "send":
				chatMap.Action(func() {
					chatSession.Room.BroadcastOthers(sessionID, "message", data["message"].(string))
				})
			case "exit":
				chatMap.Action(func() {
					glog.Info("Chat session ended")
					chatSession.Room.BroadcastOthers(sessionID, "exit", "")

					for i := range chatSession.Room.Sessions {
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
				chatMap.Action(func() {
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

	chatMap.Close(sessionID)

	glog.Info("Session closed")
}
