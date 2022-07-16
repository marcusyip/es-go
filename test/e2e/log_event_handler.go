package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
	"go.uber.org/zap"
)

type LogEventHandler struct {
	es.BaseEventHandler
	logger *zap.Logger
}

func NewLogEventHandler(config *es.Config, logger *zap.Logger) *LogEventHandler {
	return &LogEventHandler{
		BaseEventHandler: *es.NewBaseEventHandler(config),
		logger:           logger,
	}
}

func (p *LogEventHandler) Handle(ctx context.Context, event es.Event) error {
	p.logger.Info("Event handled", zap.Any("event", event))
	return nil
}
