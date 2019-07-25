package go_pg

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "go-postgres",
		Name:         "Postgres connection",
		Package:      "github.com/go-pg",
		Dependencies: []string{"postgres"},

		Environments: []*templates.Environment{},
		Properties: []*templates.Property{
			{Name: "Postgres", Type: "*pg.DB"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/go-pg/pg", Version: "v7.1.0"},
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
	app.Logger.Debug("PG connect")

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
	defer app.Closer(app.Postgres, "Postgres connection")`
			return
		},

		Templates: templates.BlankFunctionMap,
	}
}
