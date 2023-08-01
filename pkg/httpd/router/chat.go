package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/global"
	"github.com/jumpserver/kael/pkg/httpd"
	"github.com/jumpserver/kael/pkg/jms"
	"github.com/jumpserver/kael/pkg/manager"
	"github.com/jumpserver/kael/pkg/schemas"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"time"
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
				response := schemas.AskResponse{
					Type:           schemas.Error,
					ConversationID: askRequest.ConversationID,
					SystemMessage:  "current session not found",
				}
				jsonResponse, _ := json.Marshal(response)
				_ = jmss.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)
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
		go manager.ChatGPT(askChatGPT, jmss)
		lastContent := ""
		messageID := uuid.New()
		for {
			select {
			case answer := <-answerCh:
				// TODO 包装一下 websocket send 方法
				response := schemas.AskResponse{
					Type:           schemas.Message,
					ConversationID: askRequest.ConversationID,
					Message: &schemas.ChatGPTMessage{
						Content:    answer,
						ID:         uuid.New(),
						Parent:     messageID,
						CreateTime: time.Now(),
						Type:       schemas.Message,
						Role:       openai.ChatMessageRoleAssistant,
					},
				}
				jsonResponse, _ := json.Marshal(response)
				_ = jmss.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)
				lastContent = answer
			case <-doneCh:
				response := schemas.AskResponse{
					Type:           schemas.Message,
					ConversationID: askRequest.ConversationID,
					Message: &schemas.ChatGPTMessage{
						Content:    lastContent,
						ID:         uuid.New(),
						Parent:     messageID,
						CreateTime: time.Now(),
						Type:       schemas.Finish,
						Role:       openai.ChatMessageRoleAssistant,
					},
				}
				jsonResponse, _ := json.Marshal(response)
				_ = jmss.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)
				close(answerCh)
				return lastContent
			}
		}
	}
}
