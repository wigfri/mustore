package redisotp

import (
	"context"
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wigfri/mustore/domain/services"
)

type Store struct {
	rdb *redis.Client
}

func New(addr, password string, db int) *Store {
	return &Store{
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.rdb.Ping(ctx).Err()
}

var _ services.OTPStore = (*Store)(nil)

func keyEmailVerify(id uuid.UUID) string {
	return fmt.Sprintf("mustore:email_verify:%s", id.String())
}

func keyLoginOTP(id uuid.UUID) string {
	return fmt.Sprintf("mustore:login_otp:%s", id.String())
}

func (s *Store) SetEmailVerification(ctx context.Context, userID uuid.UUID, code string, ttl time.Duration) error {
	return s.rdb.Set(ctx, keyEmailVerify(userID), code, ttl).Err()
}

func (s *Store) VerifyEmailCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	key := keyEmailVerify(userID)
	stored, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if subtle.ConstantTimeCompare([]byte(stored), []byte(code)) != 1 {
		return false, nil
	}
	_ = s.rdb.Del(ctx, key).Err()
	return true, nil
}

func (s *Store) SetLoginOTP(ctx context.Context, userID uuid.UUID, code string, ttl time.Duration) error {
	return s.rdb.Set(ctx, keyLoginOTP(userID), code, ttl).Err()
}

func (s *Store) VerifyLoginOTP(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	key := keyLoginOTP(userID)
	stored, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if subtle.ConstantTimeCompare([]byte(stored), []byte(code)) != 1 {
		return false, nil
	}
	_ = s.rdb.Del(ctx, key).Err()
	return true, nil
}

func (s *Store) InvalidateUserOTP(ctx context.Context, userID uuid.UUID) error {
	return s.rdb.Del(ctx, keyEmailVerify(userID), keyLoginOTP(userID)).Err()
}
