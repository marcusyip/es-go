package e2e

import (
	"fmt"
	"time"

	"github.com/es-go/es-go/es"
	esvalidator "github.com/es-go/es-go/es/validator"
	"github.com/go-playground/validator/v10"
)

var transactionStates []es.State = []es.State{
	"processing",
	"completed",
}

var transactionTransitions []es.Transition = []es.Transition{
	{FromState: "initialized", ToState: "processing", EventName: "created_event"},
	{FromState: "processing", ToState: "completed", EventName: "completed_event"},
}

type Transaction struct {
	es.BaseAggregateRoot

	Status    string
	Currency  string
	Amount    float64
	DoneBy    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTransaction() *Transaction {
	return &Transaction{Status: "initialized"}
}

func (t *Transaction) StateMachineEnabled() bool { return true }
func (t *Transaction) GetCurrentState() es.State {
	return es.State(t.Status)
}

func (t *Transaction) GetStates() []es.State {
	return transactionStates
}

func (t *Transaction) GetTransitions() []es.Transition {
	return transactionTransitions
}

func (t *Transaction) Create(currency string, amount float64) {
	t.applyChange(&CreatedEvent{
		Currency: currency,
		Amount:   amount,
	})
}

func (t *Transaction) Complete(doneBy string) {
	t.applyChange(&CompletedEvent{
		DoneBy: doneBy,
	})
}

func (t *Transaction) applyChange(event es.Event) {
	event.SetAggregateID(t.ID)
	event.SetVersion(t.NextVersion())
	event.SetCreatedAt(time.Now())

	validate := esvalidator.New()
	err := validate.Struct(event)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			panic(err)
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.Field())     // by passing alt name to ReportError like below
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		panic(err)
	}

	prevState := t.GetCurrentState()
	t.ApplyEvent(event)

	nextState := t.GetCurrentState()

	eventName := event.GetEventName()
	if !es.IsValidStateMachineTransition(t, prevState, nextState, eventName) {
		panic(fmt.Errorf("invalid state transition, from_state=%s, to_state=%s, event_name=%s", prevState, nextState, eventName))
	}

	t.AppendChange(event)
}

func (t *Transaction) ApplyEvent(event es.Event) {
	switch event := event.(type) {
	case *CreatedEvent:
		t.ID = event.GetAggregateID()
		t.Status = "processing"
		t.Currency = event.Currency
		t.Amount = event.Amount
		t.CreatedAt = event.GetCreatedAt()
		t.UpdatedAt = event.GetCreatedAt()
	case *CompletedEvent:
		t.Status = "completed"
		t.DoneBy = event.DoneBy
		t.UpdatedAt = event.GetCreatedAt()
	}
	t.SetVersion(event.GetVersion())
}
