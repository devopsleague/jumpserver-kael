package main

import (
	"flag"
	"fmt"
	"github.com/jumpserver/kael/pkg/config"
	"github.com/jumpserver/kael/pkg/global"
	"github.com/jumpserver/kael/pkg/httpd"
	"github.com/jumpserver/kael/pkg/logger"
	"net"
	"net/http"
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
func main() {
	config.Setup(configPath)
	logger.Setup()
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	eng := httpd.Routers()
	conf := global.Config
	addr := net.JoinHostPort(conf.Host, conf.Port)
	srv := &http.Server{Addr: addr, Handler: eng}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("gin error: ", err)
		return
	}
	<-gracefulStop
	//srv.Shutdown()

}
