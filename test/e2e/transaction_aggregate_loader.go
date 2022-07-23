package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
)

type TransactionAggregateLoader struct {
	transactionRepository *TransactionRepository
}

func NewTransactionAggregateLoader(transactionRepository *TransactionRepository) *TransactionAggregateLoader {
	return &TransactionAggregateLoader{transactionRepository: transactionRepository}
}

func (l *TransactionAggregateLoader) Load(
	ctx context.Context,
	aggregateID string,
	loadOption *es.LoadOption,
) (*Transaction, error) {
	return l.transactionRepository.GetTransaction(ctx, aggregateID)
}
