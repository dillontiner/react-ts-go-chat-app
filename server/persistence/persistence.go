package persistence

import (
	"errors"
	"fmt"
	"os"
	"server/entities"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient() (*Client, error) {
	DB_HOST := os.Getenv("DB_HOST")
	if DB_HOST == "" {
		return nil, errors.New("missing DB_HOST env var")
	}

	DB_PORT := os.Getenv("DB_PORT")
	if DB_PORT == "" {
		return nil, errors.New("missing DB_HOST env var")
	}

	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	if DB_PASSWORD == "" {
		return nil, errors.New("missing DB_HOST env var")
	}

	dsn := fmt.Sprintf("host=%v user=postgres password=%v dbname=chat_app port=%v sslmode=disable", DB_HOST, DB_PASSWORD, DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	client := Client{
		db: db,
	}
	return &client, nil
}

func (c *Client) CreateUser(user entities.User) (*entities.User, error) {

	user.UUID = uuid.NewV4()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user.Password = string(hashedPassword)
	result := c.db.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (c *Client) GetUserByEmail(email string) (*entities.User, error) {
	user := entities.User{}
	result := c.db.First(&user, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (c *Client) AuthorizeUser(email string, password string) (*uuid.UUID, error) {
	user, err := c.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	return &user.UUID, nil
}

type MessageVoteResult struct {
}

func (c *Client) GetMessages() (*[]entities.Message, error) {
	fmt.Println("GETTING MESSAGES")
	messageRecords := []entities.MessageRecord{}

	// TODO: paginated loading when user scrolls
	// latest N messages
	result := c.db.Order("sent_at desc").Limit(30).Find(&messageRecords)
	if result.Error != nil {
		return nil, result.Error
	}

	// original plan to use array of strings caused issues with gorm library, hacking the association for MVP
	messages := []entities.Message{}
	messageByUUID := map[uuid.UUID]*entities.Message{}
	messageUUIDs := []uuid.UUID{}
	for _, m := range messageRecords {
		messageUUIDs = append(messageUUIDs, m.UUID)
		message := entities.Message{
			UUID:       m.UUID,
			SentAt:     m.SentAt,
			SenderUUID: m.SenderUUID,
			Body:       m.Body,
		}
		messageByUUID[m.UUID] = &message
		messages = append(messages, message)
	}

	votes := []entities.Vote{}
	result = c.db.Where("message_uuid IN ?", messageUUIDs).Find(&votes)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, v := range votes {
		if messageByUUID[v.MessageUUID] == nil {
			continue
		}

		// add the voter
		if v.Vote {
			upvotes := messageByUUID[v.MessageUUID].UpvoteUserUUIDS
			messageByUUID[v.MessageUUID].UpvoteUserUUIDS = append(upvotes, v.VoterUUID)
		} else {
			downvotes := messageByUUID[v.MessageUUID].DownvoteUserUUIDS
			messageByUUID[v.MessageUUID].DownvoteUserUUIDS = append(downvotes, v.VoterUUID)
		}

	}

	return &messages, nil
}

func (c *Client) CreateMessage(message entities.Message) (*entities.Message, error) {
	message.UUID = uuid.NewV4()
	messageRecord := entities.MessageRecord{
		UUID:       message.UUID,
		SentAt:     message.SentAt,
		SenderUUID: message.SenderUUID,
		Body:       message.Body,
	}

	result := c.db.Create(&messageRecord)

	if result.Error != nil {
		return nil, result.Error
	}

	return &message, nil
}

func (c *Client) VoteOnMessage(vote entities.Vote) (*entities.Message, error) {
	existingVote := entities.Vote{}
	result := c.db.Where("message_uuid = ?", vote.MessageUUID).Where("voter_uuid = ?", vote.VoterUUID).First(&existingVote)
	if result.Error == nil {
		var trueFalseNull string
		if existingVote.Vote == vote.Vote {
			trueFalseNull = "NULL" // effectively remove vote
		} else if vote.Vote {
			trueFalseNull = "True"
		} else {
			trueFalseNull = "False"
		}

		sql := fmt.Sprintf("UPDATE votes SET vote = %s WHERE message_uuid = '%s' AND voter_uuid = '%s'", trueFalseNull, vote.MessageUUID.String(), vote.VoterUUID.String())
		update := c.db.Exec(sql)
		if update.Error != nil {
			return nil, update.Error
		}
	} else {
		// record not found okay
		vote.UUID = uuid.NewV4()
		create := c.db.Create(vote)
		if create.Error != nil {
			return nil, create.Error
		}
	}

	// SELECT
	// MIN(m.uuid) uuid,
	// MIN(m.body) body,
	// MIN(m.sent_at) sent_at,
	// MIN(m.sender_uuid) sender_uuid,
	// STRING_AGG(v.voter_uuid, ',') voter_uuids,
	// v.vote
	// from messages m
	// join votes v
	// on m.uuid = v.message_uuid
	// WHERE v.vote IS NOT NULL
	// GROUP BY message_UUID, v.vote

	return nil, nil
}
