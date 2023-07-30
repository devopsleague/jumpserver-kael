package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jumpserver/kael/pkg/jms"
	"github.com/jumpserver/kael/pkg/manager"
	"github.com/jumpserver/kael/pkg/schemas"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"net/http"
)

var ChatApi = new(_ChatApi)

type _ChatApi struct{}

func (s *_ChatApi) ChatHandler(ctx *gin.Context) {
	status := make(map[string]interface{})

	ctx.JSON(http.StatusOK, status)
}

func chatFunc(authInfo *protobuf.TokenAuthInfo, askRequest schemas.AskRequest) func(jmss *jms.JMSSession) string {
	return func(jmss *jms.JMSSession) string {
		doneCh := make(chan bool)
		answerCh := make(chan string)
		model := authInfo.Platform.Protocols[0].Settings["api_mode"]
		jmss.HistoryAsks = append(jmss.HistoryAsks, askRequest.Content)
		c := manager.NewClient(
			authInfo.Account.Secret,
			authInfo.Asset.Address,
			authInfo.Asset.Specific.HttpProxy,
		)
		askChatGPT := &manager.AskChatGPT{
			Client:   c,
			Model:    model,
			Contents: jmss.HistoryAsks,
			AnswerCh: answerCh,
			DoneCh:   doneCh,
		}
		go manager.ChatGPT(askChatGPT)
		lastContent := ""
		for {
			select {
			case answer := <-answerCh:
				// websocket send
				lastContent = answer
				fmt.Printf(answer)
			case <-doneCh:
				close(answerCh)
				return lastContent
			}
		}
	}
}
