package e2e

import (
	"context"

	"github.com/es-go/es-go/es"
	"github.com/es-go/es-go/es/database"
	"github.com/jackc/pgx/v4/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
)

type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) Handle(ctx context.Context, event es.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

var _ = Describe("EventHandler", func() {
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

	Context("with using mock event handler", func() {
		var mockEventHandler *MockEventHandler

		BeforeEach(func() {
			mockEventHandler = new(MockEventHandler)
			mockEventHandler.On("Handle", mock.Anything, mock.Anything).Return(nil)

			aggregateRepository.Subscribe("completed_event", mockEventHandler)
		})

		It("calls event handler once", func() {
			var command es.Command
			command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
			err := commandService.Execute(context.Background(), command)
			Expect(err).ToNot(HaveOccurred())
			mockEventHandler.AssertNumberOfCalls(GinkgoT(), "Handle", 0)

			command = &CompleteCommand{TransactionID: testID, DoneBy: "marcusyip"}
			err = commandService.Execute(context.Background(), command)
			Expect(err).ToNot(HaveOccurred())
			mockEventHandler.AssertNumberOfCalls(GinkgoT(), "Handle", 1)
		})
	})
})
