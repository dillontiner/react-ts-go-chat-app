package persistence

import (
	"errors"
	"fmt"
	"os"
	"server/entities"

	uuid "github.com/satori/go.uuid"
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
	// TODO: caller or this should set uuid

	user.UUID = uuid.NewV4()
	result := c.db.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (c *Client) GetUserByEmail(email string) (*entities.User, error) {
	// TODO: caller or this should set uuid
	user := entities.User{}
	result := c.db.First(&user, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (c *Client) AuthorizeUser(email string, password string) (*uuid.UUID, error) {
	// TODO: caller or this should set uuid
	user, err := c.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("unauthorized")
	}

	return &user.UUID, nil
}

func (c *Client) GetMessages() (*[]entities.Message, error) {
	// TODO: caller or this should set uuid
	messages := []entities.Message{}

	// TODO: pagination
	result := c.db.Find(&messages, "LIMIT 100")
	if result.Error != nil {
		return nil, result.Error
	}

	return &messages, nil
}
