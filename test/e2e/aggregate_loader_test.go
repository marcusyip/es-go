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

type MockAggregateLoader struct {
	mock.Mock
}

func (m *MockAggregateLoader) Load(ctx context.Context, aggregateID string) (*Transaction, error) {
	args := m.Called(ctx, aggregateID)
	return args.Get(0).(*Transaction), args.Error(1)
}

var _ = Describe("AggregateLoader", func() {
	var config *es.Config
	var db *pgxpool.Pool
	var aggregateRepository es.AggregateRepository[*Transaction]
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
		aggregateRepository = es.NewAggregateRepository(config, NewTransaction, db, transactor, eventRegistry)
	})

	Context("with using mock event handler", func() {
		var mockAggregateLoader *MockAggregateLoader

		BeforeEach(func() {
			mockAggregateLoader = new(MockAggregateLoader)
			mockAggregateLoader.On("Load", mock.Anything, testID).Return(NewTransaction(), nil)

			aggregateRepository.WithLoader(mockAggregateLoader)
		})

		It("calls aggregate loader once", func() {
			_, err := aggregateRepository.Load(context.Background(), testID)
			Expect(err).NotTo(HaveOccurred())
			mockAggregateLoader.AssertNumberOfCalls(GinkgoT(), "Load", 1)
		})
	})

	// TODO: Test without aggregate loader

})
