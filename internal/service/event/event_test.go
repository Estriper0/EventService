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
)

func TestEventService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

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
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

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
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Minute}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

	ctx := context.Background()
	event := &models.EventResponse{Id: 1, Title: "Event"}

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
					GetEvent(ctx, 1).
					Return(event, nil)
			},
			want:    event,
			wantErr: nil,
		},
		{
			name: "cache miss, repo success, cache set success",
			id:   2,
			setup: func() {
				mockCache.EXPECT().
					GetEvent(ctx, 2).
					Return(nil, assert.AnError)

				mockRepo.EXPECT().
					GetById(ctx, 2).
					Return(&models.EventResponse{Id: 2, Title: "DB Event"}, nil)

				mockCache.EXPECT().
					SetEvent(ctx, gomock.Any(), cfg.Redis.CacheTTL).
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
					GetEvent(ctx, 3).
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
					GetEvent(ctx, 4).
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
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

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
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

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
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

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
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

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

func TestEventService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		eventID int
		setup   func()
		wantErr error
	}{
		{
			name:    "success",
			userID:  "user1",
			eventID: 1,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user1", 1).
					Return(false, nil)
				mockRepo.EXPECT().
					IncreaseCurrentAttedance(ctx, 1).
					Return(nil)
				mockEURepo.EXPECT().
					Create(ctx, "user1", 1).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "already registered",
			userID:  "user2",
			eventID: 2,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user2", 2).
					Return(true, nil)
			},
			wantErr: service.ErrRegistered,
		},
		{
			name:    "max registered",
			userID:  "user3",
			eventID: 3,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user3", 3).
					Return(false, nil)
				mockRepo.EXPECT().
					IncreaseCurrentAttedance(ctx, 3).
					Return(repositories.ErrMaxRegistered)
			},
			wantErr: service.ErrMaxRegistered,
		},
		{
			name:    "event not found on increase",
			userID:  "user4",
			eventID: 4,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user4", 4).
					Return(false, nil)
				mockRepo.EXPECT().
					IncreaseCurrentAttedance(ctx, 4).
					Return(repositories.ErrRecordNotFound)
			},
			wantErr: service.ErrRecordNotFound,
		},
		{
			name:    "repository error on exists",
			userID:  "user5",
			eventID: 5,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user5", 5).
					Return(false, assert.AnError)
			},
			wantErr: service.ErrRepositoryError,
		},
		{
			name:    "repository error on increase",
			userID:  "user6",
			eventID: 6,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user6", 6).
					Return(false, nil)
				mockRepo.EXPECT().
					IncreaseCurrentAttedance(ctx, 6).
					Return(assert.AnError)
			},
			wantErr: assert.AnError,
		},
		{
			name:    "repository error on create",
			userID:  "user7",
			eventID: 7,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user7", 7).
					Return(false, nil)
				mockRepo.EXPECT().
					IncreaseCurrentAttedance(ctx, 7).
					Return(nil)
				mockEURepo.EXPECT().
					Create(ctx, "user7", 7).
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

			err := eventService.Register(ctx, tt.userID, tt.eventID)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEventService_CancellRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		eventID int
		setup   func()
		wantErr error
	}{
		{
			name:    "success",
			userID:  "user1",
			eventID: 1,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user1", 1).
					Return(true, nil)
				mockRepo.EXPECT().
					DecreaseCurrentAttedance(ctx, 1).
					Return(nil)
				mockEURepo.EXPECT().
					Delete(ctx, "user1", 1).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "not registered",
			userID:  "user2",
			eventID: 2,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user2", 2).
					Return(false, nil)
			},
			wantErr: service.ErrNotRegistered,
		},
		{
			name:    "event not found on decrease",
			userID:  "user3",
			eventID: 3,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user3", 3).
					Return(true, nil)
				mockRepo.EXPECT().
					DecreaseCurrentAttedance(ctx, 3).
					Return(repositories.ErrRecordNotFound)
			},
			wantErr: service.ErrRecordNotFound,
		},
		{
			name:    "repository error on exists",
			userID:  "user4",
			eventID: 4,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user4", 4).
					Return(false, assert.AnError)
			},
			wantErr: service.ErrRepositoryError,
		},
		{
			name:    "repository error on decrease",
			userID:  "user5",
			eventID: 5,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user5", 5).
					Return(true, nil)
				mockRepo.EXPECT().
					DecreaseCurrentAttedance(ctx, 5).
					Return(assert.AnError)
			},
			wantErr: service.ErrRepositoryError,
		},
		{
			name:    "event not found on delete",
			userID:  "user6",
			eventID: 6,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user6", 6).
					Return(true, nil)
				mockRepo.EXPECT().
					DecreaseCurrentAttedance(ctx, 6).
					Return(nil)
				mockEURepo.EXPECT().
					Delete(ctx, "user6", 6).
					Return(repositories.ErrRecordNotFound)
			},
			wantErr: service.ErrRecordNotFound,
		},
		{
			name:    "repository error on delete",
			userID:  "user7",
			eventID: 7,
			setup: func() {
				mockEURepo.EXPECT().
					Exists(ctx, "user7", 7).
					Return(true, nil)
				mockRepo.EXPECT().
					DecreaseCurrentAttedance(ctx, 7).
					Return(nil)
				mockEURepo.EXPECT().
					Delete(ctx, "user7", 7).
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

			err := eventService.CancellRegister(ctx, tt.userID, tt.eventID)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEventService_GetAllByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		setup   func()
		want    []*models.EventResponse
		wantErr error
	}{
		{
			name:   "success",
			userID: "user1",
			setup: func() {
				mockRepo.EXPECT().
					GetAllByUser(ctx, "user1").
					Return([]*models.EventResponse{{Id: 1, Title: "Event"}}, nil)
			},
			want:    []*models.EventResponse{{Id: 1, Title: "Event"}},
			wantErr: nil,
		},
		{
			name:   "repository error",
			userID: "user2",
			setup: func() {
				mockRepo.EXPECT().
					GetAllByUser(ctx, "user2").
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

			got, err := eventService.GetAllByUser(ctx, tt.userID)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_GetAllUsersByEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocksRepo.NewMockIEventRepository(ctrl)
	mockEURepo := mocksRepo.NewMockIEventUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	logger := logger.GetLogger("test")
	cfg := &config.Config{Redis: config.Redis{CacheTTL: time.Hour}}

	eventService := New(mockRepo, mockEURepo, mockCache, logger, cfg)

	ctx := context.Background()

	tests := []struct {
		name    string
		eventID int
		setup   func()
		want    *[]string
		wantErr error
	}{
		{
			name:    "success",
			eventID: 1,
			setup: func() {
				mockEURepo.EXPECT().
					GetAllByEvent(ctx, 1).
					Return(&[]string{"id_1", "id_2"}, nil)
			},
			want: &[]string{"id_1", "id_2"},
			wantErr: nil,
		},
		{
			name:    "repository error",
			eventID: 2,
			setup: func() {
				mockEURepo.EXPECT().
					GetAllByEvent(ctx, 2).
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

			got, err := eventService.GetAllUsersByEvent(ctx, tt.eventID)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
