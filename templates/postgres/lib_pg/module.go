package lib_pg

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "lib-postgres",
		Name:         "Postgres connection",
		Package:      "github.com/lib/pq",
		Dependencies: []string{"postgres"},

		Environments: []*templates.Environment{},
		Properties: []*templates.Property{
			{Name: "Postgres", Type: "*sql.DB"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/lib/pq", Alias: "_", Version: "^1.0.0"},
			{Name: "database/sql"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.setDataBasePostgres(); err != nil {
		logger.Panic("cannot connect to Postgres", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// PG connect
func (app *Application) setDataBasePostgres() (err error) {
	app.Logger.Debug("PG connect")
	app.Postgres, err = sql.Open("postgres", app.Config.PgConnect)
	if err == nil {
		err = app.Postgres.Ping()
	}
	return
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers: func(*templates.Config) (s string) {
			s = `
	defer func() {
		if app.Postgres != nil {
			app.Closer("Postgres connection", app.Postgres)
		}
	}()`
			return
		},

		Templates: templates.BlankFunctionMap,
	}
}
