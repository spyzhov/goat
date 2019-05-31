package ses_event

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:   "ses_event",
		Name: "SES Events",

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
				"app/lambda.go": `package app

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

// TODO Implement AWS-Lambda Handler
func (app *Application) Lambda(ctx context.Context, sesEvent events.SimpleEmailEvent) error {
	for _, record := range sesEvent.Records {
		ses := record.SES
		fmt.Printf("[%s - %s] Mail = %+v, Receipt = %+v \n", record.EventVersion, record.EventSource, ses.Mail, ses.Receipt)
	}

	return nil
}
`,
			}
			return
		},
	}
}
