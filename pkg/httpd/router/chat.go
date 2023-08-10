package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jumpserver/kael/pkg/httpd/ws"
	"github.com/jumpserver/kael/pkg/jms"
	"github.com/jumpserver/kael/pkg/logger"
	"github.com/jumpserver/kael/pkg/manager"
	"github.com/jumpserver/kael/pkg/schemas"
	"github.com/jumpserver/wisp/protobuf-go/protobuf"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var ChatApi = new(_ChatApi)

type _ChatApi struct{}

func (s *_ChatApi) ChatHandler(ctx *gin.Context) {
	conn, err := ws.UpgradeWsConn(ctx)
	if err != nil {
		logger.GlobalLogger.Error("Websocket upgrade err", zap.Error(err))
		return
	}

	token, ok := ctx.GetQuery("token")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token"})
		return
	}

	currentJMSS := make([]*jms.JMSSession, 0)
	tokenHandler := jms.NewTokenHandler()
	sessionHandler := jms.NewSessionHandler(conn)
	authInfo, err := tokenHandler.GetTokenAuthInfo(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "auth fail"})
		return
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.GlobalLogger.Info("Accept message error or connect closed")
			if len(currentJMSS) != 0 {
				for _, jmss := range currentJMSS {
					reason := "Websocket已关闭, 会话中断"
					jmss.Close(reason)
				}
			}
			return
		}

		if string(msg) == "ping" {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}

		var askRequest schemas.AskRequest
		_ = json.Unmarshal(msg, &askRequest)
		jmss := &jms.JMSSession{}
		if askRequest.ConversationID == "" {
			jmss = sessionHandler.CreateNewSession(authInfo)
			jmss.ActiveSession()
			currentJMSS = append(currentJMSS, jmss)
		} else {
			conversationID := askRequest.ConversationID
			jmss = jms.GlobalSessionManager.GetJMSSession(conversationID)
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
		doneCh := make(chan string)
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
		messageID := uuid.New()
		for {
			select {
			case answer := <-answerCh:
				response := schemas.AskResponse{
					Type:           schemas.Message,
					ConversationID: jmss.Session.Id,
					Message: &schemas.ChatGPTMessage{
						Content:    answer,
						ID:         messageID,
						CreateTime: time.Now(),
						Type:       schemas.Message,
						Role:       openai.ChatMessageRoleAssistant,
					},
				}
				jsonResponse, _ := json.Marshal(response)
				_ = jmss.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)
			case answer := <-doneCh:
				response := schemas.AskResponse{
					Type:           schemas.Message,
					ConversationID: jmss.Session.Id,
					Message: &schemas.ChatGPTMessage{
						Content:    answer,
						ID:         messageID,
						CreateTime: time.Now(),
						Type:       schemas.Finish,
						Role:       openai.ChatMessageRoleAssistant,
					},
				}
				jsonResponse, _ := json.Marshal(response)
				_ = jmss.Websocket.WriteMessage(websocket.TextMessage, jsonResponse)
				close(doneCh)
				close(answerCh)
				return answer
			}
		}
	}
}
