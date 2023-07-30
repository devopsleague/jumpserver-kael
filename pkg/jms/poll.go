package jms

import (
	"context"
	"fmt"
	"github.com/jumpserver/kael/pkg/config"
	"github.com/jumpserver/kael/pkg/global"
	"github.com/jumpserver/wisp/pkg/logger"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"io"
)

type PollJMSEvent struct{}

func NewPollJMSEvent() *PollJMSEvent {
	return &PollJMSEvent{}
}

func (p *PollJMSEvent) ClearZombieSession() {
	ctx := context.Background()
	req := &protobuf.RemainReplayRequest{
		ReplayDir: config.GlobalConfig.ReplayFolderPath,
	}

	resp, err := global.GrpcClient.Client.ScanRemainReplays(ctx, req)
	if err != nil || !resp.Status.Ok {
		errorMessage := fmt.Sprintf("Failed to scan remain replay")
		fmt.Println(errorMessage)
	} else {
		logger.Info("Scan remain replay success")
	}
}

func (p *PollJMSEvent) WaitForKillSessionMessage() {
	stream, err := global.GrpcClient.Client.DispatchTask(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	waitChan := make(chan struct{})
	for {
		taskResponse, err := stream.Recv()
		if err == io.EOF {
			_ = stream.CloseSend()
			close(waitChan)
			break
		}
		if err != nil {
			fmt.Printf("Failed to receive a note : %s\n", err)
		}

		task := taskResponse.Task
		sessionId := task.SessionId
		taskAction := task.Action
		targetSession := global.SessionManager.GetJMSSession(sessionId)
		if targetSession != nil {
			if taskAction == protobuf.TaskAction_KillSession {
				targetSession.Close()
			}
			req := &protobuf.SessionFinishRequest{
				Id: task.SessionId,
			}

			resp, _ := global.GrpcClient.Client.FinishSession(context.Background(), req)
			if !resp.Status.Ok {
				errorMessage := fmt.Sprintf("Failed to finish session: %s", resp.Status.Err)
				fmt.Println(errorMessage)
			}
		}
	}
	<-waitChan
}

func (p *PollJMSEvent) Start() {
	p.ClearZombieSession()
	p.WaitForKillSessionMessage()
}

func setupPollJMSEvent() {
	jmsEvent := NewPollJMSEvent()
	go jmsEvent.Start()
}
