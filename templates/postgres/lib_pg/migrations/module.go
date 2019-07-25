package migrations

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "lib-postgres/migrations",
		Name:         "Postgres migrations",
		Package:      "github.com/rubenv/sql-migrate",
		Dependencies: []string{"lib-postgres"},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries: []*templates.Library{
			{Name: "{{.Repo}}/migrations"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.migratePostgres(); err != nil {
		logger.Panic("cannot migrate on Postgres", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// PG migrations up
func (app *Application) migratePostgres() error {
	app.Logger.Debug("PG migrate")
	return migrations.Postgres(app.Postgres, app.Logger)
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers:     templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"migrations/postgres.go": `package migrations

import (
	"database/sql"
	"github.com/gobuffalo/packr/v2"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

func Postgres(db *sql.DB, logger *zap.Logger) error {
	migrate.SetTable("_{{.Name}}_migrations")
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("postgres", "./postgres"),
	}
	logger.Debug("Postgres migrations: start")
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	rows, err := migrate.GetMigrationRecords(db, "postgres")
	if err != nil {
		return err
	}
	cnt := len(rows)
	last := ""
	if cnt > 0 {
		last = rows[cnt-1].Id
	}

	logger.Info("Postgres migrations: migrated", zap.Int("count", n), zap.String("current", last))
	return nil
}
`,
				"migrations/postgres/1-init.sql": `-- +migrate Up
SELECT NOW();

-- +migrate Down
SELECT NOW();
`,
			}
			return
		},
	}
}
