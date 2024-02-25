package diploma

import (
	"context"
	"database/sql"
	"errors"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}

func (ths OrderRepository) FindById(ctx context.Context, id uint64) (*Order, error) {
	u := new(Order)
	err := ths.db.
		QueryRowContext(
			ctx,
			`SELECT "id", "user_id" FROM "order" WHERE "id" = $1`,
			id,
		).
		Scan(&u.Id, &u.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ths OrderRepository) InsertNew(ctx context.Context, order *Order) (err error) {
	_, err = ths.db.ExecContext(
		ctx,
		`INSERT INTO "order" ("id", "user_id") VALUES ($1, $2)`,
		order.Id,
		order.UserId,
	)
	return
}

func (ths OrderRepository) FindByStatus(ctx context.Context, status OrderStatus) ([]*Order, error) {
	rows, err := ths.db.
		QueryContext(
			ctx,
			`SELECT "id", "user_id", "status" FROM "order" WHERE "status" = $1`,
			status,
		)
	defer rows.Close()

	result := make([]*Order, 0)
	for rows.Next() {
		o := &Order{}
		var status string
		err = rows.Scan(&o.Id, &o.UserId, status)
		if err != nil {
			return nil, err
		}
		o.Status, err = OrderStatusFromString(status)
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

func (ths OrderRepository) UpdateStatus(ctx context.Context, orderId uint64, status OrderStatus, accural float64) (err error) {
	_, err = ths.db.ExecContext(
		ctx,
		`UPDATE "order" SET "status" = $1, "accural" = $2 WHERE "id" = $3`,
		status,
		accural,
		orderId,
	)
	return
}
