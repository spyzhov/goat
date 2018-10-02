package mysql

import "github.com/spyzhov/goat/templates"

var Env = []templates.Environment{
	{Name: "MySQLConnect", Type: "string", Env: "MYSQL_CONNECTION", Default: "root:password@tcp(localhost:3306)/database?parseTime=true"},
	{Name: "MySQLIdleConnections", Type: "int", Env: "MYSQL_IDLE_CONNECTIONS", Default: "2"},
	{Name: "MySQLMaxConnections", Type: "int", Env: "MYSQL_MAX_CONNECTIONS", Default: "2"},
}
var Props = []templates.Property{
	{Name: "MySQL", Type: "*sql.DB"},
}
var Libs = []templates.Library{
	{Name: "github.com/go-sql-driver/mysql", Alias: "_"},
	{Name: "database/sql"},
}
var Models = map[string]string{}

var TemplateSetter = `
	if err = app.setDataBaseMySQL(&app.MySQL); err != nil {
		logger.Fatal("cannot connect to MySQL", zap.Error(err))
		return nil, err
	}`
var TemplateSetterFunction = `
// MySQL connect
func (a *Application) setDataBaseMySQL(db **sql.DB) (err error) {
	a.Logger.Debug("MySQL connect", zap.String("connect", a.Config.MySQLConnect))
	*db, err = sql.Open("mysql", a.Config.MySQLConnection)
	if err != nil {
		return
	}
	a.MySQL.SetMaxOpenConns(a.Config.MySQLMaxConnections)
	a.MySQL.SetMaxIdleConns(a.Config.MySQLIdleConnections)
	return
}`
var TemplateRunFunction = ""
var Templates = map[string]string{}
