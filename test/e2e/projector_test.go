package e2e

import (
	"context"
	"errors"

	"github.com/es-go/es-go/es"
	"github.com/es-go/es-go/es/database"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
)

type MockProjector struct {
	mock.Mock
}

func (m *MockProjector) Handle(ctx context.Context, event es.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

var _ = Describe("Projector", func() {
	var config *es.Config
	var db *pgxpool.Pool
	var aggregateRepository es.AggregateRepository
	var transactionRepository *TransactionRepository
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
		transactionRepository = NewTransactionRepository(db)

		commandService = es.NewCommandService()
		commandService.Register("create_command", NewCreateCommandHandler(aggregateRepository))
		commandService.Register("complete_command", NewCompleteCommandHandler(aggregateRepository))
	})

	Context("Transaction does not exist", func() {
		BeforeEach(func() {
			transactionProjector := NewTransactionProjector(config, transactionRepository)
			aggregateRepository.AddProjector("created_event", transactionProjector)
			aggregateRepository.AddProjector("completed_event", transactionProjector)
		})

		It("creates event and transaction", func() {
			var command es.Command
			command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
			err := commandService.Execute(context.Background(), command)
			Expect(err).ToNot(HaveOccurred())

			transaction, err := transactionRepository.GetTransaction(context.Background(), testID)
			Expect(err).ToNot(HaveOccurred())
			Expect(transaction).ToNot(BeNil())
			Expect(transaction.ID).To(Equal(testID))
			Expect(transaction.Version).To(Equal(1))
			Expect(transaction.Status).To(Equal("processing"))
			Expect(transaction.Currency).To(Equal("BTC"))
			Expect(transaction.Amount).To(Equal(1.11))
			Expect(transaction.DoneBy).To(BeEmpty())
			Expect(transaction.CreatedAt).ToNot(BeNil())
			Expect(transaction.UpdatedAt).ToNot(BeNil())

			command = &CompleteCommand{TransactionID: testID, DoneBy: "marcusyip"}
			err = commandService.Execute(context.Background(), command)
			Expect(err).ToNot(HaveOccurred())

			transaction, err = transactionRepository.GetTransaction(context.Background(), testID)
			Expect(err).ToNot(HaveOccurred())
			Expect(transaction).ToNot(BeNil())
			Expect(transaction.ID).To(Equal(testID))
			Expect(transaction.Version).To(Equal(2))
			Expect(transaction.Status).To(Equal("completed"))
			Expect(transaction.Currency).To(Equal("BTC"))
			Expect(transaction.Amount).To(Equal(1.11))
			Expect(transaction.DoneBy).To(Equal("marcusyip"))
			Expect(transaction.CreatedAt).ToNot(BeNil())
			Expect(transaction.UpdatedAt).ToNot(BeNil())
		})
	})

	Context("With using mock projector", func() {
		var mockProjector *MockProjector

		BeforeEach(func() {
			mockProjector = new(MockProjector)
			mockProjector.On("Handle", mock.Anything, mock.Anything).Return(nil)

			aggregateRepository.Subscribe("created_event", mockProjector)
			aggregateRepository.Subscribe("completed_event", mockProjector)
		})

		It("calls projector once", func() {
			var command es.Command
			command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
			err := commandService.Execute(context.Background(), command)
			Expect(err).ToNot(HaveOccurred())
			mockProjector.AssertNumberOfCalls(GinkgoT(), "Handle", 1)
		})
	})

	Context("Project handles return error", func() {
		var mockProjector *MockProjector

		BeforeEach(func() {
			mockProjector = new(MockProjector)
			mockProjector.On("Handle", mock.Anything, mock.Anything).Return(errors.New("any error"))

			aggregateRepository.AddProjector("created_event", mockProjector)
		})

		It("rollback the transaction", func() {
			var command es.Command
			command = &CreateCommand{TransactionID: testID, Currency: "BTC", Amount: 1.11}
			err := commandService.Execute(context.Background(), command)
			Expect(err).To(HaveOccurred())
			mockProjector.AssertNumberOfCalls(GinkgoT(), "Handle", 1)
			// Assert no event should be created
			var events []*es.EventModel
			events, err = aggregateRepository.ListEvents(context.Background(), testID, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(0))
			// Assert no transaction should be created
			transaction, err := transactionRepository.GetTransaction(context.Background(), testID)
			Expect(err).To(MatchError(pgx.ErrNoRows))
			Expect(transaction).To(BeNil())
		})
	})
})
