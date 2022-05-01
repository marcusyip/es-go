package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/es-go/es-go/es"
	"github.com/es-go/es-go/es/database"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestE2e_Success(t *testing.T) {
	testID := ksuid.New().String()
	// Setup
	config := es.NewConfig()
	db := database.Connect()

	eventRegistry := es.NewEventRegistry()
	eventRegistry.Set("created_event", &CreatedEvent{})
	eventRegistry.Set("completed_event", &CompletedEvent{})

	repository := es.NewAggregateRepository(config, db, eventRegistry)

	// Test
	service := es.NewCommandService()
	service.Register("create_command", NewCreateCommandHandler(repository))
	service.Register("complete_command", NewCompleteCommandHandler(repository))

	var command es.Command
	command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
	err := service.Execute(context.Background(), command)
	assert.NoError(t, err)

	command = &CompleteCommand{TransactionID: testID, DoneBy: "marcusyip"}
	err = service.Execute(context.Background(), command)
	assert.NoError(t, err)

	fmt.Println("TestE2e_Success: done")
}

func TestE2e_Projector_Success(t *testing.T) {
	testID := ksuid.New().String()
	// Setup
	config := es.NewConfig()
	db := database.Connect()

	eventRegistry := es.NewEventRegistry()
	eventRegistry.Set("created_event", &CreatedEvent{})
	eventRegistry.Set("completed_event", &CompletedEvent{})

	repository := es.NewAggregateRepository(config, db, eventRegistry)
	transactionProjector := NewTransactionProjector(config, db)
	repository.AddProjector("created_event", transactionProjector)
	repository.AddProjector("completed_event", transactionProjector)

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

	fmt.Println("TestE2e_SuccessProjector: done")
}

func TestE2e_FailInvalidCommand(t *testing.T) {
	testID := ksuid.New().String()
	// Setup
	config := es.NewConfig()
	db := database.Connect()

	eventRegistry := es.NewEventRegistry()
	eventRegistry.Set("created_event", &CreatedEvent{})
	eventRegistry.Set("completed_event", &CompletedEvent{})

	repository := es.NewAggregateRepository(config, db, eventRegistry)
	transactionProjector := NewTransactionProjector(config, db)
	repository.AddProjector("created_event", transactionProjector)
	repository.AddProjector("completed_event", transactionProjector)

	// Test
	service := es.NewCommandService()
	service.Register("create_command", NewCreateCommandHandler(repository))
	service.Register("complete_command", NewCompleteCommandHandler(repository))

	var command es.Command
	command = &CreateCommand{TransactionID: testID, Currency: "", Amount: 1.11}
	err := service.Execute(context.Background(), command)
	assert.Error(t, err, "should return invalid command error")

	// TODO: assert transaction_views table

	fmt.Println("TestE2e_FailInvalidCommand: done")
}
