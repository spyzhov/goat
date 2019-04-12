package s3_event

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:   "s3_event",
		Name: "S3 Events",

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
func (app *Application) lambdaHandle(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key) 
	}
}
`,
			}
			return
		},
	}
}
