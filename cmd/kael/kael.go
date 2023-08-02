package main

import (
	"flag"
	"github.com/jumpserver/kael/pkg/config"
	"github.com/jumpserver/kael/pkg/httpd"
	"github.com/jumpserver/kael/pkg/httpd/grpc"
	"github.com/jumpserver/kael/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
var (
	configPath = ""
)

func init() {
	flag.StringVar(&configPath, "f", "config.yml", "config.yml path")
}

type Kael struct {
	webSrv     *httpd.Server
	grpcClient *grpc.Client
}

func (k *Kael) Start() {
	go k.webSrv.Start()
	k.grpcClient.Start()
}

func (k *Kael) Stop() {
	k.webSrv.Stop()
	k.grpcClient.Stop()
}

func main() {
	flag.Parse()
	config.Setup(configPath)
	logger.Setup()
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	app := &Kael{
		webSrv:     httpd.NewServer(),
		grpcClient: grpc.GlobalGrpcClient,
	}
	app.Start()
	<-gracefulStop
	app.Stop()
}
