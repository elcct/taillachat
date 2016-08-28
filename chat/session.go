package chat

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

// Session keeps current chat session
type Session struct {
	ID      string
	Region  string
	IsReady bool
	Session *sockjs.Session
	Room    *Room
}
