package mysql

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "mysql",
		Name:    "MySQL connection",
		Package: "github.com/go-sql-driver/mysql",

		Environments: []*templates.Environment{
			{Name: "MySQLConnect", Type: "string", Env: "MYSQL_CONNECTION", Default: "root:password@tcp(localhost:3306)/database?parseTime=true"},
			{Name: "MySQLIdleConnections", Type: "int", Env: "MYSQL_IDLE_CONNECTIONS", Default: "2"},
			{Name: "MySQLMaxConnections", Type: "int", Env: "MYSQL_MAX_CONNECTIONS", Default: "2"},
		},
		Properties: []*templates.Property{
			{Name: "MySQL", Type: "*sql.DB"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/go-sql-driver/mysql", Alias: "_", Version: "^1.4.1"},
			{Name: "database/sql"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.setDataBaseMySQL(&app.MySQL); err != nil {
		logger.Panic("cannot connect to MySQL", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// MySQL connect
func (a *Application) setDataBaseMySQL(db **sql.DB) (err error) {
	a.Logger.Debug("MySQL connect", zap.String("connect", a.Config.MySQLConnect))
	*db, err = sql.Open("mysql", a.Config.MySQLConnect)
	if err != nil {
		return
	}
	a.MySQL.SetMaxOpenConns(a.Config.MySQLMaxConnections)
	a.MySQL.SetMaxIdleConns(a.Config.MySQLIdleConnections)
	return
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,

		Templates: templates.BlankFunctionMap,
	}
}
