package jms

import "sync"

type SessionManager struct {
	store sync.Map
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) RegisterJMSSession(jmsSession *JMSSession) string {
	sessionID := jmsSession.Session.Id
	sm.store.Store(sessionID, jmsSession)
	return sessionID
}

func (sm *SessionManager) UnregisterJMSSession(jmsSession *JMSSession) {
	sessionID := jmsSession.Session.Id
	sm.store.Delete(sessionID)
}

func (sm *SessionManager) GetJMSSession(sessionID string) *JMSSession {
	if value, ok := sm.store.Load(sessionID); ok {
		if jmsSession, ok := value.(*JMSSession); ok {
			return jmsSession
		}
	}
	return nil
}

func (sm *SessionManager) GetStore() sync.Map {
	return sm.store
}
