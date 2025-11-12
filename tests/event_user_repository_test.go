// event_user_repository_test.go
package tests

import (
	"github.com/Estriper0/EventService/internal/models"
	"github.com/Estriper0/EventService/internal/repositories"
	"github.com/Estriper0/EventService/internal/repositories/event"
	eventuser "github.com/Estriper0/EventService/internal/repositories/event_user"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestEventUserRepository_Exists() {
	repo := eventuser.New(s.db)
	eventRepo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func() (string, int)
		userID  string
		eventID int
		want    bool
		wantErr error
	}{
		{
			name: "success - exists",
			setup: func() (string, int) {
				userID := "ea27ecf4-02b1-453d-965d-408253a874b9"
				eventID, err := eventRepo.Create(s.ctx, &models.EventCreateRequest{
					Title:   "Event",
					Creator: userID,
					Status:  models.StatusPublished,
				})
				require.NoError(s.T(), err)
				err = repo.Create(s.ctx, userID, eventID)
				require.NoError(s.T(), err)
				return userID, eventID
			},
			want:    true,
			wantErr: nil,
		},
		{
			name:    "success - does not exist",
			userID:  "ea27ecf4-02b1-453d-965d-408253a874b9",
			eventID: 999,
			want:    false,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setup != nil {
				tt.userID, tt.eventID = tt.setup()
			}

			got, err := repo.Exists(s.ctx, tt.userID, tt.eventID)

			require.ErrorIs(s.T(), err, tt.wantErr)
			require.Equal(s.T(), tt.want, got)
		})
	}
}

func (s *TestSuite) TestEventUserRepository_Create() {
	repo := eventuser.New(s.db)
	eventRepo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func() (string, int)
		userID  string
		eventID int
		wantErr error
	}{
		{
			name: "success",
			setup: func() (string, int) {
				userID := "ea27ecf4-02b1-453d-965d-408253a874b9"
				eventID, err := eventRepo.Create(s.ctx, &models.EventCreateRequest{
					Title:   "Event",
					Creator: userID,
					Status:  models.StatusPublished,
				})
				require.NoError(s.T(), err)
				return userID, eventID
			},
			wantErr: nil,
		},
		{
			name: "fail - already exists",
			setup: func() (string, int) {
				userID := "ea27ecf4-02b1-453d-965d-408253a874b9"
				eventID, err := eventRepo.Create(s.ctx, &models.EventCreateRequest{
					Title:   "Event",
					Creator: userID,
					Status:  models.StatusPublished,
				})
				require.NoError(s.T(), err)
				err = repo.Create(s.ctx, userID, eventID)
				require.NoError(s.T(), err)
				return userID, eventID
			},
			wantErr: repositories.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setup != nil {
				tt.userID, tt.eventID = tt.setup()
			}

			err := repo.Create(s.ctx, tt.userID, tt.eventID)

			require.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *TestSuite) TestEventUserRepository_Delete() {
	repo := eventuser.New(s.db)
	eventRepo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func() (string, int)
		userID  string
		eventID int
		wantErr error
	}{
		{
			name: "success",
			setup: func() (string, int) {
				userID := "ea27ecf4-02b1-453d-965d-408253a874b9"
				eventID, err := eventRepo.Create(s.ctx, &models.EventCreateRequest{
					Title:   "Event",
					Creator: userID,
					Status:  models.StatusPublished,
				})
				require.NoError(s.T(), err)
				err = repo.Create(s.ctx, userID, eventID)
				require.NoError(s.T(), err)
				return userID, eventID
			},
			wantErr: nil,
		},
		{
			name:    "fail - not found",
			userID:  "ea27ecf4-02b1-453d-965d-408253a874b9",
			eventID: 999,
			wantErr: repositories.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setup != nil {
				tt.userID, tt.eventID = tt.setup()
			}

			err := repo.Delete(s.ctx, tt.userID, tt.eventID)

			require.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *TestSuite) TestEventUserRepository_GetAllByEvent() {
	repo := eventuser.New(s.db)
	eventRepo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func() int
		eventID int
		wantLen int
		wantErr error
	}{
		{
			name: "success - multiple users",
			setup: func() int {
				user1 := "ea27ecf4-02b1-453d-965d-408253a874b9"
				user2 := "ea28ecf4-02b1-453d-965d-408253a874b9"
				eventID, err := eventRepo.Create(s.ctx, &models.EventCreateRequest{
					Title:   "Event",
					Creator: user1,
					Status:  models.StatusPublished,
				})
				require.NoError(s.T(), err)
				err = repo.Create(s.ctx, user1, eventID)
				require.NoError(s.T(), err)
				err = repo.Create(s.ctx, user2, eventID)
				require.NoError(s.T(), err)
				return eventID
			},
			wantLen: 2,
			wantErr: nil,
		},
		{
			name:    "success - no users",
			eventID: 999,
			wantLen: 0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setup != nil {
				tt.eventID = tt.setup()
			}

			got, err := repo.GetAllByEvent(s.ctx, tt.eventID)

			require.ErrorIs(s.T(), err, tt.wantErr)
			require.Len(s.T(), *got, tt.wantLen)
		})
	}
}
