package es

import "context"

type AggregateLoader[T AggregateRoot] interface {
	Load(ctx context.Context, aggregateID string) (T, error)
}
