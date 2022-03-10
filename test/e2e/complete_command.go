package e2e

import "github.com/es-go/es-go/es"

type CompleteCommand struct {
	es.BaseCommand

	TransactionID string
	DoneBy        string
}

func (c *CompleteCommand) GetCommandName() string { return "complete_command" }
func (c *CompleteCommand) GetAggregateID() string { return c.TransactionID }
