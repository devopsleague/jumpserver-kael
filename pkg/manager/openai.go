package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/jumpserver/kael/pkg/utils"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/url"
)

func NewClient(authToken, baseURL, proxy string) *openai.Client {
	config := openai.DefaultConfig(authToken)
	config.BaseURL = baseURL
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

func ChatGPT(ask *AskChatGPT) {
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
		Model:     ask.Model,
		MaxTokens: utils.GetMaxInt(),
		Messages:  messages,
		Stream:    true,
	}

	stream, err := ask.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		close(ask.DoneCh)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			close(ask.DoneCh)
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			close(ask.DoneCh)
			return
		}

		ask.AnswerCh <- response.Choices[0].Delta.Content
	}
}
