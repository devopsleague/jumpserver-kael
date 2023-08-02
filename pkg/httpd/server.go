package httpd

import (
	"context"
	"fmt"
	"github.com/jumpserver/kael/pkg/config"
	"github.com/jumpserver/kael/pkg/httpd/router"
	"github.com/jumpserver/kael/pkg/logger"
	"net"
	"net/http"
	"time"
)

func NewServer() *Server {
	srv := &Server{}
	eng := router.CreateRouter()
	conf := config.GlobalConfig
	addr := net.JoinHostPort(conf.Host, conf.Port)
	srv.Srv = &http.Server{Addr: addr, Handler: eng}
	return srv
}

type Server struct {
	Srv *http.Server
}

func (s *Server) Start() {
	logger.GlobalLogger.Info(
		fmt.Sprintf("Start HTTP Server at %s", s.Srv.Addr),
	)
	fmt.Println(s.Srv.ListenAndServe())
}

func (s *Server) Stop() {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelFunc()
	if s.Srv != nil {
		_ = s.Srv.Shutdown(ctx)
	}
}
