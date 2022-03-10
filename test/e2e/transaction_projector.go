package e2e

import (
	"context"
	"fmt"

	"github.com/es-go/es-go/es"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/k0kubun/pp/v3"
)

type TransactionProjector struct {
	es.BaseProjector

	config *es.Config
	db     *pgxpool.Pool
}

func NewTransactionProjector(config *es.Config, db *pgxpool.Pool) *TransactionProjector {
	return &TransactionProjector{
		config: config,
		db:     db,
	}
}

func (p *TransactionProjector) Handle(ctx context.Context, event es.Event) error {
	transaction, _ := ctx.Value("aggregate").(*Transaction)
	switch event.(type) {
	case *CreatedEvent:
		sql := `
INSERT INTO transaction_views
	(id, status, currency, amount, done_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)`
		result, err := p.db.Exec(ctx, sql,
			transaction.ID,
			transaction.Status,
			transaction.Currency,
			transaction.Amount,
			transaction.DoneBy,
			transaction.CreatedAt,
			transaction.UpdatedAt)
		if err != nil {
			return err
		}
		pp.Println(result)
	case *CompletedEvent:
		fmt.Println("projector: completed")
		sql := `
UPDATE transaction_views
SET
	status=$1,
	currency=$2,
	amount=$3,
	done_by=$4,
	created_at=$5,
	updated_at=$6
WHERE id=$7`
		result, err := p.db.Exec(ctx, sql,
			transaction.Status,
			transaction.Currency,
			transaction.Amount,
			transaction.DoneBy,
			transaction.CreatedAt,
			transaction.UpdatedAt,
			transaction.ID)
		if err != nil {
			return err
		}
		pp.Println(result)
	}
	return nil
}
