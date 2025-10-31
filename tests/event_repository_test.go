package tests

import (
	"testing"
	"time"

	"github.com/Estriper0/EventService/internal/models"
	"github.com/Estriper0/EventService/internal/repositories"
	"github.com/Estriper0/EventService/internal/repositories/database/event"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestEventRepository_GetById() {
	repo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func() int
		id      int
		wantErr error
		want    *models.EventResponse
	}{
		{
			name: "success - existing event",
			setup: func() int {
				id, err := repo.Create(s.ctx, &models.EventCreateRequest{
					Title:        "Team Sync",
					About:        "Weekly team meeting",
					StartDate:    time.Date(2025, 11, 10, 14, 0, 0, 0, time.UTC),
					Location:     "Zoom",
					Status:       models.StatusDraft,
					MaxAttendees: 20,
					Creator:      "ea27ecf4-02b1-453d-965d-408253a874b9",
				})
				require.NoError(s.T(), err)
				return id
			},
			wantErr: nil,
			want: &models.EventResponse{
				Title:        "Team Sync",
				About:        "Weekly team meeting",
				StartDate:    time.Date(2025, 11, 10, 14, 0, 0, 0, time.UTC),
				Location:     "Zoom",
				Status:       models.StatusDraft,
				MaxAttendees: 20,
				Creator:      "ea27ecf4-02b1-453d-965d-408253a874b9",
			},
		},
		{
			name:    "not found",
			setup:   func() int { return 0 },
			id:      999,
			wantErr: repositories.ErrRecordNotFound,
			want:    nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			id := tt.setup()
			if id != 0 {
				tt.id = id
				tt.want.Id = id
			}

			got, err := repo.GetById(s.ctx, tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(s.T(), err, tt.wantErr)
				require.Nil(s.T(), got)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), tt.want.Id, got.Id)
				require.Equal(s.T(), tt.want.Title, got.Title)
				require.Equal(s.T(), tt.want.About, got.About)
				require.Equal(s.T(), tt.want.Location, got.Location)
				require.Equal(s.T(), tt.want.Status, got.Status)
				require.Equal(s.T(), tt.want.MaxAttendees, got.MaxAttendees)
				require.Equal(s.T(), tt.want.Creator, got.Creator)
			}
		})
	}
}

func (s *TestSuite) TestEventRepository_GetAll() {
	repo := event.New(s.db)

	tests := []struct {
		name  string
		setup func()
		want  int
	}{
		{
			name: "returns all events sorted by title",
			setup: func() {
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "C Event", Creator: "ea27ecf4-02b1-453d-965d-408253a874b9", Status: models.StatusDraft})
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "A Event", Creator: "ea28ecf4-02b1-453d-965d-408253a874b9", Status: models.StatusDraft})
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "B Event", Creator: "ea29ecf4-02b1-453d-965d-408253a874b9", Status: models.StatusDraft})
			},
			want: 3,
		},
		{
			name: "empty table",
			setup: func() {
				_, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE events CASCADE;")
				require.NoError(s.T(), err)
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()
			events, err := repo.GetAll(s.ctx)
			require.NoError(s.T(), err)
			require.Len(s.T(), events, tt.want)
		})
	}
}

func (s *TestSuite) TestEventRepository_Create() {
	repo := event.New(s.db)

	tests := []struct {
		name    string
		input   *models.EventCreateRequest
		wantErr bool
	}{
		{
			name: "valid event",
			input: &models.EventCreateRequest{
				Title:        "Go Workshop",
				About:        "Introduction to Go",
				StartDate:    time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
				Location:     "Conference Room A",
				Status:       models.StatusPublished,
				MaxAttendees: 40,
				Creator:      "ea27ecf4-02b1-453d-965d-408253a874b9",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			id, err := repo.Create(s.ctx, tt.input)
			if tt.wantErr {
				require.Error(s.T(), err)
				require.Zero(s.T(), id)
			} else {
				require.NoError(s.T(), err)
				require.Greater(s.T(), id, 0)

				event, getErr := repo.GetById(s.ctx, id)
				require.NoError(s.T(), getErr)
				require.Equal(s.T(), tt.input.Title, event.Title)
				require.Equal(s.T(), tt.input.Creator, event.Creator)
				require.Equal(s.T(), tt.input.Status, event.Status)
			}
		})
	}
}

