package db

import (
	"context"
	"time"

	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/go-redis/redis/v8"
)

type REDIS struct {
	db *redis.Client
}

// defiine database connection for redis
var (
	Redis = REDIS{
		db: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: types.REDIS_PASSWORD,
			DB:       0,
		}),
	}
	ctx = context.Background()
)

func (r *REDIS) SetCookie(userId, cookie string, expires *time.Time) error {
	err := r.db.Set(ctx, userId, cookie, time.Until(*expires)).Err()
	return err
}

func (r *REDIS) GetCookie(userId string) (string, error) {
	val, err := r.db.Get(ctx, userId).Result()
	return val, err
}

func (r *REDIS) GetUserId(cookie string) (string, error) {
	keys, _, err := r.db.Scan(ctx, 0, "*", 0).Result()

	for _, key := range keys {
		val, err := r.db.Get(ctx, key).Result()
		if val == cookie {
			return key, err
		}
	}

	return "", err
}

func (r *REDIS) DeleteCookie(userId string) error {
	_, err := r.db.Del(ctx, userId).Result()
	return err
}
