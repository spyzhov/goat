package webserver

import (
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/webserver/fasthttp"
	"github.com/spyzhov/goat/templates/webserver/http"
)

func New() *templates.Template {
	return &templates.Template{
		ID:   "webserver",
		Name: "WebServer",
		Select: []*templates.Template{
			http.New(),
			fasthttp.New(),
		},

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

		Templates: templates.BlankFunctionMap,
	}
}
