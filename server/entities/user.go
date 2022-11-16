package entities

import (
	uuid "github.com/satori/go.uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	UUID     uuid.UUID
	Name     string
	Email    string
	Password string
}
