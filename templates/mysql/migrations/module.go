package migrations

import "github.com/spyzhov/goat/templates"

var Env []templates.Environment
var Props []templates.Property
var Libs = []templates.Library{
	{Name: "{{.Repo}}/migrations"},
}
var Models = map[string]string{}

var TemplateSetter = `
	if err = app.migrateMySQL(); err != nil {
		logger.Fatal("cannot migrate on MySQL", zap.Error(err))
		return nil, err
	}`
var TemplateSetterFunction = `
// MySQL migrations up
func (a *Application) migrateMySQL() error {
	a.Logger.Debug("mySQL migrate")
	migrations.MySQL(a.Postgres, a.Logger)
	return nil
}`
var TemplateRunFunction = ""
var Templates = map[string]string{
	"migrations/mysql.go": `package migrations

import (
	"database/sql"
	"github.com/gobuffalo/packr"
	"github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

func MySQL(db *sql.DB, logger *zap.Logger) {
	migrate.SetTable("_{{.Name}}_migrations")
	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./mysql"),
	}
	logger.Debug("MySQL migrations: start")
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		return err
	}

	rows, err := migrate.GetMigrationRecords(db, "mysql")
	if err != nil {
		return err
	}
	cnt := len(rows)
	last := ""
	if cnt > 0 {
		last = rows[cnt-1].Id
	}

	logger.Info("MySQL migrations: migrated", zap.Int("count", n), zap.String("current", last))
	return nil
}
`,
	"migrations/mysql/1-init.sql": `-- +migrate Up
SELECT NOW();

-- +migrate Down
SELECT NOW();
`,
}
