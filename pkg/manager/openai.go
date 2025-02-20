package manager

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/jumpserver/kael/pkg/jms"
	"github.com/jumpserver/kael/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type TransportOptions struct {
	UseProxy        bool
	ProxyURL        *url.URL
	SkipCertificate bool
}

type TransportOption func(*TransportOptions)

func WithProxy(proxyURL string) TransportOption {
	UseProxy := proxyURL != ""
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		proxy = nil
		UseProxy = false
		logger.GlobalLogger.Error(err.Error(), zap.Error(err))
	}
	return func(opts *TransportOptions) {
		opts.UseProxy = UseProxy
		opts.ProxyURL = proxy
	}
}

func WithSkipCertificate(skip bool) TransportOption {
	return func(opts *TransportOptions) {
		opts.SkipCertificate = skip
	}
}

func NewCustomTransport(options ...TransportOption) *http.Transport {
	transportOpts := &TransportOptions{}

	for _, opt := range options {
		opt(transportOpts)
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: transportOpts.SkipCertificate}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	if transportOpts.UseProxy {
		transport.Proxy = http.ProxyURL(transportOpts.ProxyURL)
	}

	return transport
}

func NewClient(authToken, baseURL, proxy string) *openai.Client {
	config := openai.DefaultConfig(authToken)
	config.BaseURL = strings.TrimRight(baseURL, "/")
	transport := NewCustomTransport(
		WithProxy(proxy), WithSkipCertificate(true),
	)
	config.HTTPClient = &http.Client{
		Transport: transport,
	}
	return openai.NewClientWithConfig(config)
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
