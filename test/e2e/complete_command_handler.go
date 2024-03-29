package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
)

type CompleteCommandHandler struct {
	es.BaseCommandHandler

	repository es.AggregateRepository[*Transaction]
}

func NewCompleteCommandHandler(repository es.AggregateRepository[*Transaction]) *CompleteCommandHandler {
	return &CompleteCommandHandler{repository: repository}
}

func (h *CompleteCommandHandler) Handle(ctx context.Context, command es.Command) error {
	completeCommand := command.(*CompleteCommand)

	transaction, err := h.repository.Load(context.TODO(), completeCommand.GetAggregateID())
	if err != nil {
		return err
	}
	transaction.Complete(completeCommand.DoneBy)

	err = h.repository.Save(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
