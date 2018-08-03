package babex

import "github.com/spyzhov/goat/templates"

var Env = []templates.Environment{
	{Name: "BabexAddr", Type: "string", Env: "BABEX_ADDR", Default: "amqp://guest:guest@localhost:5672/"},
	{Name: "BabexExchange", Type: "string", Env: "BABEX_EXCHANGE"},
	{Name: "BabexName", Type: "string", Env: "BABEX_NAME"},
}
var Props = []templates.Property{
	{Name: "Service", Type: "*babex.Service"},
}
var Libs = []templates.Library{
	{Name: "github.com/matroskin13/babex"},
}
var Models = map[string]string{}

var TemplateSetter = `
	if err = app.setBabex(app.Service); err != nil {
		logger.Fatal("cannot create Babex node", zap.Error(err))
		return nil, err
	}`
var TemplateSetterFunction = `
// Babex node connect
func (a *Application) setBabex(service *babex.Service) (err error) {
	a.Logger.Debug("Babex connect", zap.String("address", a.Config.BabexAddr))
	service, err = babex.NewService(&babex.ServiceConfig{
		Name:    a.Config.BabexName,
		Address: a.Config.BabexAddr,
	})
	if err != nil {
		return err
	}
	err = service.BindToExchange(a.Config.BabexExchange, a.Config.BabexName)
	if err != nil {
		return err
	}
	return nil
}`
var TemplateRunFunction = `	// Run Babex-node
	if err = application.RunBabex(); err != nil {
		application.Logger.Fatal("babex-node start error", zap.Error(err))
	}`
var Templates = map[string]string{
	"app/babex.go": `package app

import (
	"github.com/matroskin13/babex"
	"go.uber.org/zap"
)

func (a *Application) RunBabex() error {
	msgs, err := a.Service.GetMessages()
	if err != nil {
		a.Logger.Fatal("cannot get messages", zap.Error(err))
		return err
	}
	errs := a.Service.GetErrors()
	a.Logger.Info("service listen queue")

	go func() {
		for {
			select {
			case msg := <-msgs:
				err := a.receive(msg)
				if err != nil {
					a.Logger.Error("error on processing message", zap.Error(err))
				}
			case err := <-errs:
				a.Error <- err
			}
		}
	}()
	return nil
}

func (a *Application) receive(msg *babex.Message) (err error) {
	//region Implementation
	//TODO: Implement me
	//endregion
	//region Next
	err = a.Service.Next(msg, msg.Data, nil)
	if err != nil && err != babex.ErrorNextIsNotDefined {
		a.Logger.Error("cannot next", zap.Error(err))
		return err
	}
	//endregion
	return nil
}
`,
}
