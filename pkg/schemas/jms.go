package schemas

import "github.com/jumpserver/wisp/protobuf-go/protobuf"

type CommandRecord struct {
	Input     string             `json:"input,omitempty"`
	Output    string             `json:"output,omitempty"`
	RiskLevel protobuf.RiskLevel `json:"risk_level"`
}
type ReviewState int

const (
	Wait ReviewState =  -1 + iota
	Rejected
	Approve
)

type JMSState struct {
	ID             string      `json:"id"`
	ActivateReview ReviewState `json:"activate_review,omitempty"`
	NewDialogue    bool        `json:"new_dialogue,omitempty"`
}

type Conversation struct {
	ID string `json:"id"`
}
