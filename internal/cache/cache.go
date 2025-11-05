package cache

import (
	"context"
	"time"

	"github.com/Estriper0/EventService/internal/models"
)

type Cache interface {
	Del(ctx context.Context, key string) error
	GetEvent(ctx context.Context, id int) (*models.EventResponse, error)
	SetEvent(ctx context.Context, event *models.EventResponse, ttl time.Duration) error
}
