package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
	"github.com/jackc/pgx/v4/pgxpool"
)

type EventLogger struct {
	config *es.Config
	db     *pgxpool.Pool
}

func NewEventLogger(config *es.Config, db *pgxpool.Pool) *EventLogger {
	return &EventLogger{
		config: config,
		db:     db,
	}
}

func (p *EventLogger) Handle(ctx context.Context, event es.Event) error {
	switch event.(type) {
	case *CompletedEvent:
		return nil
	}
	return nil
}
