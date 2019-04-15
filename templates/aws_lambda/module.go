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

		Environments: []*templates.Environment{
			{Name: "LambdaServerPort", Type: "int", Env: "_LAMBDA_SERVER_PORT"},
		},
		Properties: []*templates.Property{},
		Libraries:  []*templates.Library{},
		Models:     map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction: func(config *templates.Config) (s string) {
			s = `	// Run AWS-Lambda
	if err = app.RunLambda(); err != nil {
		app.Logger.Panic("AWS-Lambda start error", zap.Error(err))
	}`
			return
		},
		TemplateClosers: templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/lambda_handle.go": `package app

import "fmt"

// TODO Implement AWS-Lambda Handler
func (app *Application) lambdaHandle() (err error) {
	_, err = fmt.Print("Not implemented")
	return
}
`,
				"app/lambda.go": `package app

import (
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
)

// Start AWS-Lambda
func (app *Application) RunLambda() error {
	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()
		app.Logger.Info("AWS-Lambda starts")

		// app.WaitGroup.Add(1)
		go func() {
			// defer app.WaitGroup.Done()
			lambda.Start(app.lambdaHandle)
			app.Error <- errors.New("aws-lambda stops")
		}()

		select {
		case <-app.Ctx.Done():
			// todo: find the way to resolve lambda rpc.Accept releases
			app.Logger.Warn("aws-lambda won't stop")
			return
		}
	}()
	return nil
}
`,
			}
			return
		},
	}
}
