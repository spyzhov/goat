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
	if err = app.setDataBasePostgres(&app.Postgres); err != nil {
		logger.Panic("cannot connect to Postgres", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// PG connect
func (a *Application) setDataBasePostgres(db **pg.DB) error {
	a.Logger.Debug("PG connect", zap.String("connect", a.Config.PgConnect))

	options, err := pg.ParseURL(a.Config.PgConnect)
	if err != nil {
		return err
	}

	*db = pg.Connect(options)
	return nil
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,

		Templates: templates.BlankFunctionMap,
	}
}
