package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	UUID              uuid.UUID
	SentAt            time.Time `json:"sentAt"`
	SenderUUID        uuid.UUID `json:"senderUuid"`
	Body              string    `json:"body"`
	UpvoteUserUUIDS   []uuid.UUID
	DownvoteUserUUIDS []uuid.UUID
}
