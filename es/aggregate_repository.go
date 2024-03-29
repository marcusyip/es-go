package es

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type aggregateKey struct{}

// injectTx injects transaction to context
func WithContextAggregate(ctx context.Context, aggregate AggregateRoot) context.Context {
	return context.WithValue(ctx, aggregateKey{}, aggregate)
}

// extractTx extracts transaction from context
func GetContextAggregate(ctx context.Context) AggregateRoot {
	if aggregate, ok := ctx.Value(aggregateKey{}).(AggregateRoot); ok {
		return aggregate
	}
	return nil
}

type AggregateRepository[T AggregateRoot] interface {
	WithLoader(aggregateLoader AggregateLoader[T])
	ListEvents(ctx context.Context, aggregateID string, gteVersion int) ([]*EventModel, error)
	Load(ctx context.Context, aggregateID string) (T, error)
	Save(ctx context.Context, aggregate AggregateRoot) error
	AddProjector(eventName EventName, projector Projector)
	Subscribe(eventName EventName, eventHandler EventHandler)
}

type AggregateRepositoryImpl[T AggregateRoot] struct {
	config *Config
	// logger
	// logger *zap.Logger
	// custom aggregate load method
	aggregateLoader AggregateLoader[T]
	// new aggregate callback
	newAggregateFn func() T
	// placeholder of Load SQL statement
	loadSQL string
	db      *pgxpool.Pool
	// Transactor
	transactor *Transactor
	// Registry of event name and reflect.Type
	eventRegistry *EventRegistry
	// Projectors
	projectors map[EventName]([]Projector)
	// Event handler
	eventHandlers map[EventName]([]EventHandler)
}

func NewAggregateRepository[T AggregateRoot](config *Config, newAggregateFn func() T,
	db *pgxpool.Pool, transactor *Transactor, eventRegistry *EventRegistry,
) AggregateRepository[T] {
	loadSQL := fmt.Sprintf(
		`-- name: ListEvents :list
SELECT aggregate_id, version, event_type, payload, created_at
FROM %s 
WHERE aggregate_id = $1 and version > $2
ORDER BY version ASC
`,
		config.TableName)

	return &AggregateRepositoryImpl[T]{
		config:          config,
		aggregateLoader: nil,
		newAggregateFn:  newAggregateFn,
		loadSQL:         loadSQL,
		db:              db,
		transactor:      transactor,
		eventRegistry:   eventRegistry,
		eventHandlers:   map[EventName]([]EventHandler){},
		projectors:      map[EventName]([]Projector){},
	}
}

func (r *AggregateRepositoryImpl[T]) debug(format string, a ...any) {
	// fmt.Printf(format, a...)
}

type LoadFn func(ctx context.Context, aggregateID string, aggregate AggregateRoot) error

func (r *AggregateRepositoryImpl[T]) WithLoader(aggregateLoader AggregateLoader[T]) {
	r.aggregateLoader = aggregateLoader
}

func (r *AggregateRepositoryImpl[T]) GetTx(ctx context.Context) DBTX {
	tx := GetContextTx(ctx)
	if tx == nil {
		return r.db
	}
	return tx
}

func (r *AggregateRepositoryImpl[T]) ListEvents(ctx context.Context,
	aggregateID string, gteVersion int,
) ([]*EventModel, error) {
	tx := r.GetTx(ctx)
	rows, err := tx.Query(context.TODO(), r.loadSQL, aggregateID, gteVersion)
	if err != nil {
		r.debug("query error, err=%+v\n", err)
		return nil, fmt.Errorf("ListEvents tx.Query error: %w", err)
	}
	defer rows.Close()
	// don't know the size of rows
	eventModels := []*EventModel{}
	for rows.Next() {
		var m EventModel
		if err := rows.Scan(&m.AggregateID, &m.Version, &m.EventType,
			&m.Payload, &m.CreatedAt); err != nil {
			r.debug("ListEvents - scan err err=%+v\n", err)
			return nil, fmt.Errorf("ListEvents rows.Scan error: %w", err)
		}
		eventModels = append(eventModels, &m)
	}
	return eventModels, nil
}

