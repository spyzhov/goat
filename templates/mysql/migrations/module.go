package migrations

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "mysql/migrations",
		Name:         "MySQL migrations",
		Package:      "github.com/rubenv/sql-migrate",
		Dependencies: []string{"mysql"},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries: []*templates.Library{
			{Name: "{{.Repo}}/migrations"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.migrateMySQL(); err != nil {
		app.Logger.Error("cannot migrate on MySQL", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// MySQL migrations up
func (app *Application) migrateMySQL() error {
	app.Logger.Debug("MySQL migrate")
	return migrations.MySQL(app.MySQL, app.Logger)
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers:     templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"migrations/mysql.go": `package migrations

import (
	"database/sql"
	"github.com/gobuffalo/packr/v2"
	"github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
	"strings"
)

func MySQL(db *sql.DB, logger *zap.Logger) error {
	migrate.SetTable("_{{.Name}}_migrations")
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("mysql", "./mysql"),
	}
	logger.Debug("MySQL migrations: start")
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		if strings.HasSuffix(err.Error(), "unknown migration in database") {
			logger.Warn("MySQL migrations: SKIPPED", zap.Error(err))
			return nil
		}
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
			return
		},
	}
}
