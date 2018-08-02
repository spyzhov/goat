package prometheus

import "github.com/spyzhov/goat/templates"

var Env []templates.Environment
var Props = []templates.Property{
	{Name: "Echo", Type: "*echo.Echo", Default: "echo.New()"},
}
var Libs = []templates.Library{
	{Name: "github.com/labstack/echo"},
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
	a.Echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	return nil
}`
var TemplateRunFunction = `	// Run HTTP Server
	if err = application.RunHttp(); err != nil {
		application.Logger.Fatal("Echo start error", zap.Error(err))
	}`
var Templates = map[string]string{
	"app/http.go": `package app

func (a *Application) RunHttp() error {
	go func() {
		a.Error <- a.Echo.Start(":4000")
	}()
	return nil
}
`,
}
