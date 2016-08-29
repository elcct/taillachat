package chat

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockSession struct {
	LastMessage string
}

func (s *mockSession) ID() string {
	return ""
}

func (s *mockSession) Request() *http.Request {
	return nil
}

func (s *mockSession) Recv() (string, error) {
	return "", nil
}

func (s *mockSession) Send(m string) error {
	s.LastMessage = m
	return nil
}

func (s *mockSession) Close(status uint32, reason string) error {
	return nil
}

type message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func TestBroadcast(t *testing.T) {
	s1 := &Session{
		ID:      "s1",
		Session: &mockSession{},
	}
	s2 := &Session{
		ID:      "s2",
		Session: &mockSession{},
	}

	room := &Room{
		Sessions: []*Session{s1, s2},
	}

	room.Broadcast("test", "test")

	var (
		result1 message
		result2 message
	)

	err := json.Unmarshal([]byte(s1.Session.(*mockSession).LastMessage), &result1)
	assert.Nil(t, err)

	err = json.Unmarshal([]byte(s2.Session.(*mockSession).LastMessage), &result2)
	assert.Nil(t, err)

	expected := message{
		Event: "test",
		Data:  "test",
	}

	assert.Equal(t, expected, result1)
	assert.Equal(t, expected, result2)
}

func TestBroadcastOthers(t *testing.T) {
	s1 := &Session{
		ID:      "s1",
		Session: &mockSession{},
	}
	s2 := &Session{
		ID:      "s2",
		Session: &mockSession{},
	}

	room := &Room{
		Sessions: []*Session{s1, s2},
	}

	var (
		result1 message
		result2 message
	)

	room.BroadcastOthers("s2", "test", "test2")
	room.BroadcastOthers("s1", "test", "test1")

	err := json.Unmarshal([]byte(s1.Session.(*mockSession).LastMessage), &result1)
	assert.Nil(t, err)

	err = json.Unmarshal([]byte(s2.Session.(*mockSession).LastMessage), &result2)
	assert.Nil(t, err)

	expected1 := message{
		Event: "test",
		Data:  "test2",
	}

	expected2 := message{
		Event: "test",
		Data:  "test1",
	}

	assert.Equal(t, expected1, result1)
	assert.Equal(t, expected2, result2)
}
