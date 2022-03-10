package es

import (
	"fmt"
	"reflect"
	"time"
)

type EventName string

type Event interface {
	GetEventName() EventName
	GetParentID() string
	GetAggregateID() string
	SetAggregateID(aggregateID string)
	GetVersion() int
	SetVersion(version int)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetPayload() map[string]any
}

type BaseEvent struct {
	ParentID       string
	AggregateID    string
	SubAggregateID string
	ReferenceID    string
	Version        int
	CreatedAt      time.Time
	Payload        interface{}
}

func (c *BaseEvent) GetVersion() int                  { return c.Version }
func (c *BaseEvent) SetVersion(version int)           { c.Version = version }
func (c *BaseEvent) GetCreatedAt() time.Time          { return c.CreatedAt }
func (c *BaseEvent) SetCreatedAt(createdAt time.Time) { c.CreatedAt = createdAt }

type EventRegistry struct {
	eventTypes map[string]reflect.Type
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		eventTypes: make(map[string]reflect.Type),
	}
}

func (reg *EventRegistry) Set(eventName string, event Event) {
	if reg.eventTypes[eventName] != nil {
		panic(fmt.Errorf("EventRegistry: event %s already exists", eventName))
	}
	reg.eventTypes[eventName] = reflect.TypeOf(event).Elem()
}

func (reg *EventRegistry) Get(eventName string) reflect.Type {
	t, ok := reg.eventTypes[eventName]
	if !ok {
		panic("invalid eventName " + eventName)
	}
	return t
}
