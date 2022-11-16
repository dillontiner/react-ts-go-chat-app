package persistence

import (
	"fmt"
	"server/entities"

	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient() (*Client, error) {
	// TODO password to env var
	dsn := "host=localhost user=postgres password=TfwePfOzum dbname=chat_app port=5432 sslmode=disable"
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
	fmt.Println(user.UUID)

	user.UUID = uuid.NewV4()
	result := c.db.Create(&user) // pass pointer of data to Create
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
