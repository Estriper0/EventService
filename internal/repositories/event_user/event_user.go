package eventuser

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Estriper0/EventService/internal/repositories"
)

type EventUserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *EventUserRepository {
	return &EventUserRepository{
		db: db,
	}
}

func (r *EventUserRepository) Exists(ctx context.Context, user_id string, event_id int) (bool, error) {
	query := "SELECT user_id, event_id FROM event.event_user WHERE user_id = $1 AND event_id = $2"
	var ui string
	var ei int
	err := r.db.QueryRowContext(ctx, query, user_id, event_id).Scan(&ui, &ei)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *EventUserRepository) Create(ctx context.Context, user_id string, event_id int) error {
	query := "INSERT INTO event.event_user (user_id, event_id) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, user_id, event_id)
	if err != nil {
		return repositories.ErrAlreadyExists
	}
	return nil
}

func (r *EventUserRepository) Delete(ctx context.Context, user_id string, event_id int) error {
	query := "DELETE FROM event.event_user WHERE user_id = $1 AND event_id = $2"
	res, err := r.db.ExecContext(ctx, query, user_id, event_id)
	if err != nil {
		return err
	}
	i, _ := res.RowsAffected()
	if i == 0 {
		return repositories.ErrRecordNotFound
	}
	return nil
}

func (r *EventUserRepository) GetAllByEvent(ctx context.Context, event_id int) (*[]string, error) {
	query := "SELECT user_id FROM event.event_user WHERE event_id = $1"
	rows, err := r.db.QueryContext(ctx, query, event_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		users = append(users, id)
	}
	return &users, nil
}
