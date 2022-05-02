package e2e

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

const getTransactionSQL = `-- name: CreateTransaction :one
SELECT id, status, currency, amount, done_by, created_at, updated_at
FROM transactions
WHERE id = $1 LIMIT 1
`

func (r *TransactionRepository) GetTransaction(ctx context.Context, tx pgx.Tx, id string) (*Transaction, error) {
	row := tx.QueryRow(ctx, getTransactionSQL, id)
	var m Transaction
	if err := row.Scan(&m.ID); err != nil {
		return nil, err
	}
	return &m, nil
}

const createTransactionSQL = `-- name: CreateTransaction :one
INSERT INTO transaction (
  id, status, currency, amount, done_by, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, status, currency, amount, done_by, created_at, updated_at
`

func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx pgx.Tx, transaction *Transaction) (*Transaction, error) {
	row := tx.QueryRow(ctx, createTransactionSQL, transaction.ID, transaction.Status, transaction.Currency, transaction.Amount, transaction.DoneBy, transaction.CreatedAt, transaction.UpdatedAt)
	var m Transaction
	if err := row.Scan(&m.ID, &m.Status, &m.Currency, &m.Amount, &m.DoneBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}
