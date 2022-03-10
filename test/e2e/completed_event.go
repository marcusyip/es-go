package e2e

import (
	"github.com/es-go/es-go/es"
)

type CompletedEvent struct {
	es.BaseEvent

	TransactionID string `json:"transaction_id"`
	DoneBy        string `json:"done_by"`
}

func (c *CompletedEvent) GetEventName() es.EventName        { return "completed_event" }
func (c *CompletedEvent) GetAggregateID() string            { return c.TransactionID }
func (c *CompletedEvent) SetAggregateID(aggregateID string) { c.TransactionID = aggregateID }
func (c *CompletedEvent) GetParentID() string               { return "" }
func (c *CompletedEvent) GetPayload() map[string]any {
	return map[string]any{
		"done_by": c.DoneBy,
	}
}
