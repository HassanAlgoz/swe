package entities

import "github.com/google/uuid"

type UserProfile struct {
	ID       uuid.UUID
	Username string
}
