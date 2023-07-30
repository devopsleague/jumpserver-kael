package jms

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/global"
	"github.com/jumpserver/kael/pkg/schemas"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"regexp"
	"strings"
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

	resp, err := global.GrpcClient.Client.UploadCommand(ctx, req)
	if err != nil || !resp.Status.Ok {
		errorMessage := "Failed to upload command"
		fmt.Println(errorMessage)
	}
}

func (ch *CommandHandler) MatchRule() *protobuf.CommandACL {
	for _, commandACL := range ch.CommandACLs {
		for _, commandGroup := range commandACL.CommandGroups {
			flags := ""
			if commandGroup.IgnoreCase {
				flags = "(?i)"
			}
			re, err := regexp.Compile(flags + commandGroup.Pattern)
			if err != nil {
				errorMessage := "Failed to compile regular expression"
				fmt.Println(errorMessage)
				return nil
			}

			if re.MatchString(strings.ToLower(ch.CommandRecord.Input)) {
				ch.CmdACLID = commandACL.Id
				ch.CmdGroupID = commandGroup.Id
				return commandACL
			}
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

	resp, err := global.GrpcClient.Client.CreateCommandTicket(ctx, req)
	if err != nil || !resp.Status.Ok {
		errorMessage := "Failed to create ticket"
		fmt.Println(errorMessage)
		return false
	}

	return ch.WaitForTicketStatusChange(resp.Info)
}

func (ch *CommandHandler) WaitForTicketStatusChange(ticketInfo *protobuf.TicketInfo) bool {
	ctx := context.Background()
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(WAIT_TICKET_TIMEOUT) * time.Second)
	ticketClosed := true
	isContinue := false
	for time.Now().Before(endTime) {
		req := &protobuf.TicketRequest{Req: ticketInfo.CheckReq}
		resp, err := global.GrpcClient.Client.CheckTicketState(ctx, req)
		if err != nil || !resp.Status.Ok {
			errorMessage := "Failed to check ticket status"
			fmt.Println(errorMessage)
			break
		}
		systemMessage := ""
		switch resp.Data.State {
		case protobuf.TicketState_Approved:
			isContinue = true
			ticketClosed = false
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel_ReviewAccept
			break
		case protobuf.TicketState_Rejected:
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel_ReviewReject
			ticketClosed = false
			systemMessage = "The ticket is rejected"
		case protobuf.TicketState_Closed:
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel_ReviewCancel
			ticketClosed = false
			systemMessage = "The ticket is closed"
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
			isContinue = false
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel(protobuf.CommandACL_Reject)
		case protobuf.CommandACL_Warning:
			ch.CommandRecord.RiskLevel = protobuf.RiskLevel(protobuf.CommandACL_Warning)
		case protobuf.CommandACL_Review:
			isContinue = false
			startTime := time.Now()
			endTime := startTime.Add(time.Duration(60) * time.Second)

			for time.Now().Before(endTime) {
				switch ch.JMSState.ActivateReview {
				case 0:
					time.Sleep(1 * time.Second)
				case 1:
					// 复核
					isContinue = ch.CreateAndWaitTicket(acl)
				case 2:
					// 拒绝
					//ch.JMSState.ActivateReview 还原
					isContinue = false
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

	resp, err := global.GrpcClient.Client.CancelTicket(ctx, req)
	if err != nil || !resp.Status.Ok {
		errorMessage := "Failed to close ticket"
		fmt.Println(errorMessage)
	}
}
