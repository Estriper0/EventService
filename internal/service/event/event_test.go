package event

import (
	"context"
	"testing"
	"time"

	"github.com/Estriper0/EventService/internal/cache/mocks"
	"github.com/Estriper0/EventService/internal/config"
	"github.com/Estriper0/EventService/internal/logger"
	"github.com/Estriper0/EventService/internal/models"
	"github.com/Estriper0/EventService/internal/repositories"
	mocksRepo "github.com/Estriper0/EventService/internal/repositories/mocks"
	"github.com/Estriper0/EventService/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestEventService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func()
		want    []*models.EventResponse
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				mockRepo.EXPECT().
					GetAll(ctx).
					Return([]*models.EventResponse{
						{Id: 1, Title: "Event 1"},
					}, nil)
			},
			want:    []*models.EventResponse{{Id: 1, Title: "Event 1"}},
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				mockRepo.EXPECT().
					GetAll(ctx).
					Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: service.ErrRepositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := eventService.GetAll(ctx)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockCache, logger, cfg)

	ctx := context.Background()
	req := &models.EventCreateRequest{Title: "New Event"}

	tests := []struct {
		name    string
		setup   func()
		wantID  int
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				mockRepo.EXPECT().
					Create(ctx, req).
					Return(42, nil)
			},
			wantID:  42,
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				mockRepo.EXPECT().
					Create(ctx, req).
					Return(0, assert.AnError)
			},
			wantID:  0,
			wantErr: service.ErrRepositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			gotID, err := eventService.Create(ctx, req)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantID, gotID)
		})
	}
}

func TestEventService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Minute}}

	eventService := New(mockRepo, mockCache, logger, cfg)

	ctx := context.Background()
	event := &models.EventResponse{Id: 1, Title: "Event"}
	eventBytes, _ := msgpack.Marshal(event)

	tests := []struct {
		name    string
		id      int
		setup   func()
		want    *models.EventResponse
		wantErr error
	}{
		{
			name: "cache hit",
			id:   1,
			setup: func() {
				mockCache.EXPECT().
					GetBytes(ctx, "event:1").
					Return(eventBytes, nil)
			},
			want:    event,
			wantErr: nil,
		},
		{
			name: "cache miss, repo success, cache set success",
			id:   2,
			setup: func() {
				mockCache.EXPECT().
					GetBytes(ctx, "event:2").
					Return(nil, assert.AnError)

				mockRepo.EXPECT().
					GetById(ctx, 2).
					Return(&models.EventResponse{Id: 2, Title: "DB Event"}, nil)

				mockCache.EXPECT().
					Set(ctx, "event:2", gomock.Any(), cfg.Redis.CacheTTL).
					Return(nil)
			},
			want:    &models.EventResponse{Id: 2, Title: "DB Event"},
			wantErr: nil,
		},
		{
			name: "cache miss, repo not found",
			id:   3,
			setup: func() {
				mockCache.EXPECT().
					GetBytes(ctx, "event:3").
					Return(nil, assert.AnError)

				mockRepo.EXPECT().
					GetById(ctx, 3).
					Return(nil, repositories.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: service.ErrRecordNotFound,
		},
		{
			name: "cache miss, repo error",
			id:   4,
			setup: func() {
				mockCache.EXPECT().
					GetBytes(ctx, "event:4").
					Return(nil, assert.AnError)

				mockRepo.EXPECT().
					GetById(ctx, 4).
					Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: service.ErrRepositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := eventService.GetById(ctx, tt.id)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		id      int
		setup   func()
		wantErr error
	}{
		{
			name: "success",
			id:   1,
			setup: func() {
				mockRepo.EXPECT().
					DeleteById(ctx, 1).
					Return(nil)
				mockCache.EXPECT().Del(ctx, gomock.Any()).Return(assert.AnError)
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   2,
			setup: func() {
				mockRepo.EXPECT().
					DeleteById(ctx, 2).
					Return(repositories.ErrRecordNotFound)
			},
			wantErr: service.ErrRecordNotFound,
		},
		{
			name: "repository error",
			id:   3,
			setup: func() {
				mockRepo.EXPECT().
					DeleteById(ctx, 3).
					Return(assert.AnError)
			},
			wantErr: service.ErrRepositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := eventService.DeleteById(ctx, tt.id)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEventService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockCache, logger, cfg)

	ctx := context.Background()
	req := &models.EventUpdateRequest{Id: 1, Title: "Updated"}

	tests := []struct {
		name    string
		req     *models.EventUpdateRequest
		setup   func()
		wantErr error
	}{
		{
			name: "success",
			req:  req,
			setup: func() {
				mockRepo.EXPECT().
					Update(ctx, req).
					Return(nil)
				mockCache.EXPECT().Del(ctx, gomock.Any()).Return(assert.AnError)
			},
			wantErr: nil,
		},
		{
			name: "not found",
			req:  &models.EventUpdateRequest{Id: 2},
			setup: func() {
				mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					Return(repositories.ErrRecordNotFound)
			},
			wantErr: service.ErrRecordNotFound,
		},
		{
			name: "repository error",
			req:  &models.EventUpdateRequest{Id: 3},
			setup: func() {
				mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					Return(assert.AnError)
			},
			wantErr: service.ErrRepositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := eventService.Update(ctx, tt.req)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEventService_GetAllByCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		creator string
		setup   func()
		want    []*models.EventResponse
		wantErr error
	}{
		{
			name:    "success",
			creator: "user1",
			setup: func() {
				mockRepo.EXPECT().
					GetAllByCreator(ctx, "user1").
					Return([]*models.EventResponse{{Id: 1, Creator: "user1"}}, nil)
			},
			want:    []*models.EventResponse{{Id: 1, Creator: "user1"}},
			wantErr: nil,
		},
		{
			name:    "repository error",
			creator: "user2",
			setup: func() {
				mockRepo.EXPECT().
					GetAllByCreator(ctx, "user2").
					Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: service.ErrRepositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := eventService.GetAllByCreator(ctx, tt.creator)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_GetAllByStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		status  string
		setup   func()
		want    []*models.EventResponse
		wantErr error
	}{
		{
			name:   "success",
			status: "active",
			setup: func() {
				mockRepo.EXPECT().
					GetAllByStatus(ctx, "active").
					Return([]*models.EventResponse{{Id: 1, Status: "active"}}, nil)
			},
			want:    []*models.EventResponse{{Id: 1, Status: "active"}},
			wantErr: nil,
		},
		{
			name:   "repository error",
			status: "inactive",
			setup: func() {
				mockRepo.EXPECT().
					GetAllByStatus(ctx, "inactive").
					Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: service.ErrRepositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := eventService.GetAllByStatus(ctx, tt.status)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
