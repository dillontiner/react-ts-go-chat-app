package entities

import (
	uuid "github.com/satori/go.uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	UUID uuid.UUID `json:"uuid"`
}
type User struct {
	UUID     uuid.UUID
	Email    string
	Password string
}
