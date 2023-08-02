package jms

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/httpd/grpc"
	"github.com/jumpserver/kael/pkg/logger"
	"github.com/jumpserver/kael/pkg/schemas"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"time"
)

type SessionHandler struct {
	Websocket     *websocket.Conn
	RemoteAddress string
}

func NewSessionHandler(websocket *websocket.Conn) *SessionHandler {
	return &SessionHandler{
		Websocket:     websocket,
		RemoteAddress: getRemoteAddress(websocket),
	}
}

func getRemoteAddress(websocket *websocket.Conn) string {
	return websocket.RemoteAddr().String()
}

func (sh *SessionHandler) CreateNewSession(authInfo *protobuf.TokenAuthInfo) *JMSSession {
	session := sh.createSession(authInfo)
	return &JMSSession{
		Session:             session,
		Websocket:           sh.Websocket,
		HistoryAsks:         make([]string, 0),
		CurrentAskInterrupt: false,
		CommandACLs:         authInfo.FilterRules,
		ExpireTime:          time.Unix(authInfo.ExpireInfo.ExpireAt, 0),
		MaxIdleTimeDelta:    int(authInfo.Setting.MaxIdleTime),
		SessionHandler:      sh,
		CommandHandler:      nil,
		ReplayHandler:       nil,
		JMSState:            &schemas.JMSState{ID: session.Id},
	}
}

func (sh *SessionHandler) createSession(authInfo *protobuf.TokenAuthInfo) *protobuf.Session {
	reqSession := &protobuf.Session{
		UserId:     authInfo.User.Id,
		User:       fmt.Sprintf("%s(%s)", authInfo.User.Name, authInfo.User.Username),
		AccountId:  authInfo.Account.Id,
		Account:    fmt.Sprintf("%s(%s)", authInfo.Account.Name, authInfo.Account.Username),
		OrgId:      authInfo.Asset.OrgId,
		AssetId:    authInfo.Asset.Id,
		Asset:      authInfo.Asset.Name,
		LoginFrom:  protobuf.Session_WT,
		Protocol:   authInfo.Asset.Protocols[0].Name,
		DateStart:  time.Now().Unix(),
		RemoteAddr: sh.RemoteAddress,
	}
	ctx := context.Background()
	req := &protobuf.SessionCreateRequest{
		Data: reqSession,
	}

	resp, _ := grpc.GlobalGrpcClient.Client.CreateSession(ctx, req)
	if !resp.Status.Ok {
		errorMessage := fmt.Sprintf("Failed to create session: %s", resp.Status.Err)
		logger.GlobalLogger.Error(errorMessage)
	}
	return resp.Data
}

func (sh *SessionHandler) closeSession(session *protobuf.Session) {
	ctx := context.Background()
	req := &protobuf.SessionFinishRequest{
		Id:      session.Id,
		DateEnd: time.Now().Unix(),
	}

	resp, _ := grpc.GlobalGrpcClient.Client.FinishSession(ctx, req)
	if !resp.Status.Ok {
		errorMessage := fmt.Sprintf("Failed to close session: %s", resp.Status.Err)
		logger.GlobalLogger.Error(errorMessage)
	}
}
