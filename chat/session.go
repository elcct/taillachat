package chat

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

// Session keeps current chat session
type Session struct {
	ID      string
	Session *sockjs.Session
	Room    *Room
	Region  string
	IsReady bool
}
