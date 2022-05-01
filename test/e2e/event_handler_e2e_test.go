package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/es-go/es-go/es"
	"github.com/es-go/es-go/es/database"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) Handle(ctx context.Context, event es.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func TestE2e_MockEventHandler_Success(t *testing.T) {
	testID := ksuid.New().String()
	// Setup
	config := es.NewConfig()
	db := database.Connect()

	eventRegistry := es.NewEventRegistry()
	eventRegistry.Set("created_event", &CreatedEvent{})
	eventRegistry.Set("completed_event", &CompletedEvent{})

	repository := es.NewAggregateRepository(config, db, eventRegistry)

	mockEventHandler := new(MockEventHandler)
	mockEventHandler.On("Handle", mock.Anything, mock.Anything).Return(nil)

	repository.Subscribe("completed_event", mockEventHandler)

	// Test
	service := es.NewCommandService()
	service.Register("create_command", NewCreateCommandHandler(repository))
	service.Register("complete_command", NewCompleteCommandHandler(repository))

	var command es.Command
	command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
	err := service.Execute(context.Background(), command)
	assert.NoError(t, err)
	// TODO: assert transaction_views table

	command = &CompleteCommand{TransactionID: testID, DoneBy: "marcusyip"}
	err = service.Execute(context.Background(), command)
	assert.NoError(t, err)

	// TODO: assert transaction_views table

	fmt.Println("TestE2e_SuccessMockEventHandler: done")
}
