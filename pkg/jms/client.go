package jms

import (
	"fmt"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient() *GrpcClient {
	return &GrpcClient{}
}

type GrpcClient struct {
	Conn   *grpc.ClientConn
	Client pb.ServiceClient
}

func (c *GrpcClient) Start() {
	conn, err := grpc.Dial(
		"localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println("grpc client start error", err)
		return
	}
	client := pb.NewServiceClient(conn)

	c.Conn = conn
	c.Client = client
}

func (c *GrpcClient) Stop() {
	_ = c.Conn.Close()
}
