package e2e

import (
	"context"
)

type TransactionAggregateLoader struct {
	transactionRepository *TransactionRepository
}

func NewTransactionAggregateLoader(transactionRepository *TransactionRepository) *TransactionAggregateLoader {
	return &TransactionAggregateLoader{transactionRepository: transactionRepository}
}

func (l *TransactionAggregateLoader) Load(ctx context.Context, aggregateID string) (*Transaction, error) {
	return l.transactionRepository.GetTransaction(ctx, aggregateID)
}
