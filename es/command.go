package es

type Command interface {
	// Command name is unqiue in command service registry
	GetCommandName() string
	// AggreagateID
	GetAggregateID() string
}

type BaseCommand struct {
	AggregateID string
}
