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

func (r *TransactionRepository) GetTx(ctx context.Context) es.DBTX {
	tx := es.GetContextTx(ctx)
	if tx == nil {
		return r.db
	}
	return tx
}

const getTransactionSQL = `-- name: CreateTransaction :one
SELECT id, version, status, currency, amount, done_by, created_at, updated_at
FROM transaction_views
WHERE id = $1 LIMIT 1
`

func (r *TransactionRepository) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	row := r.GetTx(ctx).QueryRow(ctx, getTransactionSQL, id)
	var m Transaction
	if err := row.Scan(&m.ID, &m.Version, &m.Status, &m.Currency, &m.Amount, &m.DoneBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}

const createTransactionSQL = `-- name: CreateTransaction :one
INSERT INTO transaction_views (
  id, version, status, currency, amount, done_by, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, version, status, currency, amount, done_by, created_at, updated_at
`

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	tx := r.GetTx(ctx)
	row := tx.QueryRow(ctx, createTransactionSQL, transaction.ID, transaction.Version, transaction.Status, transaction.Currency, transaction.Amount, transaction.DoneBy, transaction.CreatedAt, transaction.UpdatedAt)
	var m Transaction
	if err := row.Scan(&m.ID, &m.Version, &m.Status, &m.Currency, &m.Amount, &m.DoneBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}

const updateTransactionSQL = `-- name: UpdateTransaction :one
UPDATE transaction_views
SET version = $2, status = $3, currency = $4, amount = $5, done_by = $6, created_at = $7, updated_at = $8
WHERE id = $1 
`

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	tx := r.GetTx(ctx)
	_, err := tx.Exec(ctx, updateTransactionSQL, transaction.ID, transaction.Version, transaction.Status, transaction.Currency, transaction.Amount, transaction.DoneBy, transaction.CreatedAt, transaction.UpdatedAt)
	return err
}
