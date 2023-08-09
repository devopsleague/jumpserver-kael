package jms

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/httpd/grpc"
	"github.com/jumpserver/kael/pkg/logger"
	"github.com/jumpserver/kael/pkg/schemas"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"time"
)

type CommandHandler struct {
	Session       *protobuf.Session
	Websocket     *websocket.Conn
	CommandACLs   []*protobuf.CommandACL
	CmdACLID      string
	CmdGroupID    string
	CommandRecord *schemas.CommandRecord
	JMSState      *schemas.JMSState
}

const (
	WAIT_TICKET_TIMEOUT  = 60 * 3
	WAIT_TICKET_INTERVAL = 2
)

func NewCommandHandler(
	websocket *websocket.Conn, session *protobuf.Session,
	commandACLs []*protobuf.CommandACL, jmsState *schemas.JMSState,
) *CommandHandler {
	return &CommandHandler{
		Session:       session,
		Websocket:     websocket,
		CommandACLs:   commandACLs,
		CommandRecord: nil,
		JMSState:      jmsState,
	}
}

func (ch *CommandHandler) RecordCommand() {
	ctx := context.Background()
	req := &protobuf.CommandRequest{
		Sid:        ch.Session.Id,
		OrgId:      ch.Session.OrgId,
		Asset:      ch.Session.Asset,
		Account:    ch.Session.Account,
		User:       ch.Session.User,
		Timestamp:  time.Now().Unix(),
		Input:      ch.CommandRecord.Input,
		Output:     ch.CommandRecord.Output,
		RiskLevel:  ch.CommandRecord.RiskLevel,
		CmdAclId:   ch.CmdACLID,
		CmdGroupId: ch.CmdGroupID,
	}

	resp, _ := grpc.GlobalGrpcClient.Client.UploadCommand(ctx, req)
	if !resp.Status.Ok {
		logger.GlobalLogger.Error("Failed to upload command")
	}
}

func (ch *CommandHandler) MatchRule() *protobuf.CommandACL {
	for _, commandACL := range ch.CommandACLs {
		for _, commandGroup := range commandACL.CommandGroups {
			regexp2Opt := regexp2.None
			if commandGroup.IgnoreCase {
				regexp2Opt = regexp2.IgnoreCase
			}

			re, err := regexp2.Compile(commandGroup.Pattern, regexp2Opt)
			if err != nil {
				logger.GlobalLogger.Error("Failed to compile regular expression")
				return nil
			}

			found, err := re.FindStringMatch(ch.CommandRecord.Input)
			if err != nil || found == nil {
				continue
			}

			ch.CmdACLID = commandACL.Id
			ch.CmdGroupID = commandGroup.Id
			return commandACL
		}
	}
	return nil
}

func (ch *CommandHandler) CreateAndWaitTicket(commandACL *protobuf.CommandACL) bool {
	ctx := context.Background()
	req := &protobuf.CommandConfirmRequest{
		Cmd:       ch.CommandRecord.Input,
		SessionId: ch.Session.Id,
		CmdAclId:  ch.CmdACLID,
	}

	resp, _ := grpc.GlobalGrpcClient.Client.CreateCommandTicket(ctx, req)
	if !resp.Status.Ok {
		logger.GlobalLogger.Error("Failed to create ticket")
		return false
	}

	return ch.WaitForTicketStatusChange(resp.Info)
}

func (ch *CommandHandler) WaitForTicketStatusChange(ticketInfo *protobuf.TicketInfo) bool {
	response := &schemas.AskResponse{
		Type:           schemas.Waiting,
		ConversationID: ch.Session.Id,
		SystemMessage: fmt.Sprintf(
			"复核请求已发起，请等待复核: [查看结果](%s): ",
			ticketInfo.TicketDetailUrl,
		),
	}

	jsonResponse, _ := json.Marshal(response)
	_ = ch.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)

	ctx := context.Background()
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(WAIT_TICKET_TIMEOUT) * time.Second)
	ticketClosed := true
	isContinue := false
	for time.Now().Before(endTime) {
		req := &protobuf.TicketRequest{Req: ticketInfo.CheckReq}
		resp, _ := grpc.GlobalGrpcClient.Client.CheckTicketState(ctx, req)
		if !resp.Status.Ok {
			logger.GlobalLogger.Error("Failed to check ticket status")
			break
		}
		switch resp.Data.State {
		case protobuf.TicketState_Approved:
			isContinue = true
			ticketClosed = false
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel_ReviewAccept
			break
		case protobuf.TicketState_Rejected:
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel_ReviewReject
			ticketClosed = false

		case protobuf.TicketState_Closed:
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel_ReviewCancel
			ticketClosed = false

		default:
			time.Sleep(WAIT_TICKET_INTERVAL)
		}
	}

	if ticketClosed {
		ch.CloseTicket(ticketInfo)
	}
	return isContinue
}

func (ch *CommandHandler) CommandACLFilter() bool {
	isContinue := true
	acl := ch.MatchRule()
	if acl != nil {
		switch acl.Action {
		case protobuf.CommandACL_Reject:
			response := &schemas.AskResponse{
				Type:           schemas.Reject,
				ConversationID: ch.Session.Id,
				SystemMessage:  "对话已被拒绝",
			}

			jsonResponse, _ := json.Marshal(response)
			_ = ch.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)

			isContinue = false
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel(protobuf.CommandACL_Reject)
		case protobuf.CommandACL_Warning:
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel(protobuf.CommandACL_Warning)
		case protobuf.CommandACL_Review:
			response := &schemas.AskResponse{
				Type:           schemas.Waiting,
				ConversationID: ch.Session.Id,
				SystemMessage:  "您输入命令需要复核后才可以执行，是否发起复核请求？",
				Meta:           schemas.ResponseMeta{ActivateReview: true},
			}

			jsonResponse, _ := json.Marshal(response)
			_ = ch.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)

			isContinue = false
			startTime := time.Now()
			endTime := startTime.Add(time.Duration(60) * time.Second)

			for time.Now().Before(endTime) {
				switch ch.JMSState.ActivateReview {
				case schemas.Wait:
					time.Sleep(1 * time.Second)
				case schemas.Rejected:
					isContinue = false
					ch.JMSState.ActivateReview = schemas.Wait
					break
				case schemas.Approve:
					ch.JMSState.ActivateReview = schemas.Wait
					isContinue = ch.CreateAndWaitTicket(acl)
					break
				}
			}
		}
	}
	return isContinue
}

func (ch *CommandHandler) CloseTicket(ticketInfo *protobuf.TicketInfo) {
	ctx := context.Background()
	req := &protobuf.TicketRequest{
		Req: ticketInfo.CancelReq,
	}

	resp, _ := grpc.GlobalGrpcClient.Client.CancelTicket(ctx, req)
	if !resp.Status.Ok {
		logger.GlobalLogger.Error("Failed to close ticket")
	}
}
