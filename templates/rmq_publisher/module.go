package rmq_publisher

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "rmq_publisher",
		Name:    "RMQ-publisher",
		Package: "github.com/streadway/amqp",

		Environments: []*templates.Environment{
			{Name: "PublisherAddress", Type: "string", Env: "PUBLISHER_ADDR", Default: "amqp://guest:guest@localhost:5672"},
			{Name: "PublisherExchange", Type: "string", Env: "PUBLISHER_EXCHANGE"},
			{Name: "PublisherRoutingKey", Type: "string", Env: "PUBLISHER_ROUTING_KEY"},
		},
		Properties: []*templates.Property{
			{Name: "Publisher", Type: "*RabbitMq", Default: "new(RabbitMq)"},
		},
		Libraries: []*templates.Library{
			{Name: "errors"},
			{Name: "github.com/streadway/amqp"},
		},
		Models: map[string]string{
			"RabbitMq": `
type RabbitMq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}`,
		},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.setPublisher(&app.Publisher, app.Config.PublisherAddress); err != nil {
		logger.Panic("cannot connect to publisher RabbitMQ", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// RMQ set publisher
func (app *Application) setPublisher(publisher **RabbitMq, address string) (err error) {
	app.Logger.Debug("RabbitMQ publisher connect", zap.String("address", address))
	(*publisher).Connection, err = amqp.Dial(address)
	if err != nil {
		return err
	}

	(*publisher).Channel, err = (*publisher).Connection.Channel()
	if err != nil {
		return err
	}

	// OnClose
	cerr := make(chan *amqp.Error)
	(*publisher).Channel.NotifyClose(cerr)

	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()
		select {
		case <-app.Ctx.Done():
			return
		case err, ok := <-cerr:
			if ok && err != nil {
				app.Error <- errors.New(err.Error())
			}
		}
	}()

	return nil
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers: func(*templates.Config) (s string) {
			s = `
	defer func() {
		if app.Publisher != nil && app.Publisher.Connection != nil {
			app.Closer("publisher connection", app.Publisher.Connection)
		}
	}()

	defer func() {
		if app.Publisher != nil && app.Publisher.Channel != nil {
			app.Closer("publisher channel", app.Publisher.Channel)
		}
	}()`
			return
		},

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/publish.go": `package app

import (
	"github.com/streadway/amqp"
)

func (app *Application) publish(body []byte) error {
	return app.Publisher.Channel.Publish(
		app.Config.PublisherExchange,   // publish to an exchange
		app.Config.PublisherRoutingKey, // routing to 0 or more queues
		false,                        // mandatory
		false,                        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
			// a bunch of application/implementation-specific fields
		},
	)
}
`,
			}
			return
		},
	}
}
