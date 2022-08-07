package db

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type REDIS struct {
	db *redis.Client
}

var (
	Redis = REDIS{
		db: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: os.Getenv("DATABASE_PASSWORD"),
			DB:       0,
		}),
	}
	ctx = context.Background()
)

func (r *REDIS) SetCookie(userId, cookie string, expires *time.Time) error {
	err := r.db.Set(ctx, userId, cookie, expires.Sub(time.Now())).Err()
	return err
}

func (r *REDIS) GetCookie(userId string) (string, error) {
	val, err := r.db.Get(ctx, userId).Result()
	return val, err
}

func (r *REDIS) GetUserId(cookie string) (string, error) {
	var err error
	iter := r.db.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		val, err := r.db.Get(ctx, iter.Val()).Result()
		if val == cookie {
			return iter.Val(), err
		}
	}

	return "", err
}

func (r *REDIS) DeleteCookie(userId string) error {
	_, err := r.db.Del(ctx, userId).Result()
	return err
}
