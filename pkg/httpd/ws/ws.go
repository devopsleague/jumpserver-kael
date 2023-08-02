package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/logger"
	"net/http"
)

const (
	defaultBufferSize = 1024
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  defaultBufferSize,
	WriteBufferSize: defaultBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func UpgradeWsConn(ctx *gin.Context) (*websocket.Conn, error) {
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header())
	if err != nil {
		return nil, err
	}

	conn.SetPingHandler(func(appData string) error {
		logger.GlobalLogger.Info(
			fmt.Sprintf("Accept ping message: %s", appData),
		)
		return conn.WriteMessage(websocket.PongMessage, []byte("Pong"))
	})
	conn.SetPongHandler(func(appData string) error {
		logger.GlobalLogger.Info(
			fmt.Sprintf("Accept pong message: %s", appData),
		)
		return conn.WriteMessage(websocket.PongMessage, []byte("Ping"))
	})
	conn.RemoteAddr()
	return conn, nil
}
