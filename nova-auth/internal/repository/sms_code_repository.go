package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/miiy/goc/redis"
)

const smsCodeKey = "sms_code:%s" // sms_code:{phone}

// SMSCodeRepository stores short-lived phone verification codes in Redis. It is
// deliberately separate from TokenRepository: a verification code and a refresh
// token are distinct domain concepts that just happen to share a Redis backend.
type SMSCodeRepository interface {
	Get(ctx context.Context, phone string) (string, error)
	Set(ctx context.Context, phone, code string, expiration time.Duration) error
	Del(ctx context.Context, phone string) error
}

type smsCodeRepository struct {
	rdb redis.UniversalClient
}

func NewSMSCodeRepository(rdb redis.UniversalClient) SMSCodeRepository {
	return &smsCodeRepository{rdb: rdb}
}

func (r *smsCodeRepository) Get(ctx context.Context, phone string) (string, error) {
	return r.rdb.Get(ctx, formatSMSCodeKey(phone)).Result()
}

func (r *smsCodeRepository) Set(ctx context.Context, phone, code string, expiration time.Duration) error {
	return r.rdb.Set(ctx, formatSMSCodeKey(phone), code, expiration).Err()
}

func (r *smsCodeRepository) Del(ctx context.Context, phone string) error {
	return r.rdb.Del(ctx, formatSMSCodeKey(phone)).Err()
}

func formatSMSCodeKey(phone string) string {
	return fmt.Sprintf(smsCodeKey, phone)
}
