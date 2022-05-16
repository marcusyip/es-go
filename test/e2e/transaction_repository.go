package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

const getTransactionSQL = `-- name: CreateTransaction :one
SELECT id, status, currency, amount, done_by, created_at, updated_at
FROM transaction_views
WHERE id = $1 LIMIT 1
`

func (r *TransactionRepository) GetTransaction(ctx context.Context, tx es.DBTX, id string) (*Transaction, error) {
	row := tx.QueryRow(ctx, getTransactionSQL, id)
	var m Transaction
	if err := row.Scan(&m.ID); err != nil {
		return nil, err
	}
	return &m, nil
}

const createTransactionSQL = `-- name: CreateTransaction :one
INSERT INTO transaction_views (
  id, status, currency, amount, done_by, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, status, currency, amount, done_by, created_at, updated_at
`

func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx es.DBTX, transaction *Transaction) (*Transaction, error) {
	row := tx.QueryRow(ctx, createTransactionSQL, transaction.ID, transaction.Status, transaction.Currency, transaction.Amount, transaction.DoneBy, transaction.CreatedAt, transaction.UpdatedAt)
	var m Transaction
	if err := row.Scan(&m.ID, &m.Status, &m.Currency, &m.Amount, &m.DoneBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}

const updateTransactionSQL = `-- name: UpdateTransaction :one
UPDATE transaction_views
SET status = $2, currency = $3, amount = $4, done_by = $5, created_at = $6, updated_at = $7
WHERE id = $1 
`

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, tx es.DBTX, transaction *Transaction) error {
	_, err := tx.Exec(ctx, updateTransactionSQL, transaction.ID, transaction.Status, transaction.Currency, transaction.Amount, transaction.DoneBy, transaction.CreatedAt, transaction.UpdatedAt)
	return err
}
