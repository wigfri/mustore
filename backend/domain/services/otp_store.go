package services

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTPStore interface {
	SetEmailVerification(ctx context.Context, userID uuid.UUID, code string, ttl time.Duration) error
	VerifyEmailCode(ctx context.Context, userID uuid.UUID, code string) (bool, error)
	SetLoginOTP(ctx context.Context, userID uuid.UUID, code string, ttl time.Duration) error
	VerifyLoginOTP(ctx context.Context, userID uuid.UUID, code string) (bool, error)
	InvalidateUserOTP(ctx context.Context, userID uuid.UUID) error
}
