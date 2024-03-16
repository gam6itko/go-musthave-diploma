package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ths UserRepository) FindByLogin(ctx context.Context, login string) (*diploma.User, error) {
	u := new(diploma.User)
	err := ths.db.
		QueryRowContext(
			ctx,
			`SELECT "id", "login", "password" FROM "user" WHERE "login" = $1`,
			login,
		).
		Scan(&u.ID, &u.Login, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ths UserRepository) FindByID(ctx context.Context, userID uint64) (*diploma.User, error) {
	u := new(diploma.User)
	err := ths.db.
		QueryRowContext(
			ctx,
			`SELECT "id", "login", "password", "balance_current", "balance_withdraw" FROM "user" WHERE "id" = $1`,
			userID,
		).
		Scan(&u.ID, &u.Login, &u.PasswordHash, &u.BalanceCurrent, &u.BalanceWithdraw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ths UserRepository) InsertNew(ctx context.Context, login string, hashPass string) (id uint64, err error) {
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

func (ths UserRepository) Withdraw(ctx context.Context, userID uint64, orderID uint64, sum float32) (err error) {
	tx, err := ths.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO "withdrawal" ("user_id", "order_id", "sum") VALUES ($1, $2, $3)`,
		userID,
		orderID,
		sum,
	)
	if err != nil {
		return
	}

	_, err = tx.Exec(
		`UPDATE "user" SET 
	"balance_current" = "balance_current" - $1, 
	"balance_withdraw" = "balance_withdraw" + $1 
WHERE "id" = $2`,
		sum,
		userID,
	)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}
