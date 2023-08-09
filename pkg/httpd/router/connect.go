package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ConnectApi = new(_ConnectApi)

type _ConnectApi struct{}

func (s *_ConnectApi) ConnectHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	redirectURL := fmt.Sprintf("/kael/static/?token=%s", token)
	ctx.Redirect(http.StatusFound, redirectURL)
}
