package e2e

import (
	"github.com/es-go/es-go/es"
)

type CreatedEvent struct {
	es.BaseEvent

	TransactionID string  `json:"transaction_id" validate:"required"`
	Currency      string  `json:"currency" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gte=0"`
}

func (c *CreatedEvent) GetEventName() es.EventName        { return "created_event" }
func (c *CreatedEvent) GetAggregateID() string            { return c.TransactionID }
func (c *CreatedEvent) SetAggregateID(aggregateID string) { c.TransactionID = aggregateID }
func (c *CreatedEvent) GetParentID() string               { return "" }
func (c *CreatedEvent) GetPayload() map[string]any {
	return map[string]any{
		"currency": c.Currency,
		"amount":   c.Amount,
	}
}
