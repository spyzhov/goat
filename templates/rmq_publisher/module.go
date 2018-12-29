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
func (a *Application) setPublisher(publisher **RabbitMq, address string) (err error) {
	a.Logger.Debug("RabbitMQ publisher connect", zap.String("address", address))
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
	go func() {
		err, ok := <-cerr
		if ok {
			a.Error <- errors.New(err.Error())
		}
	}()

	return nil
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/publish.go": `package app

import (
	"github.com/streadway/amqp"
)

func (a *Application) publish(body []byte) error {
	return a.Publisher.Channel.Publish(
		a.Config.PublisherExchange,   // publish to an exchange
		a.Config.PublisherRoutingKey, // routing to 0 or more queues
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
