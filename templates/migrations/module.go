package migrations

import "github.com/spyzhov/goat/templates"

var Env []templates.Environment
var Props []templates.Property
var Libs = []templates.Library{
	{Name: "github.com/go-pg/migrations"},
}
var Models = map[string]string{}

var TemplateSetter = `
	if err = app.setMigrationsUp(); err != nil {
		logger.Fatal("cannot migrate on Postgres", zap.Error(err))
		return nil, err
	}`
var TemplateSetterFunction = `
// PG migrations up
func (a *Application) setMigrationsUp() error {
	a.Logger.Debug("PG migrate")
	migrations.Init(a.Db, a.Logger)
	return nil
}
`
var TemplateRunFunction = ""
var Templates = map[string]string{
	"migrations/migrations.go": `package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

func Init(db *pg.DB, logger *zap.Logger) {
	migrations.SetTableName("_{{.Name}}_migrations")

	oldVersion, newVersion, err := migrations.Run(db, "up")

	if err != nil {
		logger.Fatal("Error on run migration", zap.Error(err))
	}

	if newVersion != oldVersion {
		logger.Info("Migrations: migrated", zap.Int64("old", oldVersion), zap.Int64("new", newVersion))
	} else {
		logger.Info("Migrations: version", zap.Int64("current", oldVersion))
	}
}
`,
	"migrations/01_init.go":`package migrations

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
