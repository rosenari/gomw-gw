package services

import (
	"sync"

	"gomw-gw/app/internal/models"
)

type SessionManager struct {
	sessions sync.Map
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) AddSession(session *models.Session) {
	sm.sessions.Store(session.ID, session)
}

func (sm *SessionManager) GetSession(connectionID models.ConnectionID) (*models.Session, bool) {
	value, exists := sm.sessions.Load(connectionID)
	if !exists {
		return nil, false
	}
	
	session, ok := value.(*models.Session)
	return session, ok
}

func (sm *SessionManager) RemoveSession(connectionID models.ConnectionID) {
	if session, exists := sm.GetSession(connectionID); exists {
		session.Close()
		sm.sessions.Delete(connectionID)
	}
}

func (sm *SessionManager) GetAllSessions() []*models.Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	var sessions []*models.Session
	sm.sessions.Range(func(key, value interface{}) bool {
		if session, ok := value.(*models.Session); ok {
			sessions = append(sessions, session)
		}
		return true
	})
	
	return sessions
}

func (sm *SessionManager) GetSessionCount() int {
	count := 0
	sm.sessions.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
} 