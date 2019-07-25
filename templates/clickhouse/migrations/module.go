package migrations

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:           "clickhouse/migrations",
		Name:         "ClickHouse migrations",
		Package:      "github.com/golang-migrate/migrate",
		Dependencies: []string{"clickhouse"},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries: []*templates.Library{
			{Name: "{{.Repo}}/migrations"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.migrateClickHouse(); err != nil {
		logger.Panic("cannot migrate on ClickHouse", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// ClickHouse migrations up
func (app *Application) migrateClickHouse() error {
	app.Logger.Debug("ClickHouse migrate")
	return migrations.ClickHouse(app.ClickHouse, app.Logger)
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers:     templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"migrations/clickhouse.go": `package migrations

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/kshvakov/clickhouse"
	_ "{{.Repo}}/migrations/source/packr"
	"go.uber.org/zap"
)

func ClickHouse(db *sql.DB, logger *zap.Logger) error {
	logger.Debug("ClickHouse migrations: start")
	driver, err := clickhouse.WithInstance(db, &clickhouse.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"packr://migrations/clickhouse",
		"clickhouse", driver)
	if err != nil {
		return err
	}
	oldVersion, _, err := m.Version()
	if err == migrate.ErrNilVersion {
		oldVersion = 0
		err = nil
	} else if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	newVersion, _, err := m.Version()
	if err == migrate.ErrNilVersion {
		oldVersion = 0
		err = nil
	} else if err != nil {
		return err
	}
	logger.Info("ClickHouse migrations: migrated",
		zap.Uint("old_version", oldVersion),
		zap.Uint("new_version", newVersion))
	return nil
}
`,
				"migrations/source/packr/packr.go": `package packr

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/golang-migrate/migrate/v4/source"
	"io"
	nurl "net/url"
	"os"
)

func init() {
	source.Register("packr", &Packr{})
}

type Packr struct {
	url        string
	path       string
	box        *packr.Box
	migrations *source.Migrations
}

func (p *Packr) Open(url string) (source.Driver, error) {
	_, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	// Packr MUST get const string
	box := packr.New("clickhouse", "../../clickhouse")

	nf := &Packr{
		url:        url,
		path:       box.Path,
		box:        box,
		migrations: source.NewMigrations(),
	}

	items := box.List()

	for _, name := range items {
		m, err := source.DefaultParse(name)
		if err != nil {
			continue // ignore files that we can't parse
		}
		if !nf.migrations.Append(m) {
			return nil, fmt.Errorf("unable to parse file %v", name)
		}
	}
	return nf, nil
}

func (p *Packr) Close() error {
	// nothing do to here
	return nil
}

func (p *Packr) First() (version uint, err error) {
	if v, ok := p.migrations.First(); !ok {
		return 0, &os.PathError{Op: "first", Path: p.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (p *Packr) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := p.migrations.Prev(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: p.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (p *Packr) Next(version uint) (nextVersion uint, err error) {
	if v, ok := p.migrations.Next(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: p.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (p *Packr) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := p.migrations.Up(version); ok {
		r, err := p.box.Open(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: p.path, Err: os.ErrNotExist}
}

func (p *Packr) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := p.migrations.Down(version); ok {
		r, err := p.box.Open(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: p.path, Err: os.ErrNotExist}
}
`,
				"migrations/clickhouse/1_init.up.sql":   `SELECT NOW();`,
				"migrations/clickhouse/1_init.down.sql": `SELECT NOW();`,
			}
			return
		},
	}
}
