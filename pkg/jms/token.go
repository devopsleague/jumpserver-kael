package jms

import (
	"context"
	"errors"
	"github.com/jumpserver/kael/pkg/httpd/grpc"
	"github.com/jumpserver/kael/pkg/logger"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
)

type TokenHandler struct{}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

func (th *TokenHandler) GetTokenAuthInfo(token string) (*protobuf.TokenAuthInfo, error) {
	ctx := context.Background()
	req := &protobuf.TokenRequest{Token: token}

	resp, _ := grpc.GlobalGrpcClient.Client.GetTokenAuthInfo(ctx, req)
	if !resp.Status.Ok {
		errorMessage := "Failed to get token: " + resp.Status.Err
		logger.GlobalLogger.Error(errorMessage)
		return nil, errors.New(errorMessage)
	}
	return resp.Data, nil
}
