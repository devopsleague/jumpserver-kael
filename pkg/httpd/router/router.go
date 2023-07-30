package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jumpserver/kael/pkg/httpd/middlewares"
)

func CreateRouter() *gin.Engine {
	eng := gin.Default()
	eng.Use(middlewares.CORSMiddleware())
	karlGroup := eng.Group("/kael")
	karlGroup.GET("/health/", HealthApi.HealthStatusHandler)
	karlGroup.GET("/chat/", ChatApi.ChatHandler)
	karlGroup.GET("/jms_state/", HandlerApi.JmsStateHandler)
	karlGroup.GET("/interrupt_current_ask/", HandlerApi.InterruptCurrentAskHandler)
	return eng
}
