package repository

import (
	"context"
	"time"

	"github.com/miiy/goc/redis"
)

type TokenRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// SetKeepTTL updates an existing key's value without changing its TTL.
	// Missing keys are treated as a no-op.
	SetKeepTTL(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, key string) error
	// CompareAndSet atomically sets key to newVal only if its current value equals
	// oldVal. Returns true if applied, false if the current value did not match.
	// The existing key TTL is preserved. Used for race-free refresh-token
	// rotation (reuse detection).
	CompareAndSet(ctx context.Context, key, oldVal, newVal string) (bool, error)
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

// SetKeepTTL updates an existing key without extending or removing its TTL.
func (r *tokenRepository) SetKeepTTL(ctx context.Context, key string, value interface{}) error {
	const script = `
local ttl = redis.call('PTTL', KEYS[1])
if ttl > 0 then
  redis.call('SET', KEYS[1], ARGV[1], 'PX', ttl)
  return 1
elseif ttl == -1 then
  redis.call('SET', KEYS[1], ARGV[1])
  return 1
else
  return 0
end`
	_, err := r.rdb.Eval(ctx, script, []string{key}, value).Int()
	return err
}

// DeleteToken delete token from redis
func (r *tokenRepository) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

// CompareAndSet atomically sets key to newVal only when its current value equals
// oldVal. Implemented as a Redis Lua script (GET + compare + PTTL + SET in one
// round-trip) so concurrent rotations of the same refresh token cannot both
// succeed, and the revoked record keeps the original expiration.
func (r *tokenRepository) CompareAndSet(ctx context.Context, key, oldVal, newVal string) (bool, error) {
	const script = `
if redis.call('GET', KEYS[1]) ~= ARGV[1] then
  return 0
end
local ttl = redis.call('PTTL', KEYS[1])
if ttl > 0 then
  redis.call('SET', KEYS[1], ARGV[2], 'PX', ttl)
elseif ttl == -1 then
  redis.call('SET', KEYS[1], ARGV[2])
else
  return 0
end
return 1`
	res, err := r.rdb.Eval(ctx, script, []string{key}, oldVal, newVal).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
