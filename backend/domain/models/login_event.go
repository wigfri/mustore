package models

import (
	"time"

	"github.com/google/uuid"
)

type LoginEvent struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

