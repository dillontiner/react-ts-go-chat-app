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
	// TODO: use env vars
	// Deffered Improvement: handling different env configurations
	// DB_PORT := "64174" // localhost
	DB_PORT := "32076" // http://192.168.49.2:32076
	dsn := fmt.Sprintf("host=192.168.49.2 user=postgres password=TfwePfOzum dbname=chat_app port=%v sslmode=disable", DB_PORT)
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
	result := c.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
