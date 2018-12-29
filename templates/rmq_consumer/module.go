package rmq_consumer

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "rmq_consumer",
		Name:    "RMQ-consumer",
		Package: "github.com/streadway/amqp",

		Environments: []*templates.Environment{
			{Name: "ConsumerAddress", Type: "string", Env: "CONSUMER_ADDR", Default: "amqp://guest:guest@localhost:5672"},
			{Name: "ConsumerExchange", Type: "string", Env: "CONSUMER_EXCHANGE"},
			{Name: "ConsumerQueue", Type: "string", Env: "CONSUMER_QUEUE"},
			{Name: "ConsumerRoutingKey", Type: "string", Env: "CONSUMER_ROUTING_KEY"},
		},
		Properties: []*templates.Property{
			{Name: "Consumer", Type: "*RabbitMq", Default: "new(RabbitMq)"},
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
	if err = app.setConsumer(&app.Consumer, app.Config.ConsumerAddress, app.Config.ConsumerExchange, app.Config.ConsumerQueue, app.Config.ConsumerRoutingKey); err != nil {
		logger.Panic("cannot connect to consumer RabbitMQ", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// RMQ set consumer
func (a *Application) setConsumer(consumer **RabbitMq, address, exchange, queueName, routingKey string) (err error) {
	a.Logger.Debug("RabbitMQ consumer connect", zap.String("address", address))
	(*consumer).Connection, err = amqp.Dial(address)
	if err != nil {
		return err
	}

	(*consumer).Channel, err = (*consumer).Connection.Channel()
	if err != nil {
		return err
	}

	queue, err := (*consumer).Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return err
	}
	(*consumer).Queue = &queue

	err = (*consumer).Channel.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		return err
	}

	// OnClose
	cerr := make(chan *amqp.Error)
	(*consumer).Channel.NotifyClose(cerr)
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
		TemplateRunFunction: func(config *templates.Config) (s string) {
			s = `	// Run RabbitMQ Consumer
	if err = application.RunConsumer(application.Consumer, application.Config.ConsumerQueue); err != nil {
		application.Logger.Panic("RabbitMQ consumer start error", zap.Error(err))
	}`
			return
		},

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/consumer.go": `package app

import (
	"go.uber.org/zap"
	"github.com/streadway/amqp"
)

func (a *Application) RunConsumer(consumer *RabbitMq, queue string) (err error) {
	a.Logger.Info("consumer start")
	msgs, err := consumer.Channel.Consume(
		queue,
		"{{.Name}}-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			consumer.Connection.Close()
		}()
		for {
			select {
			case msg, ok := <-msgs:
				if ok {
					a.Logger.Debug("income message", zap.ByteString("message", msg.Body))
					if len(msg.Body) > 0 {
						go a.consumerHandle(&msg)
					}
				}
			}
		}
	}()

	return nil
}

func (a *Application) consumerHandle(msg *amqp.Delivery) {
	//TODO: Implement me
	msg.Ack(false)
}
`,
			}
			return
		},
	}
}
