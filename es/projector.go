package es

import (
	"context"
)

type Projector interface {
	Handle(ctx context.Context, event Event) error
}

type BaseProjector struct {
	config *Config
}

func NewBaseProjector(config *Config) *BaseProjector {
	return &BaseProjector{config: config}
}
