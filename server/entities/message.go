package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	UUID              uuid.UUID
	SentAt            time.Time   `json:"sentAt"`
	SenderUUID        uuid.UUID   `json:"senderUuid"`
	Body              string      `json:"body"`
	UpvoteUserUUIDS   []uuid.UUID `gorm:"type:text;column:upvote_user_uuids"`
	DownvoteUserUUIDS []uuid.UUID `gorm:"type:text;column:downvote_user_uuids"`
}
