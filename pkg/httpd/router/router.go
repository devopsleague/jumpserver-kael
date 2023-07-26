package router

import "github.com/gin-gonic/gin"

func CreateRouter() *gin.Engine {
	eng := gin.Default()
	karlGroup := eng.Group("/kael")
	karlGroup.GET("/health/", HealthApi.HealthStatusHandler)
	karlGroup.GET("/chat/", ChatApi.ChatHandler)
	return eng
}
