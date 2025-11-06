package service

import (
	"context"

	"github.com/Estriper0/EventService/internal/models"
	pb "github.com/Estriper0/protobuf/gen/event"
)

type IEventService interface {
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
	Register(
		ctx context.Context,
		user_id string,
		event_id int,
	) error
	CancellRegister(
		ctx context.Context,
		user_id string,
		event_id int,
	) error
	GetAllByUser(
		ctx context.Context,
		user_id string,
	) (*pb.GetAllByUserResponse, error)
	GetAllUsersByEvent(
		ctx context.Context,
		event_id int,
	) (*pb.GetAllUsersByEventResponse, error)
}
