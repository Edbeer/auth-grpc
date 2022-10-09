package core

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	RefreshToken string    `json:"refresh_token" redis:"refresh_token"`
	Uuid         uuid.UUID `json:"user_id" redis:"user_id"`
	ExpireAt     time.Time `json:"expire_at" redis:"expire_at"`
}
