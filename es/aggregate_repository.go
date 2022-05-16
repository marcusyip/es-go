package es

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/k0kubun/pp/v3"
)

type AggregateRepository interface {
	Load(aggregateID string, aggregate AggregateRoot) error
	Save(ctx context.Context, aggregate AggregateRoot) error
	SaveInTransaction(ctx context.Context, tx DBTX, aggregate AggregateRoot) error
	AddProjector(eventName EventName, projector Projector)
	Subscribe(eventName EventName, eventHandler EventHandler)
}

type AggregateRepositoryImpl struct {
	config *Config
	// placeholder of Load SQL statement
	loadSQL string
	db      *pgxpool.Pool
	// Registry of event name and reflect.Type
	eventRegistry *EventRegistry
	// Projectors
	projectors map[EventName]([]Projector)
	// Event handler
	eventHandlers map[EventName]([]EventHandler)
}

func NewAggregateRepository(config *Config, db *pgxpool.Pool, eventRegistry *EventRegistry) AggregateRepository {
	loadSQL := fmt.Sprintf(
		"SELECT aggregate_id, version, event_type, payload, created_at FROM %s WHERE aggregate_id = $1 ORDER BY version ASC",
		config.TableName)

	return &AggregateRepositoryImpl{
		config:        config,
		loadSQL:       loadSQL,
		db:            db,
		eventRegistry: eventRegistry,
		eventHandlers: map[EventName]([]EventHandler){},
		projectors:    map[EventName]([]Projector){},
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

func (r *AggregateRepositoryImpl) debug(format string, a ...any) {
	fmt.Printf(format, a...)
}

func (r *AggregateRepositoryImpl) Load(aggregateID string, aggregate AggregateRoot) error {
	r.debug("Load aggregateID %s, sql=%s\n", aggregateID, r.loadSQL)
	aggregate.SetAggregateID(aggregateID)

	// TODO: load aggregate by ID
	rows, err := r.db.Query(context.Background(), r.loadSQL, aggregateID)
	if err != nil {
		r.debug("query error, err=%+v\n", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var m EventModel
		if err := rows.Scan(&m.AggregateID, &m.Version, &m.EventType,
			&m.Payload, &m.CreatedAt); err != nil {
			r.debug("scan err err=%+v\n", err)

			return err
		}
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
	return nil
}

func (r *AggregateRepositoryImpl) Save(ctx context.Context, aggregate AggregateRoot) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		panic(err)
	}
	err = r.doSave(ctx, tx, aggregate, func(ctx context.Context, sql string, args ...any) error {
		result, err := r.db.Exec(ctx, sql, args...)
		pp.Println(result)
		return err
	})
	if err != nil {
		panic(err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *AggregateRepositoryImpl) SaveInTransaction(ctx context.Context, tx DBTX, aggregate AggregateRoot) error {
	return r.doSave(ctx, tx, aggregate, func(ctx context.Context, sql string, args ...any) error {
		result, err := tx.Exec(ctx, sql, args...)
		pp.Println(result)
		return err
	})
}

func (r *AggregateRepositoryImpl) doSave(ctx context.Context, tx DBTX, aggregate AggregateRoot, dbExecFn func(ctx context.Context, sql string, args ...any) error) error {
	changes := aggregate.GetChanges()
	pp.Println(changes)

	ctx = context.WithValue(ctx, "aggregate", aggregate)
	for _, change := range changes {
		commitSQL := fmt.Sprintf(
			"INSERT INTO %s (aggregate_id, version, event_type, payload, created_at) VALUES ($1, $2, $3, $4, $5)",
			r.config.TableName)

		payloadStr, _ := json.Marshal(change.GetPayload())
		err := dbExecFn(ctx, commitSQL,
			change.GetAggregateID(),
			change.GetVersion(),
			change.GetEventName(),
			string(payloadStr),
			change.GetCreatedAt())
		if err != nil {
			return err
		}
	}
	// projectView runs synchronously
	for _, change := range changes {
		r.projectView(ctx, tx, change)
	}
	// publishEvent runs synchronously
	for _, change := range changes {
		r.publishEvent(ctx, change)
	}
	return nil
}

func (r *AggregateRepositoryImpl) AddProjector(eventName EventName, projector Projector) {
	if r.projectors[eventName] == nil {
		r.projectors[eventName] = make([]Projector, 0, 2)
	}
	r.projectors[eventName] = append(r.projectors[eventName], projector)
}

func (r *AggregateRepositoryImpl) projectView(ctx context.Context, tx DBTX, event Event) {
	eventName := event.GetEventName()
	if r.projectors[eventName] == nil {
		return
	}
	for _, projector := range r.projectors[eventName] {
		err := projector.Handle(ctx, tx, event)
		if err != nil {
			panic(err)
		}
	}
}

func (r *AggregateRepositoryImpl) Subscribe(eventName EventName, eventHandler EventHandler) {
	if r.eventHandlers[eventName] == nil {
		r.eventHandlers[eventName] = make([]EventHandler, 0, 2)
	}
	r.eventHandlers[eventName] = append(r.eventHandlers[eventName], eventHandler)
}

func (r *AggregateRepositoryImpl) publishEvent(ctx context.Context, event Event) {
	eventName := event.GetEventName()
	if r.eventHandlers[eventName] == nil {
		return
	}
	for _, handler := range r.eventHandlers[eventName] {
		err := handler.Handle(ctx, event)
		if err != nil {
			panic(err)
		}
	}
}
