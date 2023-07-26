package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jumpserver/kael/pkg/manager"
	"github.com/sashabaranov/go-openai"
	"net/http"
)

var ChatApi = new(_ChatApi)

type _ChatApi struct{}

func (s *_ChatApi) ChatHandler(ctx *gin.Context) {
	status := make(map[string]interface{})

	answerCh := make(chan string)
	doneCh := make(chan bool)
	authToken := ""
	baseURL := ""
	proxy := ""
	model := ""
	content := "你好"
	c := manager.NewClient(authToken, baseURL, proxy)
	askChatGPT := &manager.AskChatGPT{
		Client:         c,
		Model:          model,
		Content:        content,
		HistoryContent: make([]openai.ChatCompletionMessage, 0),
		AnswerCh:       answerCh,
		DoneCh:         doneCh,
	}
	go manager.ChatGPT(askChatGPT)

	fmt.Printf("回答内容：")
	for {
		select {
		case answer := <-answerCh:
			fmt.Printf(answer)
		case <-doneCh:
			fmt.Println("\n循环结束")
			return
		}
	}

	ctx.JSON(http.StatusOK, status)
}
