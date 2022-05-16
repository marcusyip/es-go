package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
)

type TransactionProjector struct {
	es.BaseProjector

	config                *es.Config
	transactionRepository *TransactionRepository
}

func NewTransactionProjector(config *es.Config, transactionRepository *TransactionRepository) *TransactionProjector {
	return &TransactionProjector{
		config:                config,
		transactionRepository: transactionRepository,
	}
}

func (p *TransactionProjector) Handle(ctx context.Context, tx es.DBTX, event es.Event) error {
	transaction, _ := ctx.Value("aggregate").(*Transaction)
	switch event.(type) {
	case *CreatedEvent:
		_, err := p.transactionRepository.CreateTransaction(ctx, tx, transaction)
		return err
	case *CompletedEvent:
		return p.transactionRepository.UpdateTransaction(ctx, tx, transaction)
	}
	return nil
}
