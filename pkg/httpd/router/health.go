package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var HealthApi = new(_HealthApi)

type _HealthApi struct{}

var upTime = time.Now()

func (s *_HealthApi) HealthStatusHandler(ctx *gin.Context) {
	status := make(map[string]interface{})
	now := time.Now()
	status["timestamp"] = now.UTC()
	status["uptime"] = now.Sub(upTime).String()
	ctx.JSON(http.StatusOK, status)
}
