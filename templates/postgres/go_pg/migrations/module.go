package migrations

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "go-postgres/migrations",
		Name:         "Postgres migrations",
		Package:      "github.com/go-pg/migrations",
		Dependencies: []string{"go-postgres"},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries: []*templates.Library{
			{Name: "{{.Repo}}/migrations"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.migratePostgres(); err != nil {
		app.Logger.Error("cannot migrate on Postgres", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// PG migrations up
func (app *Application) migratePostgres() error {
	app.Logger.Debug("PG migrate")
	migrations.Postgres(app.Postgres, app.Logger)
	return nil
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers:     templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"migrations/postgres.go": `package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

func Postgres(db *pg.DB, logger *zap.Logger) {
	migrations.SetTableName("_{{.Name}}_migrations")
	_, _, err := migrations.Run(db, "init")
	if err != nil {
		logger.Debug("Postgres migrations: initialized")
	} else {
		logger.Info("Postgres migrations: initialize")
	}

	oldVersion, newVersion, err := migrations.Run(db, "up")

	if err != nil {
		logger.Panic("Error on run migration", zap.Error(err))
	}

	if newVersion != oldVersion {
		logger.Info("Postgres migrations: migrated", zap.Int64("old", oldVersion), zap.Int64("new", newVersion))
	} else {
		logger.Info("Postgres migrations: version", zap.Int64("current", oldVersion))
	}
}
`,
				"migrations/01_init.go": `package migrations

import (
	"fmt"
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("init migration...")

		_, err := db.Exec("SELECT 1")

		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("SELECT 1")

		return err
	})
}
`,
			}
			return
		},
	}
}
