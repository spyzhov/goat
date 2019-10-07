package aws_lambda

import (
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/aws_lambda/api_gateway"
	"github.com/spyzhov/goat/templates/aws_lambda/config_event"
	"github.com/spyzhov/goat/templates/aws_lambda/s3_event"
	"github.com/spyzhov/goat/templates/aws_lambda/ses_event"
	"github.com/spyzhov/goat/templates/aws_lambda/simple"
	"github.com/spyzhov/goat/templates/aws_lambda/sns_event"
	"github.com/spyzhov/goat/templates/aws_lambda/sqs_event"
)

func New() *templates.Template {
	return &templates.Template{
		ID:   "aws_lambda",
		Name: "AWS Lambda",
		Select: []*templates.Template{
			simple.New(),
			api_gateway.New(),
			config_event.New(),
			s3_event.New(),
			ses_event.New(),
			sns_event.New(),
			sqs_event.New(),
		},
		Conflicts: []string{
			"webserver",
			"rmq_consumer",
			"console",
		},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries: []*templates.Library{
			{Name: "github.com/aws/aws-lambda-go/lambda"},
		},
		Models: map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction:    templates.BlankFunction,
		TemplateClosers:        templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/lambda.go": `package app

import "context"

// TODO Implement AWS-Lambda Handler
func (app *Application) Lambda(ctx context.Context) (err error) {
	app.Logger.Warn("Not implemented")
	return
}
`,
			}
			return
		},
	}
}
