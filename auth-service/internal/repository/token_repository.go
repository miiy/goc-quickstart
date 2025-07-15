package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type tokenRepository struct {
	rdb redis.UniversalClient
}

func NewTokenRepository(rdb redis.UniversalClient) TokenRepository {
	return &authRepository{
		rdb: rdb,
	}
}

// GetToken get token from redis
func (r *authRepository) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

// SetToken set token to redis
func (r *authRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

// DeleteToken delete token from redis
func (r *authRepository) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}
