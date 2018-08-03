package rmq_consumer

import "github.com/spyzhov/goat/templates"

var Env = []templates.Environment{
	{Name: "ConsumerAddress", Type: "string", Env: "CONSUMER_ADDR", Default: "\"amqp://guest:guest@localhost:5672\""},
	{Name: "ConsumerExchange", Type: "string", Env: "CONSUMER_EXCHANGE"},
	{Name: "ConsumerQueue", Type: "string", Env: "CONSUMER_QUEUE"},
	{Name: "ConsumerRoutingKey", Type: "string", Env: "CONSUMER_ROUTING_KEY"},
}
var Props = []templates.Property{
	{Name: "Consumer", Type: "*RabbitMq", Default: "new(RabbitMq)"},
}
var Libs = []templates.Library{
	{Name: "errors"},
	{Name: "github.com/streadway/amqp"},
}
var Models = map[string]string{
	"RabbitMq": `
type RabbitMq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}`,
}

var TemplateSetter = `
	if err = app.setConsumer(app.Consumer, app.Config.ConsumerAddress, app.Config.ConsumerExchange, app.Config.ConsumerQueue, app.Config.ConsumerRoutingKey); err != nil {
		logger.Fatal("cannot connect to consumer RabbitMQ", zap.Error(err))
		return nil, err
	}`
var TemplateSetterFunction = `
// RMQ set consumer
func (a *Application) setConsumer(consumer *RabbitMq, address, exchange, queueName, routingKey string) (err error) {
	a.Logger.Debug("RabbitMQ consumer connect", zap.String("address", address))
	consumer.Connection, err = amqp.Dial(address)
	if err != nil {
		return err
	}

	consumer.Channel, err = consumer.Connection.Channel()
	if err != nil {
		return err
	}

	queue, err := consumer.Channel.QueueDeclare(
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
	consumer.Queue = &queue

	err = consumer.Channel.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		return err
	}

	// OnClose
	cerr := make(chan *amqp.Error)
	consumer.Channel.NotifyClose(cerr)
	go func() {
		a.Error <- errors.New((<-cerr).Error())
		close(cerr)
	}()

	return nil
}`
var TemplateRunFunction = `	// Run RabbitMQ Consumer
	if err = application.RunConsumer(application.Consumer, application.Config.ConsumerQueue); err != nil {
		application.Logger.Fatal("RabbitMQ consumer start error", zap.Error(err))
	}`
var Templates = map[string]string{
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
			case msg := <-msgs:
				a.Logger.Debug("income message", zap.ByteString("message", msg.Body))
				if len(msg.Body) > 0 {
					go a.consumerHandle(&msg)
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
