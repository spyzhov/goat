package webserver

import (
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/webserver/fasthttp"
	"github.com/spyzhov/goat/templates/webserver/http"
	"github.com/spyzhov/goat/templates/webserver/httprouter"
)

func New() *templates.Template {
	return &templates.Template{
		ID:   "webserver",
		Name: "WebServer",
		Select: []*templates.Template{
			http.New(),
			httprouter.New(),
			fasthttp.New(),
		},
		Conflicts: []string{"aws_lambda"},

		Environments: []*templates.Environment{
			{Name: "Port", Type: "int", Env: "PORT", Default: "4000"},
		},
		Properties: []*templates.Property{},
		Libraries:  []*templates.Library{},
		Models:     map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction:    templates.BlankFunction,
		TemplateClosers:        templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{}
			strings["app/healthcheck.go"] = `package app

import (
	"net/http"
	"time"
)

// Handle function for health-check
func (app *Application) healthCheck() (info map[string]string, status int) {
	status = http.StatusOK
	info = map[string]string{
		"service": "{{.Name}}",
		"time":    time.Now().String(),
` + templates.Str(config.IsEnabled("lib-postgres"), `

		"postgres": (func() string {
			if err := app.Postgres.Ping(); err != nil {
				status = http.StatusInternalServerError
				return err.Error()
			}
			return "OK"
		})(),`, "") + `
	}

	return info, status
}
`
			return
		},
	}
}
