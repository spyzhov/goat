package mysql

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "mysql",
		Name:    "MySQL connection",
		Package: "github.com/go-sql-driver/mysql",

		Environments: []*templates.Environment{
			{Name: "MySQLConnect", Type: "string", Env: "MYSQL_CONNECTION", Default: "root:password@tcp(localhost:3306)/database?parseTime=true"},
		},
		Properties: []*templates.Property{
			{Name: "MySQL", Type: "*sql.DB"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/go-sql-driver/mysql", Alias: "_", Version: "v1.4.1"},
			{Name: "database/sql"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.setDataBaseMySQL(); err != nil {
		logger.Panic("cannot connect to MySQL", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// MySQL connect
func (app *Application) setDataBaseMySQL() (err error) {
	app.Logger.Debug("MySQL connect")
	app.MySQL, err = sql.Open("mysql", app.Config.MySQLConnect)
	return
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers: func(*templates.Config) (s string) {
			s = `
	defer app.Closer(app.MySQL, "MySQL connection")`
			return
		},

		Templates: templates.BlankFunctionMap,
	}
}
