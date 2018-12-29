package http

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "http",
		Name:    "HTTP server",
		Package: "net/http",

		Environments: []*templates.Environment{
			{Name: "Port", Type: "int", Env: "PORT", Default: "4000"},
		},
		Properties: []*templates.Property{
			{Name: "Http", Type: "*http.ServeMux", Default: "http.NewServeMux()"},
		},
		Libraries: []*templates.Library{
			{Name: "net/http"},
		},
		Models: map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction: func(config *templates.Config) (s string) {
			s = `	// Run HTTP server
	if err = application.RunHttp(); err != nil {
		application.Logger.Panic("HTTP Server start error", zap.Error(err))
	}`
			return
		},

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/http.go": `package app

import (
	"fmt"
	"net/http"
	"strconv"
)

func (a *Application) RunHttp() error {
	a.Http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Not implemented")
		w.WriteHeader(http.StatusNotImplemented)
	})
	go func() {
		a.Logger.Info("http server started on [::]:"+strconv.Itoa(a.Config.Port))
		a.Error <- http.ListenAndServe(":"+strconv.Itoa(a.Config.Port), a.Http)
	}()
	return nil
}
`,
			}
			return
		},
	}
}
