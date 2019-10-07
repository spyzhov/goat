package prometheus

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "prometheus",
		Name:         "Prometheus",
		Package:      "github.com/prometheus/client_golang",
		Dependencies: []string{"webserver"},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries: []*templates.Library{
			{Name: "github.com/prometheus/client_golang/prometheus/promhttp", Repo: "github.com/prometheus/client_golang", Version: "v0.9.4"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			if config.IsEnabled("http") || config.IsEnabled("httprouter") {
				s = `
	if err = app.setPrometheus(); err != nil {
		app.Logger.Error("cannot register Prometheus", zap.Error(err))
		return nil, err
	}`
			}
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			if config.IsEnabled("http") {
				s = `
// Set metrics
func (app *Application) setPrometheus() error {
	app.Logger.Debug("Prometheus registered")
	app.Http.Handle("/metrics", promhttp.Handler())
	return nil
}`
			} else if config.IsEnabled("httprouter") {
				s = `
// Set metrics
func (app *Application) setPrometheus() error {
	app.Logger.Debug("Prometheus registered")
	app.Router.Handler("GET", "/metrics", promhttp.Handler())
	return nil
}`
			}
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers:     templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{}
			if config.IsEnabled("fasthttp") {
				strings["metrics/metrics.go"] = `package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var (
	RequestTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_request_count_total",
		Help: "Count of online connections",
	})
	RequestStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_status_total",
		Help: "Requests status",
	}, []string{"method", "code"})
)

func init() {
	prometheus.MustRegister(RequestTotal)
	prometheus.MustRegister(RequestStatus)
}

func RequestInc() {
	RequestTotal.Inc()
}

func RequestStatusInc(method string, code int) {
	RequestStatus.WithLabelValues(method, strconv.Itoa(code)).Inc()
}
`
			}
			return
		},
	}
}
