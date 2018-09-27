package prometheus

import "github.com/spyzhov/goat/templates"

var Env []templates.Environment
var Props = []templates.Property{
	{Name: "Http", Type: "*http.ServeMux", Default: "http.NewServeMux()"},
}
var Libs = []templates.Library{
	{Name: "net/http"},
	{Name: "github.com/prometheus/client_golang/prometheus/promhttp"},
}
var Models = map[string]string{}

var TemplateSetter = `
	if err = app.setPrometheus(); err != nil {
		logger.Fatal("cannot register Prometheus", zap.Error(err))
		return nil, err
	}`
var TemplateSetterFunction = `
// Set metrics
func (a *Application) setPrometheus() error {
	a.Logger.Debug("Prometheus registered")
	a.Http.Handle("/metrics", promhttp.Handler())
	return nil
}`
var TemplateRunFunction = `	// Run HTTP Server
	if err = application.RunHttp(); err != nil {
		application.Logger.Fatal("HTTP Server start error", zap.Error(err))
	}`
var Templates = map[string]string{
	"app/http.go": `package app

import (
	"net/http"
)

func (a *Application) RunHttp() error {
	// TODO: Implement me
	go func() {
		a.Logger.Info("http server started on [::]:4000")
		a.Error <- http.ListenAndServe(":4000", a.Http)
	}()
	return nil
}
`,
}
