package es

type State string

type Transition struct {
	FromState State
	ToState   State
	EventName EventName
}

type StateMachine interface {
	StateMachineEnabled() bool
	GetCurrentState() State
	GetStates() []State
	GetTransitions() []Transition
}

func IsValidStateMachineTransition(stateMachine StateMachine, fromState, toState State, eventName EventName) bool {
	transitions := stateMachine.GetTransitions()
	for _, transition := range transitions {
		if transition.FromState == fromState &&
			transition.ToState == toState &&
			transition.EventName == eventName {
			return true
		}
	}
	return false
}
