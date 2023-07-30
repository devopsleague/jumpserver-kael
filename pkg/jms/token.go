package jms

import (
	"context"
	"fmt"
	"github.com/jumpserver/kael/pkg/global"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"log"
)

type TokenHandler struct{}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

func (th *TokenHandler) GetTokenAuthInfo(token string) (*protobuf.TokenAuthInfo, error) {
	ctx := context.Background()
	req := &protobuf.TokenRequest{Token: token}

	resp, err := global.GrpcClient.Client.GetTokenAuthInfo(ctx, req)
	if err != nil {
		fmt.Println("Failed to get token")
		return nil, err
	}
	if !resp.Status.Ok {
		errorMessage := "Failed to get token: " + resp.Status.Err
		log.Printf(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}
	return resp.Data, nil
}
