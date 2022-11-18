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
	UpvoteUserUUIDS   []uuid.UUID `json:"upvoteUserUuids"`
	DownvoteUserUUIDS []uuid.UUID `json:"downvoteUserUuids"`
}

type MessageRecord struct {
	UUID       uuid.UUID `json:"uuid"`
	SentAt     time.Time `json:"sentAt"`
	SenderUUID uuid.UUID `json:"senderUuid"`
	Body       string    `json:"body"`
}

func (MessageRecord) TableName() string {
	return "messages"
}

type Vote struct {
	UUID        uuid.UUID `json:"uuid"`
	MessageUUID uuid.UUID `json:"messageUuid"`
	VoterUUID   uuid.UUID `json:"voterUuid"`
	Vote        bool      `json:"vote"`
}

type GetChatResponse struct {
	Messages []Message `json:"messages"`
}
