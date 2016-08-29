package chat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSetCloseMap(t *testing.T) {
	m := NewMap()

	assert.NotNil(t, m)
	s1 := &Session{
		ID: "s1",
	}
	s2 := &Session{
		ID: "s2",
	}
	m.Set("s1", s1)
	m.Set("s2", s2)
	assert.Equal(t, s1, m.Get("s1"))
	assert.Equal(t, s2, m.Get("s2"))
	m.Close("s1")
	m.Close("s2")
	assert.Nil(t, m.Get("s1"))
	assert.Nil(t, m.Get("s2"))
}

func TestGetReadyIds(t *testing.T) {
	m := NewMap()

	s1 := &Session{
		ID:      "s1",
		IsReady: false,
	}
	s2 := &Session{
		ID:      "s2",
		IsReady: true,
	}
	m.Set("s1", s1)
	m.Set("s2", s2)

	result := m.GetReadyIds()

	expected := map[string]bool{
		"s2": true,
	}

	assert.Equal(t, expected, result)
}

func TestGetNumberOfReadyAndChatting(t *testing.T) {
	m := NewMap()

	s1 := &Session{
		ID:      "s1",
		IsReady: true,
	}
	s2 := &Session{
		ID:      "s2",
		IsReady: false,
		Room:    &Room{},
	}
	s3 := &Session{
		ID:      "s3",
		IsReady: false,
		Room:    &Room{},
	}
	m.Set("s1", s1)
	m.Set("s2", s2)
	m.Set("s3", s3)

	ready, chatting := m.GetNumberOfReadyAndChatting()
	assert.Equal(t, 1, ready)
	assert.Equal(t, 2, chatting)
}
