package http

import "github.com/spyzhov/goat/templates"

var Env []templates.Environment
var Props = []templates.Property{
	{Name: "Http", Type: "*http.ServeMux", Default: "http.NewServeMux()"},
}
var Libs = []templates.Library{
	{Name: "net/http"},
}
var Models = map[string]string{}

var TemplateSetter = ""
var TemplateSetterFunction = ""
var TemplateRunFunction = `	// Run HTTP server
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
