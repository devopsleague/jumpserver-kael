package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/jumpserver/kael/pkg/jms"
	"github.com/jumpserver/kael/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewClient(authToken, baseURL, proxy string) *openai.Client {
	config := openai.DefaultConfig(authToken)
	config.BaseURL = strings.TrimRight(baseURL, "/")
	if proxy != "" {
		AddProxy(&config, proxy)
	}
	return openai.NewClientWithConfig(config)
}

func AddProxy(config *openai.ClientConfig, proxy string) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		fmt.Println(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}
}

func ChatGPT(ask *AskChatGPT, jmss *jms.JMSSession) {
	// TODO 做超时处理
	ctx := context.Background()
	messages := make([]openai.ChatCompletionMessage, 0)

	for _, content := range ask.Contents {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	}

	req := openai.ChatCompletionRequest{
		Model:    ask.Model,
		Messages: messages,
		Stream:   true,
	}

	stream, err := ask.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		ask.DoneCh <- err.Error()
		return
	}
	defer stream.Close()
	content := ""
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) || jmss.CurrentAskInterrupt {
			jmss.CurrentAskInterrupt = false
			ask.DoneCh <- content
			return
		}

		if err != nil {
			logger.GlobalLogger.Error("openai stream error", zap.Error(err))
			ask.DoneCh <- content
			return
		}
		content += response.Choices[0].Delta.Content
		ask.AnswerCh <- content
	}
}
