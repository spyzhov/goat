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
func (app *Application) setConsumer(consumer **RabbitMq, address, exchange, queueName, routingKey string) (err error) {
	app.Logger.Debug("RabbitMQ consumer connect", zap.String("address", address))
	(*consumer).Connection, err = amqp.Dial(address)
	if err != nil {
		return err
	}

	(*consumer).Channel, err = (*consumer).Connection.Channel()
	if err != nil {
		return err
	}

	err = (*consumer).Channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
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
		TemplateRunFunction: func(config *templates.Config) (s string) {
			s = `	// Run RabbitMQ Consumer
	if err = app.RunConsumer(app.Consumer, app.Config.ConsumerQueue); err != nil {
		app.Logger.Panic("RabbitMQ consumer start error", zap.Error(err))
	}`
			return
		},
		TemplateClosers: func(*templates.Config) (s string) {
			s = `
	defer func() {
		if app.Consumer != nil && app.Consumer.Connection != nil {
			app.closer("consumer connection", app.Consumer.Connection)
		}
	}()

	defer func() {
		if app.Consumer != nil && app.Consumer.Channel != nil {
			app.closer("consumer channel", app.Consumer.Channel)
		}
	}()`
			return
		},

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/consumer.go": `package app

import (
	"errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"runtime"
)

var (
	ConsumerStopError = errors.New("consumer stops")
)

func (app *Application) RunConsumer(consumer *RabbitMq, queue string) (err error) {
	app.Logger.Info("consumer start")
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

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		app.WaitGroup.Add(1)
		go func(i int) {
			defer app.WaitGroup.Done()
			app.Logger.Debug("consumer worker add", zap.Int("worker", i))
			for {
				select {
				case <-app.Ctx.Done():
					app.Logger.Debug("consumer worker close", zap.Int("worker", i))
					return
				case msg, ok := <-msgs:
					if ok {
						app.Logger.Debug("income message", zap.ByteString("message", msg.Body))
						if len(msg.Body) > 0 {
							app.consumerHandle(&msg)
						}
					} else {
						app.Error <- ConsumerStopError
						return
					}
				}
			}
		}(i)
	}

	return nil
}

func (app *Application) consumerHandle(msg *amqp.Delivery) {
	//TODO: Implement me
	defer func() {
		if err := msg.Ack(false); err != nil {
			app.Logger.Warn("error on ACK message", 
				zap.Error(err), 
				zap.Uint64("delivery_tag", msg.DeliveryTag))
		}
	}()
}
`,
			}
			return
		},
	}
}
