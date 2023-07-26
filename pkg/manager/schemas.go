package manager

import "github.com/sashabaranov/go-openai"

type AskChatGPT struct {
	Client         *openai.Client
	Model          string
	Content        string
	HistoryContent []openai.ChatCompletionMessage
	AnswerCh       chan string
	DoneCh         chan bool
}
