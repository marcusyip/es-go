package es

import (
	"context"
)

type Projector interface {
	Handle(ctx context.Context, tx DBTX, event Event) error
}

type BaseProjector struct {
	config *Config
}

func NewBaseProjector(config *Config) *BaseProjector {
	return &BaseProjector{config: config}
}

func (h *BaseProjector) Handle(ctx context.Context, tx *DBTX, event Event) error {
	return nil
}