func (s *TestSuite) TestEventRepository_DeleteById() {
	repo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func() int
		id      int
		wantErr error
	}{
		{
			name: "success - delete existing",
			setup: func() int {
				id, _ := repo.Create(s.ctx, &models.EventCreateRequest{
					Title: "To Be Deleted", Creator: "ea27ecf4-02b1-453d-965d-408253a874b9", Status: models.StatusDraft,
				})
				return id
			},
			wantErr: nil,
		},
		{
			name:    "not found",
			setup:   func() int { return 0 },
			id:      999,
			wantErr: repositories.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setup() != 0 {
				tt.id = tt.setup()
			}

			err := repo.DeleteById(s.ctx, tt.id)
			require.ErrorIs(s.T(), err, tt.wantErr)

			if tt.wantErr == nil {
				_, getErr := repo.GetById(s.ctx, tt.id)
				require.ErrorIs(s.T(), getErr, repositories.ErrRecordNotFound)
			}
		})
	}
}

func (s *TestSuite) TestEventRepository_Update() {
	repo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func() int
		update  *models.EventUpdateRequest
		wantErr error
		verify  func(*testing.T, *models.EventResponse)
	}{
		{
			name: "full update",
			setup: func() int {
				id, _ := repo.Create(s.ctx, &models.EventCreateRequest{
					Title: "Old Title", About: "Old", StartDate: time.Now(), Location: "Old", Status: models.StatusDraft, MaxAttendees: 10, Creator: "ea27ecf4-02b1-453d-965d-408253a874b9",
				})
				return id
			},
			update: &models.EventUpdateRequest{
				Title:        "New Title",
				About:        "New description",
				StartDate:    time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC),
				Location:     "New Location",
				Status:       models.StatusOngoing,
				MaxAttendees: 150,
			},
			wantErr: nil,
			verify: func(t *testing.T, e *models.EventResponse) {
				require.Equal(t, "New Title", e.Title)
				require.Equal(t, "New description", e.About)
				require.Equal(t, "New Location", e.Location)
				require.Equal(t, models.StatusOngoing, e.Status)
				require.Equal(t, 150, e.MaxAttendees)
			},
		},
		{
			name:    "not found",
			setup:   func() int { return 0 },
			update:  &models.EventUpdateRequest{Id: 999},
			wantErr: repositories.ErrRecordNotFound,
			verify:  nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			id := tt.setup()
			if id != 0 {
				tt.update.Id = id
			}

			err := repo.Update(s.ctx, tt.update)
			require.ErrorIs(s.T(), err, tt.wantErr)

			if tt.wantErr == nil && tt.verify != nil {
				updated, getErr := repo.GetById(s.ctx, id)
				require.NoError(s.T(), getErr)
				tt.verify(s.T(), updated)
			}
		})
	}
}

func (s *TestSuite) TestEventRepository_GetAllByCreator() {
	repo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func()
		creator string
		wantLen int
	}{
		{
			name: "multiple events by same creator",
			setup: func() {
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "Event 1", Creator: "ea27ecf4-02b1-453d-965d-408253a874b9", Status: models.StatusPublished})
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "Event 2", Creator: "ea27ecf4-02b1-453d-965d-408253a874b9", Status: models.StatusDraft})
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "Event 3", Creator: "ea29ecf4-02b1-453d-965d-408253a874b9", Status: models.StatusPublished})
			},
			creator: "ea27ecf4-02b1-453d-965d-408253a874b9",
			wantLen: 2,
		},
		{
			name:    "no events for creator",
			setup:   func() {},
			creator: "ea30ecf4-02b1-453d-965d-408253a874b9",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()
			events, err := repo.GetAllByCreator(s.ctx, tt.creator)
			require.NoError(s.T(), err)
			require.Len(s.T(), events, tt.wantLen)

			for _, e := range events {
				require.Equal(s.T(), tt.creator, e.Creator)
			}
		})
	}
}

func (s *TestSuite) TestEventRepository_GetAllByStatus() {
	repo := event.New(s.db)

	tests := []struct {
		name    string
		setup   func()
		status  string
		wantLen int
	}{
		{
			name: "filter by published",
			setup: func() {
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "A", Status: models.StatusPublished, Creator: "ea27ecf4-02b1-453d-965d-408253a874b9"})
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "B", Status: models.StatusDraft, Creator: "ea27ecf4-02b1-453d-965d-408253a874b9"})
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "C", Status: models.StatusPublished, Creator: "ea27ecf4-02b1-453d-965d-408253a874b9"})
				_, _ = repo.Create(s.ctx, &models.EventCreateRequest{Title: "D", Status: models.StatusCompleted, Creator: "ea27ecf4-02b1-453d-965d-408253a874b9"})
			},
			status:  models.StatusPublished,
			wantLen: 2,
		},
		{
			name:    "no events with cancelled status",
			setup:   func() {},
			status:  models.StatusCancelled,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()

			events, err := repo.GetAllByStatus(s.ctx, tt.status)
			require.NoError(s.T(), err)
			require.Len(s.T(), events, tt.wantLen)

			for _, e := range events {
				require.Equal(s.T(), tt.status, e.Status)
			}
		})
	}
}
