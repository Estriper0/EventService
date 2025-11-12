package event

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Estriper0/EventService/internal/models"
	"github.com/Estriper0/EventService/internal/repositories"
)

type EventRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) GetById(
	ctx context.Context,
	id int,
) (*models.EventResponse, error) {
	query := "SELECT * FROM events WHERE id = $1"
	event := &models.EventResponse{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.Id,
		&event.Title,
		&event.About,
		&event.StartDate,
		&event.Location,
		&event.Status,
		&event.MaxAttendees,
		&event.CurrentAttendance,
		&event.Creator,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrRecordNotFound
		}
		return nil, err
	}
	return event, nil
}

func (r *EventRepository) GetAll(
	ctx context.Context,
) ([]*models.EventResponse, error) {
	query := "SELECT * FROM events ORDER BY title"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []*models.EventResponse{}

	for rows.Next() {
		event := &models.EventResponse{}
		err := rows.Scan(
			&event.Id,
			&event.Title,
			&event.About,
			&event.StartDate,
			&event.Location,
			&event.Status,
			&event.MaxAttendees,
			&event.CurrentAttendance,
			&event.Creator,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *EventRepository) Create(
	ctx context.Context,
	event *models.EventCreateRequest,
) (int, error) {
	var id int
	query := "INSERT INTO events (title, about, start_date, location, status, max_attendees, creator) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err := r.db.QueryRowContext(
		ctx,
		query,
		event.Title,
		event.About,
		event.StartDate,
		event.Location,
		event.Status,
		event.MaxAttendees,
		event.Creator,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *EventRepository) DeleteById(
	ctx context.Context,
	id int,
) error {
	query := "DELETE FROM events WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	i, _ := res.RowsAffected()
	if i == 0 {
		return repositories.ErrRecordNotFound
	}
	return nil
}

func (r *EventRepository) Update(
	ctx context.Context,
	event *models.EventUpdateRequest,
) error {
	query := "UPDATE events SET title = $1, about = $2, start_date = $3, location = $4, status = $5, max_attendees = $6 WHERE id = $7"
	res, err := r.db.ExecContext(
		ctx,
		query,
		event.Title,
		event.About,
		event.StartDate,
		event.Location,
		event.Status,
		event.MaxAttendees,
		event.Id,
	)

	if err != nil {
		return repositories.ErrRecordNotFound
	}

	i, _ := res.RowsAffected()
	if i == 0 {
		return repositories.ErrRecordNotFound
	}

	return nil
}

func (r *EventRepository) GetAllByCreator(
	ctx context.Context,
	creator string,
) ([]*models.EventResponse, error) {
	query := "SELECT * FROM events WHERE creator = $1 ORDER BY title"
	rows, err := r.db.QueryContext(ctx, query, creator)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []*models.EventResponse{}

	for rows.Next() {
		event := &models.EventResponse{}
		err := rows.Scan(
			&event.Id,
			&event.Title,
			&event.About,
			&event.StartDate,
			&event.Location,
			&event.Status,
			&event.MaxAttendees,
			&event.CurrentAttendance,
			&event.Creator,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *EventRepository) GetAllByStatus(
	ctx context.Context,
	status string,
) ([]*models.EventResponse, error) {
	query := "SELECT * FROM events WHERE status = $1 ORDER BY title"
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []*models.EventResponse{}

	for rows.Next() {
		event := &models.EventResponse{}
		err := rows.Scan(
			&event.Id,
			&event.Title,
			&event.About,
			&event.StartDate,
			&event.Location,
			&event.Status,
			&event.MaxAttendees,
			&event.CurrentAttendance,
			&event.Creator,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *EventRepository) IncreaseCurrentAttedance(ctx context.Context, event_id int) error {
	query := "UPDATE events SET current_attendance = current_attendance + 1 WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, event_id)

	if err != nil {
		return repositories.ErrMaxRegistered
	}

	i, _ := res.RowsAffected()
	if i == 0 {
		return repositories.ErrRecordNotFound
	}

	return nil
}

func (r *EventRepository) DecreaseCurrentAttedance(ctx context.Context, event_id int) error {
	query := "UPDATE events SET current_attendance = current_attendance - 1 WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, event_id)

	if err != nil {
		return err
	}

	i, _ := res.RowsAffected()
	if i == 0 {
		return repositories.ErrRecordNotFound
	}

	return nil
}

func (r *EventRepository) GetAllByUser(ctx context.Context, user_id string) ([]*models.EventResponse, error) {
	query := "SELECT events.* FROM events JOIN event_user ON events.id = event_user.event_id WHERE event_user.user_id = $1"
	rows, err := r.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []*models.EventResponse{}

	for rows.Next() {
		event := &models.EventResponse{}
		err := rows.Scan(
			&event.Id,
			&event.Title,
			&event.About,
			&event.StartDate,
			&event.Location,
			&event.Status,
			&event.MaxAttendees,
			&event.CurrentAttendance,
			&event.Creator,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
