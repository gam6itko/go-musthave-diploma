package diploma

import (
	"context"
	"database/sql"
	"errors"
)

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
		Scan(&u.ID, &u.Login, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ths UserRepository) FindByID(ctx context.Context, userID uint64) (*User, error) {
	u := new(User)
	err := ths.db.
		QueryRowContext(
			ctx,
			`SELECT "id", "login", "password", "balance_current", "balance_withdraw" FROM "user" WHERE "login" = $1`,
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
		`UPDATE "user SET 
	"balance_current" = "balance_current" - $1, 
	"balance_withdraw" = "balance_withdraw" + $1 
WHERE "id" = $2`,
		sum,
		userID,
	)
	//todo insert into withdrawals

	return
}
