package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// init db driver.
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type store struct {
	db *sql.DB
}

func (s *store) Connect(ctx context.Context, connect string) error {
	db, err := sql.Open("pgx", connect)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	s.db = db
	return s.db.PingContext(ctx)
}

func (s *store) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *store) Create(ctx context.Context, event storage.Event) (int, error) {
	var query string
	var args []interface{}
	if event.Notification != nil {
		query = `
			INSERT INTO event (title, start, stop, description, user_id, notification)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING event_id
		`
		args = []interface{}{event.Title, event.Start, event.Stop, event.Description, event.UserID, event.Notification}
	} else {
		query = `
			INSERT INTO event (title, start, stop, description, user_id)
			VALUES($1, $2, $3, $4, $5)
			RETURNING event_id
		`
		args = []interface{}{event.Title, event.Start, event.Stop, event.Description, event.UserID}
	}
	var id int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("db exec: %w", err)
	}
	return id, nil
}

func (s *store) Update(ctx context.Context, id int, change storage.Event) error {
	var query string
	var args []interface{}
	if change.Notification != nil {
		query = `
			UPDATE event
			SET title = $1,
				start = $2,
				stop = $3,
				description = $4,
				notification = $5
			WHERE event_id = $6;
		`
		args = []interface{}{change.Title, change.Start, change.Stop, change.Description, change.Notification, id}
	} else {
		query = `
			UPDATE event
			SET title = $1,
				start = $2,
				stop = $3,
				description = $4,
				notification = null
			WHERE event_id = $5;
		`
		args = []interface{}{change.Title, change.Start, change.Stop, change.Description, id}
	}
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db rows affected: %w", err)
	}
	if count != 1 {
		return storage.ErrNotExistsEvent
	}
	return nil
}

func (s *store) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM event
		WHERE event_id = $1
	`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}

func (s *store) DeleteAll(ctx context.Context) error {
	query := `
		TRUNCATE TABLE event RESTART IDENTITY
	`
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}

func (s *store) ListAll(ctx context.Context) ([]storage.Event, error) {
	query := `
		SELECT event_id, title, start, stop, description, user_id, notification
		FROM event
		ORDER BY start
	`
	return s.queryList(ctx, query)
}

func (s *store) ListDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()
	query := `
		SELECT event_id, title, start, stop, description, user_id, notification
		FROM event
		WHERE extract(year from start) = $1 AND extract(month from start) = $2 AND extract(day from start) = $3
		ORDER BY start
	`
	return s.queryList(ctx, query, year, month, day)
}

func (s *store) ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	year, week := date.ISOWeek()
	query := `
		SELECT event_id, title, start, stop, description, user_id, notification
		FROM event
		WHERE extract(isoyear from start) = $1 AND extract(week from start) = $2
		ORDER BY start
	`
	return s.queryList(ctx, query, year, week)
}

func (s *store) ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	year, month, _ := date.Date()
	query := `
		SELECT event_id, title, start, stop, description, user_id, notification
		FROM event
		WHERE extract(year from start) = $1 AND extract(month from start) = $2
		ORDER BY start
	`
	return s.queryList(ctx, query, year, month)
}

func (s *store) queryList(ctx context.Context, query string, args ...interface{}) (result []storage.Event, resultErr error) {
	// проверка есть, чего линтер хочет непонятно
	//nolint:rowserrcheck
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil && resultErr == nil {
			resultErr = err
		}
	}()

	for rows.Next() {
		var event storage.Event
		var notification sql.NullInt64
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Start,
			&event.Stop,
			&event.Description,
			&event.UserID,
			&notification,
		)
		if err != nil {
			resultErr = fmt.Errorf("db scan: %w", err)
			return
		}
		if notification.Valid {
			event.Notification = (*time.Duration)(&notification.Int64)
		}
		result = append(result, event)
	}
	if err := rows.Err(); err != nil {
		resultErr = fmt.Errorf("db rows: %w", err)
		return
	}
	return
}

func (s *store) IsTimeBusy(ctx context.Context, start, stop time.Time, excludeID int) (bool, error) {
	var count int
	query := `
		SELECT Count(*) AS count
		FROM event
		WHERE start < $1 AND stop > $2 AND event_id != $3
	`
	err := s.db.QueryRowContext(ctx, query, stop, start, excludeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("db query: %w", err)
	}
	return count > 0, nil
}
