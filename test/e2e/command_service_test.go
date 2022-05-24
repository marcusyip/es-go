package e2e

import (
	"context"
	"encoding/json"
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
			// 1. Execute CreateCommand
			var command es.Command
			command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
			err := commandService.Execute(context.Background(), command)
			Expect(err).NotTo(HaveOccurred())

			// 2. Assert version 1 created_event
			var events []*es.EventModel
			events, err = aggregateRepository.ListEventsByAggregateIDVersion(context.Background(), testID, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			event1 := events[0]
			Expect(event1.Version).To(Equal(1))
			Expect(event1.EventType).To(Equal("created_event"))
			var payload1 map[string]any
			_ = json.Unmarshal([]byte(event1.Payload), &payload1)
			Expect(payload1).To(Equal(map[string]any{
				"currency": "BTC",
				"amount":   1.11,
			}))
			Expect(event1.CreatedAt).ToNot(BeNil())

			// 3. Execute CompleteCommand
			command = &CompleteCommand{TransactionID: testID, DoneBy: "marcusyip"}
			err = commandService.Execute(context.Background(), command)
			Expect(err).NotTo(HaveOccurred())

			// 4. Assert version 2 completed_event
			events, err = aggregateRepository.ListEventsByAggregateIDVersion(context.Background(), testID, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			event2 := events[0]
			Expect(event2.Version).To(Equal(2))
			Expect(event2.EventType).To(Equal("completed_event"))
			var payload2 map[string]any
			_ = json.Unmarshal([]byte(event2.Payload), &payload2)
			Expect(payload2).To(Equal(map[string]any{
				"done_by": "marcusyip",
			}))
			Expect(event2.CreatedAt).ToNot(BeNil())
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
