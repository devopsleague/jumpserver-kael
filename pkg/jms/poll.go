package jms

import (
	"fmt"
	"github.com/jumpserver/wisp/pkg/logger"
	"path/filepath"
	"sync"
)

type PollJMSEvent struct {
	mutex sync.Mutex
}

func NewPollJMSEvent() *PollJMSEvent {
	return &PollJMSEvent{}
}

func (p *PollJMSEvent) CloseSession(targetSession *JMSSession) {
	// Implement the logic to close the targetSession
	// You may need to call the Close() method of JMSSession
}

func (p *PollJMSEvent) ClearZombieSession() {
	// Implement the logic to clear zombie sessions
	// You may need to use the appropriate APIs for scanning and removing zombie sessions
	replayDir := filepath.Join(globals.PROJECT_DIR, "data/replay")
	req := &RemainReplayRequest{
		ReplayDir: replayDir,
	}
	resp, err := p.Stub.ScanRemainReplays(context.Background(), req)
	if err != nil || !resp.Status.Ok {
		errorMessage := fmt.Sprintf("Failed to scan remain replay: %s", resp.Status.Err)
		// Handle the error
	} else {
		logger.Info("Scan remain replay success")
	}
}

func (p *PollJMSEvent) WaitForKillSessionMessage() {
	from
	api.jms.session import
	SessionManager
	q := make(chan *YourResponseType, 1000)
	go func() {
		for resp := range q {
			task := resp.Task
			taskID := task.ID
			sessionID := task.SessionID
			taskAction := task.Action
			targetSession := SessionManager.GetJMSSession(sessionID)
			if targetSession != nil {
				if taskAction == KillSession {
					p.CloseSession(targetSession)
				}
				req := &FinishedTaskRequest{
					TaskID: targetSession.Session.ID,
				}
				p.Stub.FinishSession(context.Background(), req)
			}
		}
	}()
	p.Stub.DispatchTask(q)
}

func (p *PollJMSEvent) StartSessionKiller() {
	p.WaitForKillSessionMessage()
}

func (p *PollJMSEvent) Start() {
	p.ClearZombieSession()
	p.StartSessionKiller()
}

func setupPollJMSEvent() {
	jmsEvent := NewPollJMSEvent()
	go jmsEvent.Start()
}