func (r *AggregateRepositoryImpl[T]) Load(ctx context.Context, aggregateID string) (T, error) {
	if r.aggregateLoader != nil {
		return r.aggregateLoader.Load(ctx, aggregateID)
	}

	r.debug("Load aggregateID %s, sql=%s\n", aggregateID, r.loadSQL)
	aggregate := r.newAggregateFn()
	aggregate.SetAggregateID(aggregateID)

	mList, err := r.ListEvents(ctx, aggregateID, 0)
	if err != nil {
		var result T
		return result, fmt.Errorf("Load error: %w", err)
	}

	for _, m := range mList {
		eventType := r.eventRegistry.Get(m.EventType)
		event, ok := reflect.New(eventType).Interface().(Event)
		if !ok {
			panic("invalid event type casting")
		}

		err := json.Unmarshal(m.Payload, event)
		if err != nil {
			panic(err)
		}
		event.SetAggregateID(m.AggregateID)
		event.SetVersion(m.Version)
		event.SetCreatedAt(m.CreatedAt)

		// r.debug("============= Load - 1\n")
		// pp.Println(event)
		// pp.Println(aggregate)
		aggregate.ApplyEvent(event)
		// r.debug("============= Load - 2\n")
		// pp.Println(aggregate)
	}
	return aggregate, nil
}

func (r *AggregateRepositoryImpl[T]) Save(ctx context.Context, aggregate AggregateRoot) error {
	// TODO: use transactor
	return r.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		return r.doSave(ctx, aggregate)
	})
}

func (r *AggregateRepositoryImpl[T]) doSave(ctx context.Context, aggregate AggregateRoot) error {
	changes := aggregate.GetChanges()

	tx := r.GetTx(ctx)
	ctx = WithContextAggregate(ctx, aggregate)
	for _, change := range changes {
		commitSQL := fmt.Sprintf(
			"INSERT INTO %s (aggregate_id, version, event_type, payload, created_at) VALUES ($1, $2, $3, $4, $5)",
			r.config.TableName)
		payloadStr, err := json.Marshal(change.GetPayload())
		if err != nil {
			panic(err)
		}

		_, err = tx.Exec(ctx, commitSQL,
			change.GetAggregateID(),
			change.GetVersion(),
			change.GetEventName(),
			string(payloadStr),
			change.GetCreatedAt())
		if err != nil {
			return fmt.Errorf("Save tx.Exec error: %w", err)
		}
	}
	// projectView runs synchronously
	for _, change := range changes {
		r.debug("projecting view\n")
		if err := r.projectView(ctx, change); err != nil {
			return fmt.Errorf("Save projectView error: %w", err)
		}
	}
	// publishEvent runs synchronously
	for _, change := range changes {
		r.debug("publishing events\n")
		r.publishEvent(ctx, change)
	}
	return nil
}

func (r *AggregateRepositoryImpl[T]) AddProjector(eventName EventName, projector Projector) {
	if r.projectors[eventName] == nil {
		r.projectors[eventName] = make([]Projector, 0, 2)
	}
	r.projectors[eventName] = append(r.projectors[eventName], projector)
}

func (r *AggregateRepositoryImpl[T]) projectView(ctx context.Context, event Event) error {
	eventName := event.GetEventName()
	if r.projectors[eventName] == nil {
		return nil
	}
	for _, projector := range r.projectors[eventName] {
		err := projector.Handle(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AggregateRepositoryImpl[T]) Subscribe(eventName EventName, eventHandler EventHandler) {
	if r.eventHandlers[eventName] == nil {
		r.eventHandlers[eventName] = make([]EventHandler, 0, 2)
	}
	r.eventHandlers[eventName] = append(r.eventHandlers[eventName], eventHandler)
}

func (r *AggregateRepositoryImpl[T]) publishEvent(ctx context.Context, event Event) {
	eventName := event.GetEventName()
	if r.eventHandlers[eventName] == nil {
		return
	}
	for _, handler := range r.eventHandlers[eventName] {
		err := handler.Handle(ctx, event)
		// TODO: use goroutine and ignore error?
		if err != nil {
			r.debug("Event handler error %+v\n", err)
		}
	}
}

type EventModel struct {
	// ParentID       string
	AggregateID string `validate:"required"`
	// SubAggregateID string
	// ReferenceID    string
	EventType string    `validate:"required"`
	Version   int       `validate:"gt=0"`
	Payload   []byte    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
}
