package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/AayushVinayDev/URL-micro/shortener"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.client.HMSet(ctx, key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
