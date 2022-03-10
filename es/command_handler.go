package es

import "context"

type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}

type BaseCommandHandler struct {
	config *Config
}

func NewBaseCommandHandler(config *Config) *BaseCommandHandler {
	return &BaseCommandHandler{config: config}
}

func (h *BaseCommandHandler) Execute(ctx context.Context, command Command) error {
	return nil
}
