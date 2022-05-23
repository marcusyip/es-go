package es

import (
	"context"
	"fmt"

	esvalidator "github.com/es-go/es-go/es/validator"
	"github.com/go-playground/validator/v10"
)

type CommandService struct {
	commandHandlers map[string]CommandHandler
}

func NewCommandService() *CommandService {
	return &CommandService{
		commandHandlers: map[string]CommandHandler{},
	}
}

func (s *CommandService) Register(commandName string, handler CommandHandler) {
	s.commandHandlers[commandName] = handler
}

func (s *CommandService) Execute(ctx context.Context, command Command) error {
	validate := esvalidator.New()
	err := validate.Struct(command)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
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
		return err
	}
	commandName := command.GetCommandName()
	if _, ok := s.commandHandlers[commandName]; !ok {
		panic("unregistered command: " + commandName)
	}
	return s.commandHandlers[commandName].Handle(ctx, command)
}
