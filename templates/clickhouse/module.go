package clickhouse

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "clickhouse",
		Name:    "ClickHouse connection",
		Package: "github.com/kshvakov/clickhouse",

		Environments: []*templates.Environment{
			{Name: "ClickHouseDSN", Type: "string", Env: "CLICKHOUSE_CONNECTION", Default: "tcp://127.0.0.1:9000?username=&password=&database=default&read_timeout=10&write_timeout=10&debug=true"},
		},
		Properties: []*templates.Property{
			{Name: "ClickHouse", Type: "*sql.DB"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/kshvakov/clickhouse", Branch: "master"},
			{Name: "database/sql"},
			{Name: "fmt"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.setClickHouse(); err != nil {
		logger.Panic("cannot connect to ClickHouse", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// Connect to ClickHouse
func (app *Application) setClickHouse() (err error) {
	app.Logger.Debug("Connect to ClickHouse")
	app.ClickHouse, err = sql.Open("clickhouse", app.Config.ClickHouseDSN)
	if err != nil {
		return
	}
	if err = app.ClickHouse.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			err = fmt.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return
	}
	return
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers: func(*templates.Config) (s string) {
			s = `
	defer func() {
		if app.ClickHouse != nil {
			if err := app.ClickHouse.Close(); err != nil {
				app.Logger.Warn("error on ClickHouse close", zap.Error(err))
			}
		}
	}()`
			return
		},

		Templates: templates.BlankFunctionMap,
	}
}
