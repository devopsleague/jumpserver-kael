package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/global"
	"github.com/jumpserver/kael/pkg/httpd"
	"github.com/jumpserver/kael/pkg/jms"
	"github.com/jumpserver/kael/pkg/manager"
	"github.com/jumpserver/kael/pkg/schemas"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"net/http"
)

var ChatApi = new(_ChatApi)

type _ChatApi struct{}

func (s *_ChatApi) ChatHandler(ctx *gin.Context) {
	conn, err := httpd.UpgradeWsConn(ctx)
	if err != nil {
		fmt.Println("Websocket upgrade err: ", err)
		return
	}
	token, ok := ctx.GetQuery("token")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token"})
		return
	}

	tokenHandler := jms.NewTokenHandler()
	sessionHandler := jms.NewSessionHandler(conn)
	authInfo, _ := tokenHandler.GetTokenAuthInfo(token)

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Accept message error: ", err)
			continue
		}

		var askRequest schemas.AskRequest
		err = json.Unmarshal(msg, &askRequest)
		if err != nil {
			fmt.Println("Invalid ask request: ", err)
			continue
		}
		jmss := &jms.JMSSession{}
		if askRequest.ConversationID == "" {
			jmss = sessionHandler.CreateNewSession(authInfo)
			jmss.ActiveSession()
		} else {
			conversationID := askRequest.ConversationID
			jmss = global.SessionManager.GetJMSSession(conversationID)
			if jmss == nil {
				fmt.Println("-----")
				continue
			} else {
				jmss.JMSState.NewDialogue = true
			}
		}
		go jmss.WithAudit(askRequest.Content, chatFunc(authInfo, askRequest))
	}
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
				response := schemas.AskResponse{
					Type:           schemas.Message,
					ConversationID: askRequest.ConversationID,
					Message: &schemas.ChatGPTMessage{
						Content: answer,
						// 其他ChatGPTMessage的字段
					},
				}
				jsonResponse, _ := json.Marshal(response)
				_ = jmss.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)
				lastContent = answer
			case <-doneCh:
				close(answerCh)
				return lastContent
			}
		}
	}
}
