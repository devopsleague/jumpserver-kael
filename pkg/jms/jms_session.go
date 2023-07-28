package jms

import (
	"github.com/gorilla/websocket"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"time"
)

type JMSSession struct {
	Session             *protobuf.Session
	Websocket           *websocket.Conn
	HistoryAsks         []*AskResponse
	CurrentAskInterrupt bool
	CommandACLs         []*protobuf.CommandACL
	ExpireTime          time.Time
	MaxIdleTimeDelta    int
	SessionHandler      *SessionHandler
	CommandHandler      *CommandHandler
	ReplayHandler       *ReplayHandler
	JMSState            *JMSState
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
	jmss.SessionHandler.closeSession(jmss.Session)
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
