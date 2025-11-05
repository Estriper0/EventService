package event

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/Estriper0/EventService/internal/cache"
	"github.com/Estriper0/EventService/internal/config"
	"github.com/Estriper0/EventService/internal/models"
	"github.com/Estriper0/EventService/internal/repositories"
	"github.com/Estriper0/EventService/internal/service"
)

type EventService struct {
	eventRepo repositories.IEventRepository
	cache     cache.Cache
	logger    *slog.Logger
	config    *config.Config
}

func New(repo repositories.IEventRepository, cache cache.Cache, logger *slog.Logger, config *config.Config) *EventService {
	return &EventService{
		eventRepo: repo,
		cache:     cache,
		logger:    logger,
		config:    config,
	}
}

func (s *EventService) GetAll(ctx context.Context) ([]*models.EventResponse, error) {
	events, err := s.eventRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error(
			"Error getting all events",
			slog.String("err", err.Error()),
		)
		return nil, service.ErrRepositoryError
	}
	s.logger.Info(
		"Successful getting all events",
	)
	return events, nil
}

func (s *EventService) Create(ctx context.Context, event *models.EventCreateRequest) (int, error) {
	id, err := s.eventRepo.Create(ctx, event)
	if err != nil {
		s.logger.Error(
			"Error create event",
			slog.String("err", err.Error()),
		)
		return 0, service.ErrRepositoryError
	}
	s.logger.Info(
		"Successful create event",
		slog.Int("id", id),
	)
	return id, nil
}

func (s *EventService) GetById(ctx context.Context, id int) (*models.EventResponse, error) {
	event, err := s.cache.GetEvent(ctx, id)
	if err != nil && err != cache.ErrNotFound {
		s.logger.Error(
			"Error in redis getting event",
			slog.String("error", err.Error()),
		)
	} else if err == nil {
		return event, nil
	}

	event, err = s.eventRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			s.logger.Info(
				"Event not found",
				slog.Int("id", id),
			)
			return nil, service.ErrRecordNotFound
		}
		s.logger.Error(
			"Error getting event",
			slog.Int("id", id),
			slog.String("err", err.Error()),
		)
		return nil, service.ErrRepositoryError
	}
	s.logger.Info(
		"Successful getting event",
		slog.Int("id", id),
	)

	err = s.cache.SetEvent(ctx, event, s.config.Redis.CacheTTL)
	if err != nil {
		s.logger.Error(
			"Error in redis setting event",
			slog.String("error", err.Error()),
		)
	} else {
		s.logger.Info(
			"Successfully added to cache",
			slog.Int("id", id),
		)
	}

	return event, nil
}

func (s *EventService) DeleteById(ctx context.Context, id int) error {
	err := s.eventRepo.DeleteById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			s.logger.Warn(
				"Event not found",
				slog.Int("id", id),
			)
			return service.ErrRecordNotFound
		}
		s.logger.Error(
			"Error delete event",
			slog.Int("id", id),
			slog.String("err", err.Error()),
		)
		return service.ErrRepositoryError
	}
	err = s.cache.Del(ctx, "event:"+strconv.Itoa(id))
	if err != nil {
		s.logger.Error(
			"Error in redis delete event",
			slog.Int("id", id),
			slog.String("err", err.Error()),
		)
	}
	s.logger.Info(
		"Successful delete event",
		slog.Int("id", id),
	)
	return nil
}

func (s *EventService) Update(ctx context.Context, event *models.EventUpdateRequest) error {
	err := s.eventRepo.Update(ctx, event)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			s.logger.Warn(
				"Event not found",
				slog.Int("id", event.Id),
			)
			return service.ErrRecordNotFound
		}
		s.logger.Error(
			"Error update event",
			slog.Int("id", event.Id),
			slog.String("err", err.Error()),
		)
		return service.ErrRepositoryError
	}
	err = s.cache.Del(ctx, "event:"+strconv.Itoa(event.Id))
	if err != nil {
		s.logger.Error(
			"Error in redis delete event",
			slog.Int("id", event.Id),
			slog.String("err", err.Error()),
		)
	}
	s.logger.Info(
		"Successful update event",
		slog.Int("id", event.Id),
	)
	return nil
}

func (s *EventService) GetAllByCreator(ctx context.Context, creator string) ([]*models.EventResponse, error) {
	events, err := s.eventRepo.GetAllByCreator(ctx, creator)
	if err != nil {
		s.logger.Error(
			"Error getting all events by creator",
			slog.String("creator", creator),
			slog.String("err", err.Error()),
		)
		return nil, service.ErrRepositoryError
	}
	s.logger.Info(
		"Successful getting all events by creator",
		slog.String("creator", creator),
	)
	return events, nil
}

func (s *EventService) GetAllByStatus(ctx context.Context, status string) ([]*models.EventResponse, error) {
	events, err := s.eventRepo.GetAllByStatus(ctx, status)
	if err != nil {
		s.logger.Error(
			"Error getting all events by status",
			slog.String("status", status),
			slog.String("err", err.Error()),
		)
		return nil, service.ErrRepositoryError
	}
	s.logger.Info(
		"Successful getting all events",
		slog.String("status", status),
	)
	return events, nil
}
