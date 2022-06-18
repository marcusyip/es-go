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

func (p *TransactionProjector) Handle(ctx context.Context, event es.Event) error {
	transaction := es.GetContextAggregate(ctx).(*Transaction)
	switch event.(type) {
	case *CreatedEvent:
		_, err := p.transactionRepository.CreateTransaction(ctx, transaction)
		return err
	case *CompletedEvent:
		return p.transactionRepository.UpdateTransaction(ctx, transaction)
	}
	return nil
}
