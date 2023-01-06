package data

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Session struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Version int    `json:"version"`
	db      *sql.DB
}

func (s *Session) GetByID(id string) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, version from sessions where id = $1`

	var session Session
	row := s.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&session.ID,
		&session.Email,
		&session.Version,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Session) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from sessions where id = $1`

	var session Session
	row := s.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&session.ID,
		&session.Email,
		&session.Version,
	)

	return err
}

func (s *Session) Insert(email string) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	id := uuid.New()
	version := 1

	var newID string
	stmt := `insert into sessions (email, id, version)
		values ($1, $2, $3) returning id`

	err := s.db.QueryRowContext(ctx, stmt,
		email,
		id.String(),
		version,
	).Scan(&newID)

	if err != nil {
		return "", version, err
	}

	return newID, version, nil
}

func (s *Session) UpdateVersion(ID string, version int) int {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	newVersion := version + 1

	stmt := `update sessions set version = $1 where id = $2`

	s.db.QueryRowContext(ctx, stmt,
		newVersion,
		ID,
	)

	return newVersion
}
