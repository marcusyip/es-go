package es

type AggregateRoot interface {
	StateMachine

	GetAggregateID() string
	SetAggregateID(id string)
	GetVersion() int
	ApplyEvent(event Event)
	GetChanges() []Event
	AppendChange(event Event)
}

type BaseAggregateRoot struct {
	AggregateRoot

	ID      string
	version int
	changes []Event

	states      []State
	transitions []Transition
}

func (r *BaseAggregateRoot) GetAggregateID() string {
	return r.ID
}

func (r *BaseAggregateRoot) SetAggregateID(aggregateID string) {
	r.ID = aggregateID
}

func (r *BaseAggregateRoot) GetVersion() int {
	return r.version
}

func (r *BaseAggregateRoot) SetVersion(version int) {
	r.version = version
}

func (r *BaseAggregateRoot) NextVersion() int {
	return r.version + 1
}

func (r *BaseAggregateRoot) GetChanges() []Event {
	return r.changes
}

func (r *BaseAggregateRoot) AppendChange(event Event) {
	r.changes = append(r.changes, event)
}
