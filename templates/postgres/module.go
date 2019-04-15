package postgres

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "postgres",
		Name:    "Postgres connection",
		Package: "github.com/go-pg",

		Environments: []*templates.Environment{
			{Name: "PgConnect", Type: "string", Env: "POSTGRES_CONNECTION", Default: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"},
		},
		Properties: []*templates.Property{
			{Name: "Postgres", Type: "*pg.DB"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/go-pg/pg", Version: "^7.1.0"},
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
func (app *Application) setDataBasePostgres() error {
	app.Logger.Debug("PG connect", zap.String("connect", app.Config.PgConnect))

	options, err := pg.ParseURL(app.Config.PgConnect)
	if err != nil {
		return err
	}

	app.Postgres = pg.Connect(options)
	return nil
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers: func(*templates.Config) (s string) {
			s = `
	defer func() {
		if app.Postgres != nil {
			app.closer("Postgres connection", app.Postgres)
		}
	}()`
			return
		},

		Templates: templates.BlankFunctionMap,
	}
}
