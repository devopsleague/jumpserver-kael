package jms

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/global"
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
	CommandRecord *CommandRecord
	JMSState      *JMSState
}

const (
	WAIT_TICKET_TIMEOUT  = 60 * 3
	WAIT_TICKET_INTERVAL = 2
)

func NewCommandHandler(
	websocket *websocket.Conn, session *protobuf.Session,
	commandACLs []*protobuf.CommandACL, jmsState *JMSState,
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
			flags := regexp.UNICODE
			if commandGroup.IgnoreCase {
				flags |= regexp.IGNORECASE
			}

			pattern, err := regexp.Compile("(?" + flags + ")" + commandGroup.Pattern)
			if err != nil {
				errorMessage := "Failed to compile regular expression"
				// Handle the error
			}

			if pattern.MatchString(strings.ToLower(ch.CommandRecord.Input)) {
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
	}

	return ch.WaitForTicketStatusChange(resp.Info)
}

func (ch *CommandHandler) WaitForTicketStatusChange(ticketInfo *protobuf.TicketInfo) bool {
	// Implement the function here as you did in the Python code.
	// The logic for waiting and checking the ticket status can be similar.

	return true
}

func (ch *CommandHandler) CommandACLFilter() bool {
	isContinue := true
	acl := ch.MatchRule()
	if acl != nil {
		switch acl.Action {
		case protobuf.CommandACL_Reject:
			isContinue = false
			ch.CommandRecord.RiskLevel = protobuf.CommandACL_Reject
		case protobuf.CommandACL_Review:
			isContinue = false
			startTime := time.Now()
			endTime := startTime.Add(time.Duration(60) * time.Second)
		case protobuf.CommandACL_Warning:
			ch.CommandRecord.RiskLevel = protobuf.CommandACL_Warning
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
