package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
)

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
	_, err = ths.db.ExecContext(
		ctx,
		`INSERT INTO "order" ("id", "user_id", "status", "sum") VALUES ($1, $2, $3, $4)`,
		order.ID,
		order.UserID,
		order.Status,
		order.Accural,
	)
	return
}

func (ths OrderRepository) FindByStatus(ctx context.Context, status diploma.OrderStatus) ([]*diploma.Order, error) {
	rows, err := ths.db.
		QueryContext(
			ctx,
			`SELECT "id", "user_id", "status" FROM "order" WHERE "status" = $1`,
			status,
		)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return ths.rowsToOrders(rows)
}

func (ths OrderRepository) UpdateStatus(ctx context.Context, orderID uint64, status diploma.OrderStatus, accural float64) (err error) {
	_, err = ths.db.ExecContext(
		ctx,
		`UPDATE "order" SET "status" = $1, "accural" = $2 WHERE "id" = $3`,
		status,
		accural,
		orderID,
	)
	return
}

func (ths OrderRepository) FindByUserID(ctx context.Context, userID uint64) ([]*diploma.Order, error) {
	rows, err := ths.db.
		QueryContext(
			ctx,
			`SELECT "id", "user_id", "status" FROM "order" WHERE "user_id" = $1`,
			userID,
		)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return ths.rowsToOrders(rows)
}

func (ths OrderRepository) rowsToOrders(rows *sql.Rows) ([]*diploma.Order, error) {
	result := make([]*diploma.Order, 0)
	for rows.Next() {
		o := &diploma.Order{}
		err := rows.Scan(&o.ID, &o.UserID, &o.Status)
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
