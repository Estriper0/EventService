package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/Estriper0/EventService/internal/cache"
	"github.com/Estriper0/EventService/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type redisCache struct {
	client *redis.Client
}

func New(client *redis.Client) *redisCache {
	return &redisCache{
		client: client,
	}
}

func (r *redisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) GetEvent(ctx context.Context, id int) (*models.EventResponse, error) {
	data, err := r.client.Get(ctx, "event:"+strconv.Itoa(id)).Bytes()

	if err != redis.Nil && err != nil {
		return nil, err
	}

	if err == nil {
		event := &models.EventResponse{}
		err = msgpack.Unmarshal(data, event)
		if err == nil {
			return event, nil
		} else {
			return nil, err
		}
	}

	return nil, cache.ErrNotFound
}

func (r *redisCache) SetEvent(ctx context.Context, event *models.EventResponse, ttl time.Duration) error {
	data, err := msgpack.Marshal(event)
	if err == nil {
		err = r.client.Set(ctx, "event:"+strconv.Itoa(event.Id), data, ttl).Err()
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return err
	}
}
