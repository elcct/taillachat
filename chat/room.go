package chat

import (
	"encoding/json"
)

// Room keeps list of sessions in the room
type Room struct {
	Sessions []*Session
}

// Broadcast sends message to all sessions in the room
func (r *Room) Broadcast(event string, data string) {
	msg := &Message{
		"event": event,
		"data":  data,
	}
	out, _ := json.Marshal(msg)
	body := string(out)

	for _, cs := range r.Sessions {
		cs.Session.Send(body)
	}
}

// BroadcastOthers sends message to all but sessionId sessions in the room
func (r *Room) BroadcastOthers(sessionID string, event string, data string) {
	msg := &Message{
		"event": event,
		"data":  data,
	}
	out, _ := json.Marshal(msg)
	body := string(out)

	for _, cs := range r.Sessions {
		if cs.ID != sessionID {
			cs.Session.Send(body)
		}
	}
}
