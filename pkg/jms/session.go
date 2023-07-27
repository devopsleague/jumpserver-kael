package jms

import "sync"

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type JMSSession struct {
	Session             *Session
	Websocket           *websocket.Conn
	HistoryAsks         []*AskResponse
	CurrentAskInterrupt bool
	CommandACLs         []*CommandACL
	ExpireTime          time.Time
	MaxIdleTimeDelta    int
	SessionHandler      *SessionHandler
	CommandHandler      *CommandHandler
	ReplayHandler       *ReplayHandler
	JMSState            *JMSState
}

type SessionHandler struct {
	Websocket     *websocket.Conn
	RemoteAddress string
}

func NewSessionHandler(websocket *websocket.Conn) *SessionHandler {
	return &SessionHandler{
		Websocket: websocket,
	}
}

func (sh *SessionHandler) GetRemoteAddress() string {
	// Implement the logic to get the remote address from the websocket
	remoteAddress := ""
	if sh.Websocket != nil {
		// Logic to extract the remote address from the websocket connection
	}
	return remoteAddress
}

func (sh *SessionHandler) CreateNewSession(authInfo *TokenAuthInfo) *JMSSession {
	session := sh.CreateSession(authInfo)
	return &JMSSession{
		Session:             session,
		Websocket:           sh.Websocket,
		HistoryAsks:         []*AskResponse{},
		CurrentAskInterrupt: false,
		CommandACLs:         authInfo.FilterRules,
		ExpireTime:          time.Unix(authInfo.ExpireInfo.ExpireAt, 0),
		MaxIdleTimeDelta:    authInfo.Setting.MaxIdleTime,
		SessionHandler:      sh,
		CommandHandler:      nil,
		ReplayHandler:       nil,
		JMSState:            &JMSState{ID: session.ID},
	}
}

func (sh *SessionHandler) CreateSession(authInfo *TokenAuthInfo) *Session {
	// Implement the logic to create a new session based on authInfo
	// You may need to use your JMSStub instance to send the request and receive the response
	reqSession := &Session{
		UserID:     authInfo.User.ID,
		User:       fmt.Sprintf("%s(%s)", authInfo.User.Name, authInfo.User.Username),
		AccountID:  authInfo.Account.ID,
		Account:    fmt.Sprintf("%s(%s)", authInfo.Account.Name, authInfo.Account.Username),
		OrgID:      authInfo.Asset.OrgID,
		AssetID:    authInfo.Asset.ID,
		Asset:      authInfo.Asset.Name,
		LoginFrom:  Session_WT,
		Protocol:   authInfo.Asset.Protocols[0].Name,
		DateStart:  time.Now().Unix(),
		RemoteAddr: sh.GetRemoteAddress(),
	}
	req := &SessionCreateRequest{
		Data: reqSession,
	}

	resp, err := sh.Stub.CreateSession(context.Background(), req)
	if err != nil || !resp.Status.Ok {
		errorMessage := fmt.Sprintf("Failed to create session: %s", resp.Status.Err)
		// Handle the error
	}
	return resp.Data
}

func (jmss *JMSSession) ActiveSession() {
	SessionManager.RegisterJMSSession(jmss)
	jmss.ReplayHandler = NewReplayHandler(jmss.Session)
	jmss.CommandHandler = NewCommandHandler(jmss.Websocket, jmss.Session, jmss.CommandACLs, jmss.JMSState)
	go jmss.MaximumIdleTimeDetection()
}

func (jmss *JMSSession) MaximumIdleTimeDetection() {
	lastActiveTime := time.Now()

	for {
		currentTime := time.Now()
		idleTime := currentTime.Sub(lastActiveTime)

		if idleTime.Seconds() >= float64(jmss.MaxIdleTimeDelta*60) {
			jmss.Close()
			break
		}

		if jmss.JMSState.NewDialogue {
			lastActiveTime = currentTime
			jmss.JMSState.NewDialogue = false
		}

		time.Sleep(3 * time.Second)
	}
}

func (jmss *JMSSession) Close() {
	jmss.CurrentAskInterrupt = true
	time.Sleep(1 * time.Second)
	jmss.ReplayHandler.Upload()
	jmss.SessionHandler.CloseSession(jmss.Session)
	SessionManager.UnregisterJMSSession(jmss)
	jmss.NotifyToClose()
}

func (jmss *JMSSession) NotifyToClose() {
	// Implement the logic to notify to close the session
	// You may need to use your ws package to send the response
	response := &AskResponse{
		Type:           AskResponseTypeFinish,
		ConversationID: jmss.Session.ID,
		SystemMessage:  "Session interrupted",
	}
	reply(jmss.Websocket, response)
}

func (jmss *JMSSession) WithAudit(command string, chatFunc func(*JMSSession) string) (result string) {
	// Implement the logic to perform auditing
	// You may need to call the chatFunc here and handle the CommandRecord, replay writing, and other tasks accordingly.
	return result
}

// Implement other methods of JMSSession as needed.

type SessionManager struct {
	store sync.Map
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) RegisterJMSSession(jmsSession *JMSSession) string {
	sessionID := jmsSession.Session.ID
	sm.store.Store(sessionID, jmsSession)
	return sessionID
}

func (sm *SessionManager) UnregisterJMSSession(jmsSession *JMSSession) {
	sessionID := jmsSession.Session.ID
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
