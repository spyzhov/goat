package postgres

import (
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/postgres/go_pg"
	"github.com/spyzhov/goat/templates/postgres/lib_pg"
)

func New() *templates.Template {
	return &templates.Template{
		ID:   "postgres",
		Name: "Postgres",
		Select: []*templates.Template{
			go_pg.New(),
			lib_pg.New(),
		},

		Environments: []*templates.Environment{
			{Name: "PgConnect", Type: "string", Env: "POSTGRES_CONNECTION", Default: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"},
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
