package es

import "context"

type Projector interface {
	EventHandler
}

type BaseProjector struct {
	config *Config
}

func NewBaseProjector(config *Config) *BaseProjector {
	return &BaseProjector{config: config}
}

func (h *BaseProjector) Handle(ctx context.Context, event Event) error {
	return nil
}
