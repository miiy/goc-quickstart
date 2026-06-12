package repository

import (
	"context"
	"time"

	"github.com/miiy/goc/redis"
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
	return &tokenRepository{
		rdb: rdb,
	}
}

// GetToken get token from redis
func (r *tokenRepository) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

// SetToken set token to redis
func (r *tokenRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

// DeleteToken delete token from redis
func (r *tokenRepository) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}
