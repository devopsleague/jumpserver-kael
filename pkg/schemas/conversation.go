package schemas

type AskRequest struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Content        string `json:"content"`
}

type AskResponseType string

const (
	Waiting AskResponseType = "waiting"
	Reject  AskResponseType = "reject"
	Message AskResponseType = "message"
	Error   AskResponseType = "error"
	Finish  AskResponseType = "finish"
)

type ResponseMeta struct {
	ActivateReview bool `json:"activate_review" default:"false"`
}

type AskResponse struct {
	Type           AskResponseType `json:"type"`
	ConversationID string          `json:"conversation_id,omitempty"`
	Message        *ChatGPTMessage `json:"message,omitempty"`
	SystemMessage  string          `json:"system_message,omitempty"`
	Meta           ResponseMeta    `json:"meta,omitempty"`
}
