package api_gateway

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:   "api_gateway",
		Name: "API Gateway",

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
func (app *Application) lambdaHandle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}
`,
			}
			return
		},
	}
}
