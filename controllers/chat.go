package controllers

import (
	"encoding/json"
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

type ChatMessage map[string]interface{}

type ChatSession struct {
	Id      string
	Region  string
	IsReady bool
	Session *sockjs.Session

	Room *ChatRoom
}

type ChatMap struct {
	Sessions map[string]*ChatSession
	actions  chan func()
}

type ChatRoom struct {
	Sessions []*ChatSession
}

func (cr *ChatRoom) Broadcast(event string, data string) {
	msg := &ChatMessage{
		"event": event,
		"data":  data,
	}
	out, _ := json.Marshal(msg)
	body := string(out)

	for _, cs := range cr.Sessions {
		(*cs.Session).Send(body)
	}
}

func (cr *ChatRoom) BroadcastOthers(sessionId string, event string, data string) {
	msg := &ChatMessage{
		"event": event,
		"data":  data,
	}
	out, _ := json.Marshal(msg)
	body := string(out)

	for _, cs := range cr.Sessions {
		if cs.Id != sessionId {
			(*cs.Session).Send(body)
		}
	}
}

func NewChatMap() *ChatMap {
	chatMap := &ChatMap{
		Sessions: make(map[string]*ChatSession),
		actions:  make(chan func()),
	}

	go func() {
		for action := range chatMap.actions {
			action()
		}
	}()

	return chatMap
}

func (cm *ChatMap) Set(id string, session *ChatSession) {
	cm.actions <- func() {
		cm.Sessions[id] = session
	}
}

func (cm *ChatMap) Get(id string) (session *ChatSession) {
	wait := make(chan bool)
	cm.actions <- func() {
		session = cm.Sessions[id]
		wait <- true
	}
	<-wait
	return
}

func (cm *ChatMap) GetReadyIdsByRegion(region string) (sessions map[string]bool) {
	sessions = make(map[string]bool)

	wait := make(chan bool)
	cm.actions <- func() {
		for key := range cm.Sessions {
			session := cm.Sessions[key]
			if session.Region == region && session.IsReady {
				sessions[session.Id] = true
			}
		}
		wait <- true
	}
	<-wait
	return
}

func (cm *ChatMap) GetReadyIds() (sessions map[string]bool) {
	sessions = make(map[string]bool)

	wait := make(chan bool)
	cm.actions <- func() {
		for key := range cm.Sessions {
			session := cm.Sessions[key]
			if session.IsReady {
				sessions[session.Id] = true
			}
		}
		wait <- true
	}
	<-wait
	return
}

func (cm *ChatMap) GetNumberOfReadyAndChatting() (ready int, chatting int) {
	wait := make(chan bool)

	cm.actions <- func() {
		for key := range cm.Sessions {
			session := cm.Sessions[key]
			if session.IsReady {
				ready++
			}
			if session.Room != nil {
				chatting++
			}
		}
		wait <- true
	}

	<-wait
	return
}

func (cm *ChatMap) Close(id string) {
	cm.actions <- func() {
		delete(cm.Sessions, id)
	}
}

func (cm *ChatMap) Action(fn func()) {
	wait := make(chan bool)
	cm.actions <- func() {
		fn()
		wait <- true
	}
	<-wait
}

var ChatSessions = NewChatMap()

func GetRandomKey(data map[string]bool) string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	n := rand.Intn(len(keys))
	return keys[n]
}

func ChatFindMatch(sessionId string, region string) (chatSession *ChatSession) {
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

			room := &ChatRoom{
				Sessions: []*ChatSession{cs, partner},
			}

			cs.Room = room
			partner.Room = room

			room.Broadcast("join", "")
			room.BroadcastOthers(cs.Id, "message", helpers.Parse(Template, "chat/join", cs))
			room.BroadcastOthers(partner.Id, "message", helpers.Parse(Template, "chat/join", partner))
		})

	}
}

func Chat(session sockjs.Session) {
	glog.Info("Session started")
	sessionId := session.ID()

	chatSession := &ChatSession{
		Id:      sessionId,
		IsReady: false,
		Session: &session,
	}

	ChatSessions.Set(sessionId, chatSession)

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

		msg := &ChatMessage{
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
			var data ChatMessage
			json.Unmarshal([]byte(msg), &data)

			switch data["event"] {
			case "ready":
				ChatReady(sessionId, data["region"].(string))
			case "typing":
				ChatSessions.Action(func() {
					chatSession.Room.BroadcastOthers(sessionId, "typing", strconv.FormatBool(data["typing"].(bool)))
				})
			case "send":
				ChatSessions.Action(func() {
					chatSession.Room.BroadcastOthers(sessionId, "message", data["message"].(string))
				})
			case "exit":
				ChatSessions.Action(func() {
					glog.Info("Chat session ended")
					chatSession.Room.BroadcastOthers(sessionId, "exit", "")

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
					chatSession.Room.BroadcastOthers(sessionId, "picturebefore", "true")
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
						chatSession.Room.BroadcastOthers(sessionId, "picture", data["data"].(string))
					}
				})
			}

			continue
		}
		break
	}

	ticker.Stop()

	ChatSessions.Close(sessionId)

	glog.Info("Session closed")
}
