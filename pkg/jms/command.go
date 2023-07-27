package jms

import (
	"context"
	"github.com/gorilla/websocket"
	"regexp"
	"strings"
	"time"
)

type CommandHandler struct {
	Session       *Session
	Websocket     *websocket.Conn
	CommandACLs   []*CommandACL
	CmdACLID      string
	CmdGroupID    string
	CommandRecord *CommandRecord
	JMSState      *JMSState
	Stub          *JMSStub
}

const (
	WAIT_TICKET_TIMEOUT  = 60 * 3
	WAIT_TICKET_INTERVAL = 2
)

func NewCommandHandler(websocket *websocket.Conn, session *Session, commandACLs []*CommandACL, jmsState *JMSState) *CommandHandler {
	return &CommandHandler{
		Session:       session,
		Websocket:     websocket,
		CommandACLs:   commandACLs,
		CommandRecord: nil,
		JMSState:      jmsState,
		Stub:, // Initialize your JMSStub instance here,
	}
}

func (ch *CommandHandler) RecordCommand() {
	req := &CommandRequest{
		Sid:        ch.Session.ID,
		OrgID:      ch.Session.OrgID,
		Asset:      ch.Session.Asset,
		Account:    ch.Session.Account,
		User:       ch.Session.User,
		Timestamp:  time.Now().Unix(),
		Input:      ch.CommandRecord.Input,
		Output:     ch.CommandRecord.Output,
		RiskLevel:  ch.CommandRecord.RiskLevel,
		CmdACLID:   ch.CmdACLID,
		CmdGroupID: ch.CmdGroupID,
	}

	resp, err := ch.Stub.UploadCommand(context.Background(), req)
	if err != nil || !resp.Status.Ok {
		errorMessage := "Failed to upload command"
		// Handle the error
	}
}

func (ch *CommandHandler) MatchRule() *CommandACL {
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
				ch.CmdACLID = commandACL.ID
				ch.CmdGroupID = commandGroup.ID
				return commandACL
			}
		}
	}

	return nil
}

func (ch *CommandHandler) CreateAndWaitTicket(commandACL *CommandACL) bool {
	req := &CommandConfirmRequest{
		Cmd:       ch.CommandRecord.Input,
		SessionID: ch.Session.ID,
		CmdACLID:  ch.CmdACLID,
	}

	resp, err := ch.Stub.CreateCommandTicket(context.Background(), req)
	if err != nil || !resp.Status.Ok {
		errorMessage := "Failed to create ticket"
		// Handle the error
	}

	return ch.WaitForTicketStatusChange(resp.Info)
}

func (ch *CommandHandler) WaitForTicketStatusChange(ticketInfo *TicketInfo) bool {
	// Implement the function here as you did in the Python code.
	// The logic for waiting and checking the ticket status can be similar.

	return true
}

func (ch *CommandHandler) CommandACLFilter() bool {
	isContinue := true
	acl := ch.MatchRule()
	if acl != nil {
		switch acl.Action {
		case CommandACLReject:
			isContinue = false
			ch.CommandRecord.RiskLevel = RiskLevelReject
			// Handle response as you did in the Python code
		case CommandACLReview:
			isContinue = false
			startTime := time.Now()
			endTime := startTime.Add(time.Duration(60) * time.Second)
			// Handle response as you did in the Python code
		case CommandACLWarning:
			ch.CommandRecord.RiskLevel = RiskLevelWarning
		}
	}
	return isContinue
}

func (ch *CommandHandler) CloseTicket(ticketInfo *TicketInfo) {
	req := &TicketRequest{
		Req: ticketInfo.CancelReq,
	}

	resp, err := ch.Stub.CancelTicket(context.Background(), req)
	if err != nil || !resp.Status.Ok {
		errorMessage := "Failed to close ticket"
		// Handle the error
	}
}
