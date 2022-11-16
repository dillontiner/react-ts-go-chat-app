package entities

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	UUID     uuid.UUID
	Name     string
	Email    string
	Password string
}
