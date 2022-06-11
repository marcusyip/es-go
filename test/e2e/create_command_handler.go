package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
)

type CreateCommandHandler struct {
	es.BaseCommandHandler

	repository es.AggregateRepository[*Transaction]
}

func NewCreateCommandHandler(repository es.AggregateRepository[*Transaction]) *CreateCommandHandler {
	return &CreateCommandHandler{repository: repository}
}

func (h *CreateCommandHandler) Handle(ctx context.Context, command es.Command) error {
	createCommand := command.(*CreateCommand)
	transaction, err := h.repository.Load(context.TODO(), createCommand.GetAggregateID())
	if err != nil {
		return err
	}
	transaction.Create(createCommand.Currency, createCommand.Amount)

	err = h.repository.Save(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
