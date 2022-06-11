package es

import "context"

type AggregateLoader interface {
	Load(ctx context.Context, aggregateID string) (AggregateRoot, error)
}
