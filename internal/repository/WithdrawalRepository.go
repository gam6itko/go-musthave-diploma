package repository

import (
	"context"
	"database/sql"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
)

type WithdrawalRepository struct {
	db *sql.DB
}

func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		db,
	}
}

func (ths WithdrawalRepository) FindByUserID(ctx context.Context, userID uint64) ([]*diploma.Withdrawal, error) {
	rows, err := ths.db.
		QueryContext(
			ctx,
			`SELECT "id", "user_id", "order_id", "processed_at", "sum" FROM "withdrawal" WHERE "user_id" = $1`,
			userID,
		)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return ths.rowsToOrders(rows)
}

func (ths WithdrawalRepository) rowsToOrders(rows *sql.Rows) ([]*diploma.Withdrawal, error) {
	result := make([]*diploma.Withdrawal, 0)
	for rows.Next() {
		w := &diploma.Withdrawal{}
		err := rows.Scan(&w.ID, &w.UserID, &w.OrderID, &w.ProcessedAt, &w.Sum)
		if err != nil {
			return nil, err
		}

		result = append(result, w)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	if err := rows.Close(); err != nil {
		return result, err
	}

	return result, nil
}
