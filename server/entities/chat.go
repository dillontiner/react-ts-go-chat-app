package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	UUID              uuid.UUID   `json:"uuid"`
	SentAt            time.Time   `json:"sentAt"`
	SenderUUID        uuid.UUID   `json:"senderUuid"`
	Body              string      `json:"body"`
	UpvoteUserUUIDS   []uuid.UUID `json:"upvoteUserUuids" gorm:"type:text;column:upvote_user_uuids"`
	DownvoteUserUUIDS []uuid.UUID `json:"downvoteUserUuids" gorm:"type:text;column:downvote_user_uuids"`
}

type GetChatResponse struct {
	Messages []Message `json:"messages"`
}
