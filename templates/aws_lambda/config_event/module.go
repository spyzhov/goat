package config_event

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:   "config_event",
		Name: "Config Events",

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries:    []*templates.Library{},
		Models:       map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction:    templates.BlankFunction,
		TemplateClosers:        templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/lambda_handle.go": `package app

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

// TODO Implement AWS-Lambda Handler
func (app *Application) lambdaHandle(ctx context.Context, configEvent events.ConfigEvent) {
    fmt.Printf("AWS Config rule: %s\n", configEvent.ConfigRuleName)
    fmt.Printf("Invoking event JSON: %s\n", configEvent.InvokingEvent)
    fmt.Printf("Event version: %s\n", configEvent.Version)
}
`,
			}
			return
		},
	}
}
