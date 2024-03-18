package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
)

type IOrderRepository interface {
	FindByID(ctx context.Context, orderID uint64) (*diploma.Order, error)
	InsertNew(ctx context.Context, order *diploma.Order) (err error)
	UpdateStatus(ctx context.Context, orderID uint64, status diploma.OrderStatus, accrual float64) (err error)
	FindByUserID(ctx context.Context, userID uint64) ([]diploma.Order, error)
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}

// В реальных проектах, конечно должна быть пагинация.
func (ths OrderRepository) FindByID(ctx context.Context, orderID uint64) (*diploma.Order, error) {
	u := new(diploma.Order)
	err := ths.db.
		QueryRowContext(
			ctx,
			`SELECT "id", "user_id" FROM "order" WHERE "id" = $1`,
			orderID,
		).
		Scan(&u.ID, &u.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ths OrderRepository) InsertNew(ctx context.Context, order *diploma.Order) (err error) {
	tx, err := ths.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO "order" ("id", "user_id", "status", "sum") VALUES ($1, $2, $3, $4)`,
		order.ID,
		order.UserID,
		order.Status,
		order.Accrual,
	)
	if err != nil {
		return
	}

	_, err = tx.Exec(
		`UPDATE "user" SET "balance_current" = "balance_current" + $1 WHERE "id" = $2`,
		order.Accrual,
		order.UserID,
	)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (ths OrderRepository) UpdateStatus(ctx context.Context, orderID uint64, status diploma.OrderStatus, accrual float64) (err error) {
	_, err = ths.db.ExecContext(
		ctx,
		`UPDATE "order" SET "status" = $1, "accrual" = $2 WHERE "id" = $3`,
		status,
		accrual,
		orderID,
	)
	return
}

func (ths OrderRepository) FindByUserID(ctx context.Context, userID uint64) ([]diploma.Order, error) {
	rows, err := ths.db.
		QueryContext(
			ctx,
			`SELECT "id", "uploaded_at", "user_id", "status", "sum" FROM "order" WHERE "user_id" = $1`,
			userID,
		)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToOrders(rows)
}

func rowsToOrders(rows *sql.Rows) ([]diploma.Order, error) {
	result := make([]diploma.Order, 0)
	for rows.Next() {
		o := diploma.Order{}
		err := rows.Scan(&o.ID, &o.UploadedAt, &o.UserID, &o.Status, &o.Accrual)
		if err != nil {
			return nil, err
		}

		result = append(result, o)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	if err := rows.Close(); err != nil {
		return result, err
	}

	return result, nil
}
