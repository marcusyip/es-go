package e2e

import (
	"context"
	"errors"

	"github.com/es-go/es-go/es"
	"github.com/es-go/es-go/es/database"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/segmentio/ksuid"
)

var _ = Describe("CommandService", func() {
	var config *es.Config
	var db *pgxpool.Pool
	var aggregateRepository es.AggregateRepository
	var commandService *es.CommandService
	var testID string

	BeforeEach(func() {
		testID = ksuid.New().String()
		// Setup
		config = es.NewConfig()
		db = database.Connect()

		eventRegistry := es.NewEventRegistry()
		eventRegistry.Set("created_event", &CreatedEvent{})
		eventRegistry.Set("completed_event", &CompletedEvent{})

		transactor := es.NewTransactor(db)

		aggregateRepository = es.NewAggregateRepository(config, db, transactor, eventRegistry)

		commandService = es.NewCommandService()
		commandService.Register("create_command", NewCreateCommandHandler(aggregateRepository))
		commandService.Register("complete_command", NewCompleteCommandHandler(aggregateRepository))
	})

	Context("Commit", func() {
		It("should not return error", func() {
			var command es.Command
			command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
			err := commandService.Execute(context.Background(), command)
			Expect(err).NotTo(HaveOccurred())
			// TODO: assert transaction_views table

			command = &CompleteCommand{TransactionID: testID, DoneBy: "marcusyip"}
			err = commandService.Execute(context.Background(), command)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Invalid command", func() {
		It("returns error", func() {
			var command es.Command
			command = &CreateCommand{TransactionID: testID, Currency: "", Amount: 1.11}
			err := commandService.Execute(context.Background(), command)
			Expect(errors.As(err, &validator.ValidationErrors{})).To(BeTrue())
		})
	})
})
