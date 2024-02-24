package internal

import (
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ths UserRepository) IsExists(ctx context.Context, login string) (exists bool, err error) {
	var cnt int
	err = ths.db.
		QueryRowContext(
			ctx,
			`SELECT COUNT("id") FROM "user" WHERE "login" = $1`,
			login,
		).
		Scan(&cnt)
	if err != nil {
		return
	}

	exists = cnt > 0
	return
}

func (ths UserRepository) InsertNew(ctx context.Context, login string, hashPass string) (id int64, err error) {
	err = ths.db.
		QueryRowContext(
			ctx,
			`INSERT INTO "user" ("login", "password") VALUES ($1, $2) RETURNING "id"`,
			login,
			hashPass,
		).
		Scan(&id)

	return
}
