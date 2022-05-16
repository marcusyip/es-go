package es

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type txKey struct{}

// injectTx injects transaction to context
func WithContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func GetContextTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}

type Transactor struct {
	db *pgxpool.Pool
}

func NewTransactor(db *pgxpool.Pool) *Transactor {
	return &Transactor{db: db}
}

// func (r *Transactor) GetTx(ctx context.Context) DBTX {
// 	tx := GetContextTx(ctx)
// 	if tx == nil {
// 		return r.db
// 	}
// 	return tx
// }

func (t *Transactor) WithTransaction(ctx context.Context, callback func(ctx context.Context) error) error {
	// TODO: handle transaction in transaction

	// begin transaction
	tx, err := t.db.Begin(context.TODO())
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	err = callback(WithContextTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if err := tx.Rollback(context.TODO()); err != nil {
			log.Printf("rollback transaction: %v", err)
		}
		return err
	}

	if err := tx.Commit(context.TODO()); err != nil {
		log.Printf("commit transaction: %v", err)
	}
	return nil
}
