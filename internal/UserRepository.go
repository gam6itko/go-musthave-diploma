package internal

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	Id           int64
	Login        string
	PasswordHash []byte
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ths UserRepository) FindByLogin(ctx context.Context, login string) (*User, error) {
	u := new(User)
	err := ths.db.
		QueryRowContext(
			ctx,
			`SELECT "id", "login", "password" FROM "user" WHERE "login" = $1`,
			login,
		).
		Scan(&u.Id, &u.Login, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
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
