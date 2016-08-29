package chat

// Map keeps all sessions and actions
type Map struct {
	Sessions map[string]*Session
	actions  chan func()
}

// NewMap creates new map of sessions and start processing actions on them
func NewMap() *Map {
	m := &Map{
		Sessions: make(map[string]*Session),
		actions:  make(chan func()),
	}

	go func() {
		for action := range m.actions {
			action()
		}
	}()

	return m
}

// Set assigns session to the sessions map
func (m *Map) Set(id string, session *Session) {
	m.Action(func() {
		m.Sessions[id] = session
	})
}

// Get gets session from the sessions map
func (m *Map) Get(id string) (session *Session) {
	m.Action(func() {
		session = m.Sessions[id]
	})
	return
}

// GetReadyIdsByRegion gets sessions that are ready for chat in given region
func (m *Map) GetReadyIdsByRegion(region string) (sessions map[string]bool) {
	sessions = make(map[string]bool)

	m.Action(func() {
		for key := range m.Sessions {
			session := m.Sessions[key]
			if session.Region == region && session.IsReady {
				sessions[session.ID] = true
			}
		}
	})
	return
}

// GetReadyIds gets sessions that are ready to chat
func (m *Map) GetReadyIds() (sessions map[string]bool) {
	sessions = make(map[string]bool)

	m.Action(func() {
		for key := range m.Sessions {
			session := m.Sessions[key]
			if session.IsReady {
				sessions[session.ID] = true
			}
		}
	})
	return
}

// GetNumberOfReadyAndChatting gets number of sessions ready and already chatting
func (m *Map) GetNumberOfReadyAndChatting() (ready int, chatting int) {
	m.Action(func() {
		for key := range m.Sessions {
			session := m.Sessions[key]
			if session.IsReady {
				ready++
			}
			if session.Room != nil {
				chatting++
			}
		}
	})

	return
}

// Close closes defined session
func (m *Map) Close(id string) {
	m.Action(func() {
		delete(m.Sessions, id)
	})
}

// Action sends action to be perform on the sessions
func (m *Map) Action(fn func()) {
	wait := make(chan bool)
	m.actions <- func() {
		fn()
		wait <- true
	}
	<-wait
}
