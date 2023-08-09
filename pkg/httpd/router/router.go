package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jumpserver/kael/pkg/httpd/middlewares"
)

func CreateRouter() *gin.Engine {
	eng := gin.Default()
	eng.Use(middlewares.CORSMiddleware())
	karlGroup := eng.Group("/kael")
	karlGroup.Static("/static/", "ui/dist")
	karlGroup.Static("/assets/", "ui/dist/assets")

	karlGroup.GET("/chat/", ChatApi.ChatHandler)
	karlGroup.GET("/connect", ConnectApi.ConnectHandler)
	karlGroup.GET("/health/", HealthApi.HealthStatusHandler)
	karlGroup.POST("/jms_state/", HandlerApi.JmsStateHandler)
	karlGroup.POST("/interrupt_current_ask/", HandlerApi.InterruptCurrentAskHandler)
	return eng
}
