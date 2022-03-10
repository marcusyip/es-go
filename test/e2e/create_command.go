package e2e

import "github.com/es-go/es-go/es"

type CreateCommand struct {
	es.BaseCommand

	TransactionID string  `validate:"required"`
	Currency      string  `validate:"required"`
	Amount        float64 `validate:"required,gte=0"`
}

func (c *CreateCommand) GetCommandName() string { return "create_command" }
func (c *CreateCommand) GetAggregateID() string { return c.TransactionID }
