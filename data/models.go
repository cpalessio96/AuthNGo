package data

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

type Models struct {
	User    User
	Session Session
}

func New(dbPool *sql.DB) Models {
	return Models{
		User:    User{db: dbPool},
		Session: Session{db: dbPool},
	}
}
