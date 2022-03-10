package es

import "context"

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

type BaseEventHandler struct {
	config *Config
}

func NewBaseEventHandler(config *Config) *BaseEventHandler {
	return &BaseEventHandler{config: config}
}

func (h *BaseEventHandler) Handle(ctx context.Context, event Event) error {
	return nil
}
