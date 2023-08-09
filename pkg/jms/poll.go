package jms

import (
	"context"
	"fmt"
	"github.com/jumpserver/kael/pkg/config"
	"github.com/jumpserver/kael/pkg/httpd/grpc"
	"github.com/jumpserver/kael/pkg/logger"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"go.uber.org/zap"
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

	resp, err := grpc.GlobalGrpcClient.Client.ScanRemainReplays(ctx, req)
	if err != nil || !resp.Status.Ok {
		logger.GlobalLogger.Error("Failed to scan remain replay")
	} else {
		logger.GlobalLogger.Info("Scan remain replay success")
	}
}

func (p *PollJMSEvent) WaitForKillSessionMessage() {
	stream, err := grpc.GlobalGrpcClient.Client.DispatchTask(context.Background())
	if err != nil {
		logger.GlobalLogger.Error("dispatch task err", zap.Error(err))
		return
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
			logger.GlobalLogger.Error("Failed to receive a note", zap.Error(err))
			continue
		}

		task := taskResponse.Task
		sessionId := task.SessionId
		taskAction := task.Action
		targetSession := GlobalSessionManager.GetJMSSession(sessionId)
		if targetSession != nil {
			if taskAction == protobuf.TaskAction_KillSession {
				targetSession.Close()
			}
			req := &protobuf.SessionFinishRequest{
				Id: task.SessionId,
			}

			resp, _ := grpc.GlobalGrpcClient.Client.FinishSession(context.Background(), req)
			if !resp.Status.Ok {
				errorMessage := fmt.Sprintf("Failed to finish session: %s", resp.Status.Err)
				logger.GlobalLogger.Error(errorMessage)
			}
		}
	}
	<-waitChan
}

func (p *PollJMSEvent) Start() {
	p.ClearZombieSession()
	p.WaitForKillSessionMessage()
}

func SetupPollJMSEvent() {
	jmsEvent := NewPollJMSEvent()
	go jmsEvent.Start()
}
