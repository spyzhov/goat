package postgres

import "github.com/spyzhov/goat/templates"

var Env = []templates.Environment{
	{Name: "DbConnect", Type: "string", Env: "POSTGRES_CONNECTION", Default: "\"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable\""},
}
var Props = []templates.Property{
	{Name: "Db", Type: "*pg.DB"},
}
var Libs = []templates.Library{
	{Name: "github.com/go-pg/pg"},
	{Name: "time"},
}
var Models = map[string]string{}

var TemplateSetter = `
	if err = app.setDataBase(app.Db); err != nil {
		logger.Fatal("cannot connect to Postgres", zap.Error(err))
		return nil, err
	}`
var TemplateSetterFunction = `
// PG connect
func (a *Application) setDataBase(db *pg.DB) error {
	a.Logger.Debug("PG connect", zap.String("connect", a.Config.DbConnect))

	options, err := pg.ParseURL(a.Config.DbConnect)
	if err != nil {
		return err
	}

	db = pg.Connect(options)

	if a.Config.Debug {
		a.Logger.Debug("Used debug mode for database queries")
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}
			a.Logger.Debug("postgres query",
				zap.String("query", query),
				zap.Duration("durations", time.Since(event.StartTime)))
		})
	}
	return nil
}
`
var TemplateRunFunction = ""
var Templates = map[string]string{}
