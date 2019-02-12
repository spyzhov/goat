package prometheus

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "prometheus",
		Name:    "Prometheus",
		Package: "github.com/prometheus/client_golang",

		Environments: []*templates.Environment{
			{Name: "Port", Type: "int", Env: "PORT", Default: "4000"},
		},
		Properties: []*templates.Property{
			{Name: "Http", Type: "*http.ServeMux", Default: "http.NewServeMux()"},
		},
		Libraries: []*templates.Library{
			{Name: "net/http"},
			{Name: "github.com/prometheus/client_golang/prometheus/promhttp", Repo: "github.com/prometheus/client_golang", Version: "^0.9.2"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.setPrometheus(); err != nil {
		logger.Panic("cannot register Prometheus", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// Set metrics
func (a *Application) setPrometheus() error {
	a.Logger.Debug("Prometheus registered")
	a.Http.Handle("/metrics", promhttp.Handler())
	return nil
}`
			return
		},
		TemplateRunFunction: func(config *templates.Config) (s string) {
			if !config.IsEnabled("http") {
				s = `	// Run HTTP Server
	if err = application.RunHttp(); err != nil {
		application.Logger.Panic("HTTP Server start error", zap.Error(err))
	}`
			}
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
