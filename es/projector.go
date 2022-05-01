package es

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Projector interface {
	Handle(ctx context.Context, tx pgx.Tx, event Event) error
}

type BaseProjector struct {
	config *Config
}

func NewBaseProjector(config *Config) *BaseProjector {
	return &BaseProjector{config: config}
}

func (h *BaseProjector) Handle(ctx context.Context, tx *pgx.Tx, event Event) error {
	return nil
}
