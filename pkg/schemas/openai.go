package schemas

import (
	"github.com/google/uuid"
	"time"
)

//class MessageType(StrEnum):
//	message = auto()
//	finish = auto()

type ChatGPTMessage struct {
	Content    string          `json:"content"`
	ID         uuid.UUID       `json:"id"`
	CreateTime time.Time       `json:"create_time,omitempty"`
	Type       AskResponseType `json:"type"`
	Role       string          `json:"role"`
}
