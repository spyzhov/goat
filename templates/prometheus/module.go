package prometheus

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "prometheus",
		Name:         "Prometheus",
		Package:      "github.com/prometheus/client_golang",
		Dependencies: []string{"http"},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries: []*templates.Library{
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
func (app *Application) setPrometheus() error {
	app.Logger.Debug("Prometheus registered")
	app.Http.Handle("/metrics", promhttp.Handler())
	return nil
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers:     templates.BlankFunction,

		Templates: templates.BlankFunctionMap,
	}
}
