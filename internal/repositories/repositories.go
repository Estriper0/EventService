package repositories

import (
	"context"

	"github.com/Estriper0/EventService/internal/models"
)

type IEventRepository interface {
	GetById(
		ctx context.Context,
		id int,
	) (*models.EventResponse, error)
	Create(
		ctx context.Context,
		event *models.EventCreateRequest,
	) (int, error)
	GetAll(
		ctx context.Context,
	) ([]*models.EventResponse, error)
	GetAllByCreator(
		ctx context.Context,
		creator string,
	) ([]*models.EventResponse, error)
	GetAllByStatus(
		ctx context.Context,
		status string,
	) ([]*models.EventResponse, error)
	DeleteById(
		ctx context.Context,
		id int,
	) error
	Update(
		ctx context.Context,
		event *models.EventUpdateRequest,
	) error
	IncreaseCurrentAttedance(
		ctx context.Context,
		event_id int,
	) error
	DecreaseCurrentAttedance(
		ctx context.Context,
		event_id int,
	) error
	GetAllByUser(
		ctx context.Context,
		user_id string,
	) ([]*models.EventResponse, error)
}

type IEventUserRepository interface {
	Exists(
		ctx context.Context,
		user_id string,
		event_id int,
	) (bool, error)
	Create(
		ctx context.Context,
		user_id string,
		event_id int,
	) error
	Delete(
		ctx context.Context,
		user_id string,
		event_id int,
	) error
	GetAllByEvent(
		ctx context.Context,
		event_id int,
	) (*[]string, error)
}
