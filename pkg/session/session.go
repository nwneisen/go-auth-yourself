package session

import (
	"sync"

	"github.com/google/uuid"
)

type Session struct {
	ID        string
	UserInfo  map[string]string
	CreatedAt int64
	ExpiresAt int64
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession(userID string, userInfo map[string]string) *Session {
	now := 0
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		ID:        uuid.New().String(),
		UserInfo:  userInfo,
		CreatedAt: int64(now),
		ExpiresAt: int64(now) + 3600,
	}

	sm.sessions[session.ID] = session
	return session
}

func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, ok := sm.sessions[id]
	if !ok {
		return nil, false
	}

	if session.ExpiresAt < int64(0) {
		return nil, false
	}

	return session, true
}

func (sm *SessionManager) DeleteSession(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, id)
}

func (sm *SessionManager) DeleteExpired() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := int64(0)
	for id, session := range sm.sessions {
		if session.ExpiresAt < now {
			delete(sm.sessions, id)
		}
	}
}
