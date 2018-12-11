package action

import (
	"errors"
	"fmt"
	"github.com/spyzhov/goat/console"
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/babex"
	"github.com/spyzhov/goat/templates/http"
	"github.com/spyzhov/goat/templates/mysql"
	myMigrations "github.com/spyzhov/goat/templates/mysql/migrations"
	"github.com/spyzhov/goat/templates/postgres"
	pgMigrations "github.com/spyzhov/goat/templates/postgres/migrations"
	"github.com/spyzhov/goat/templates/prometheus"
	"github.com/spyzhov/goat/templates/rmq_consumer"
	"github.com/spyzhov/goat/templates/rmq_publisher"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Action struct {
	log     func(format string, v ...interface{})
	Console *console.Console
	Name    string
	Repo    string
	Path    string
	Config  *Config
	Context *cli.Context
}

type Config struct {
	Babex           bool
	Postgres        bool
	PgMigrations    bool
	MySQL           bool
	MyMigrations    bool
	RabbitConsumer  bool
	RabbitPublisher bool
	Prometheus      bool
	Http            bool
}

type Context struct {
	Env            string
	Name           string
	Repo           string
	Repos          string
	Runners        string
	Setter         string
	SetterFunction string
	Props          string
	PropsValue     string
	Models         string
	MdCode         string
}

type TemplateResponse string

func New(c *cli.Context) *Action {
	var err error
	a := &Action{
		Context: c,
		Console: console.New(),
	}
	// region Debug
	if err = a.setDebug(); err != nil {
		log.Fatal(err)
	}
	// endregion
	// region Path
	if err = a.setPath(); err != nil {
		log.Fatal(err)
	}
	// endregion
	// region Name
	if err = a.setName(); err != nil {
		log.Fatal(err)
	}
	// endregion
	// region Repository name
	if err = a.setRepo(); err != nil {
		log.Fatal(err)
	}
	// endregion
	return a
}

func (a *Action) setDebug() (err error) {
	if a.Context.Bool("debug") {
		a.log = func(format string, v ...interface{}) {
			log.Printf(format+"\n", v...)
		}
	} else {
		a.log = func(format string, v ...interface{}) {}
	}
	return
}

func (a *Action) setPath() (err error) {
	if a.Path, err = filepath.Abs(a.Context.String("path")); err != nil {
		return
	}
	a.log("found path: %s", a.Path)
	if !a.Console.PromptY("Project path [%s]?", a.Path) {
		if a.Path, err = a.Console.Scanln("Enter correct path: "); err != nil {
			return
		}
		if a.Path, err = filepath.Abs(a.Path); err != nil {
			return
		}
		if a.Path == "" {
			return errors.New("project path was not set")
		}
		// TODO: validate
	}
	a.log("use path: %s", a.Path)
	return
}

func (a *Action) setName() (err error) {
	a.Name = filepath.Base(a.Path)
	a.log("found name: %s", a.Name)
	if !a.Console.PromptY("Project name [%s]?", a.Name) {
		if a.Name, err = a.Console.Scanln("Enter correct name: "); err != nil {
			return
		}
		if a.Name == "" {
			return errors.New("project name was not set")
		}
		// TODO: validate
	}
	a.log("use name: %s", a.Name)
	return
}

func (a *Action) setRepo() (err error) {
	gopath := os.Getenv("GOPATH") + "/src/"
	if strings.HasPrefix(a.Path, gopath) {
		a.Repo = strings.TrimPrefix(a.Path, gopath)
	}
	a.log("found repository name: %s", a.Repo)
	if a.Repo != "" && !a.Console.PromptY("Repository name [%s]?", a.Repo) {
		if a.Repo, err = a.Console.Scanln("Enter correct repository name: "); err != nil {
			return
		}
		if a.Repo == "" {
			return errors.New("repository name was not set")
		}
		// TODO: validate
	}
	a.log("use repository name: %s", a.Repo)
	return
}

func (a *Action) Invoke() (err error) {
	a.log("start to generate")
	a.Config = a.getConfig()
	context := &Context{
		Env:            a.getEnv(),
		Name:           a.Name,
		Repo:           a.Repo,
		Repos:          render("repos", a.getLibs(), a),
		Runners:        a.getRunners(),
		Setter:         a.getSetter(),
		SetterFunction: a.getSetterFunction(),
		Props:          a.getProps(),
		PropsValue:     a.getPropsValue(),
		Models:         a.getModels(),
		MdCode:         "```",
	}
	files := make(map[string]*template.Template)
	for name, content := range a.getFiles() {
		files[name] = template.Must(template.New(name).Parse(content))
	}

	for name, tpl := range files {
		fileName := a.Path + "/" + name
		a.log("Process file: %s", name)
		if err := os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		err = tpl.Execute(file, context)
		err1 := file.Close()
		if err != nil {
			log.Fatal(err)
		}
		if err1 != nil {
			log.Fatal(err)
		}
	}

	return
}

func render(name, tpl string, obj interface{}) string {
	result := new(TemplateResponse)
	if err := template.Must(template.New(name).Parse(tpl)).Execute(result, obj); err != nil {
		log.Fatal(err)
	}
	return string(*result)
}

func join(parts []string, sep string) (result string) {
	for _, part := range parts {
		if part != "" {
			result += sep + part
		}
	}
	l := len(sep)
	if len(result) >= l {
		return result[l:]
	}
	return
}

func joinMap(parts map[string]string, sep string) (result string) {
	for _, part := range parts {
		if part != "" {
			result += sep + part
		}
	}
	l := len(sep)
	if len(result) >= l {
		return result[l:]
	}
	return
}

func (t *TemplateResponse) Write(p []byte) (n int, err error) {
	*t += TemplateResponse(p)
	return len(*t), nil
}

func (a *Action) getConfig() *Config {
	a.log("generate config")
	conf := &Config{}
	conf.Postgres = a.Console.Prompt("Use Postgres connection (github.com/go-pg)?")
	if conf.Postgres {
		conf.PgMigrations = a.Console.PromptY("With postgres migrations (github.com/go-pg/migrations)?")
	}
	conf.MySQL = a.Console.Prompt("Use MySQL connection (github.com/go-sql-driver/mysql)?")
	if conf.MySQL {
		conf.MyMigrations = a.Console.PromptY("With MySQL migrations (github.com/rubenv/sql-migrate)?")
	}
	conf.Http = a.Console.Prompt("Use HTTP server (het/http)?")
	conf.Prometheus = a.Console.Prompt("Use Prometheus (github.com/prometheus/client_golang)?")
	//fixme conf.Babex = a.Console.Prompt("Use Babex-service (github.com/matroskin13/babex)?")
	conf.RabbitConsumer = a.Console.Prompt("Use RMQ-consumers (github.com/streadway/amqp)?")
	conf.RabbitPublisher = a.Console.Prompt("Use RMQ-publishers (github.com/streadway/amqp)?")
	return conf
}

func (a *Action) getEnv() string {
	a.log("env: start")
	var (
		parts  []string
		length [2]int
		env    []templates.Environment
		e      templates.Environment
		l      int
	)
	env = append(env, templates.Env...)
	if a.Config.Babex {
		a.log("env: get babex")
		env = append(env, babex.Env...)
	}
	if a.Config.Postgres {
		a.log("env: get postgres")
		env = append(env, postgres.Env...)
	}
	if a.Config.PgMigrations {
		a.log("env: get pgMigrations")
		env = append(env, pgMigrations.Env...)
	}
	if a.Config.MySQL {
		a.log("env: get mysql")
		env = append(env, mysql.Env...)
	}
	if a.Config.MyMigrations {
		a.log("env: get myMigrations")
		env = append(env, myMigrations.Env...)
	}
	if a.Config.RabbitConsumer {
		a.log("env: get rmq_consumer")
		env = append(env, rmq_consumer.Env...)
	}
	if a.Config.RabbitPublisher {
		a.log("env: get rmq_publisher")
		env = append(env, rmq_publisher.Env...)
	}
	if a.Config.Prometheus {
		a.log("env: get prometheus")
		env = append(env, prometheus.Env...)
	}
	if a.Config.Http {
		a.log("env: get http")
		env = append(env, http.Env...)
	}

	a.log("env: calculate length")
	for _, e = range env {
		l = len(e.Name)
		if l > length[0] {
			length[0] = l
		}
		l = len(e.Type)
		if l > length[1] {
			length[1] = l
		}
	}
	tpl := fmt.Sprintf("\t%%-%ds %%-%ds `env:\"%%s\"`", length[0], length[1])
	tpld := fmt.Sprintf("\t%%-%ds %%-%ds `env:\"%%s\" envDefault:\"%%s\"`", length[0], length[1])
	for _, e = range env {
		if e.Default != "" {
			parts = append(parts, fmt.Sprintf(tpld, e.Name, e.Type, e.Env, e.Default))
		} else {
			parts = append(parts, fmt.Sprintf(tpl, e.Name, e.Type, e.Env))
		}
	}
	return join(parts, "\n")
}

func (a *Action) getProps() string {
	a.log("props: start")
	var (
		parts  = make(map[string]string)
		length int
		props  []templates.Property
		p      templates.Property
		l      int
	)
	props = append(props, templates.Props...)
	if a.Config.Babex {
		a.log("props: get babex")
		props = append(props, babex.Props...)
	}
	if a.Config.Postgres {
		a.log("props: get postgres")
		props = append(props, postgres.Props...)
	}
	if a.Config.PgMigrations {
		a.log("props: get pgMigrations")
		props = append(props, pgMigrations.Props...)
	}
	if a.Config.MySQL {
		a.log("props: get mysql")
		props = append(props, mysql.Props...)
	}
	if a.Config.MyMigrations {
		a.log("props: get myMigrations")
		props = append(props, myMigrations.Props...)
	}
	if a.Config.RabbitConsumer {
		a.log("props: get rmq_consumer")
		props = append(props, rmq_consumer.Props...)
	}
	if a.Config.RabbitPublisher {
		a.log("props: get rmq_publisher")
		props = append(props, rmq_publisher.Props...)
	}
	if a.Config.Prometheus {
		a.log("props: get prometheus")
		props = append(props, prometheus.Props...)
	}
	if a.Config.Http {
		a.log("props: get prometheus")
		props = append(props, http.Props...)
	}

	a.log("props: calculate length")
	for _, p = range props {
		l = len(p.Name)
		if l > length {
			length = l
		}
	}
	tpl := fmt.Sprintf("\t%%-%ds %%s", length)
	for _, p = range props {
		parts[p.Name] = fmt.Sprintf(tpl, p.Name, p.Type)
	}
	return joinMap(parts, "\n")
}

func (a *Action) getPropsValue() string {
	a.log("props-value: start")
	var (
		parts  = make(map[string]string)
		length int
		props  []templates.Property
		p      templates.Property
		l      int
	)
	props = append(props, templates.Props...)
	if a.Config.Babex {
		a.log("props: get babex")
		props = append(props, babex.Props...)
	}
	if a.Config.Postgres {
		a.log("props: get postgres")
		props = append(props, postgres.Props...)
	}
	if a.Config.PgMigrations {
		a.log("props: get pgMigrations")
		props = append(props, pgMigrations.Props...)
	}
	if a.Config.MySQL {
		a.log("props: get mysql")
		props = append(props, mysql.Props...)
	}
	if a.Config.MyMigrations {
		a.log("props: get myMigrations")
		props = append(props, myMigrations.Props...)
	}
	if a.Config.RabbitConsumer {
		a.log("props: get rmq_consumer")
		props = append(props, rmq_consumer.Props...)
	}
	if a.Config.Prometheus {
		a.log("props: get prometheus")
		props = append(props, prometheus.Props...)
	}
	if a.Config.Http {
		a.log("props: get http")
		props = append(props, http.Props...)
	}

	a.log("props: calculate length")
	for _, p = range props {
		l = len(p.Name)
		if l > length && p.Default != "" {
			length = l
		}
	}
	tpl := fmt.Sprintf("\t\t%%-%ds %%s", length+1)
	for _, p = range props {
		if p.Default != "" {
			parts[p.Name] = fmt.Sprintf(tpl, p.Name+":", p.Default+",")
		}
	}
	return joinMap(parts, "\n")
}

func (a *Action) getLibs() string {
	a.log("lib: start")
	var (
		parts = map[string]string{}
		lib   []templates.Library
		l     templates.Library
	)
	lib = append(lib, templates.Libs...)
	if a.Config.Babex {
		a.log("lib: get babex")
		lib = append(lib, babex.Libs...)
	}
	if a.Config.Postgres {
		a.log("lib: get postgres")
		lib = append(lib, postgres.Libs...)
	}
	if a.Config.PgMigrations {
		a.log("lib: get pgMigrations")
		lib = append(lib, pgMigrations.Libs...)
	}
	if a.Config.MySQL {
		a.log("lib: get mysql")
		lib = append(lib, mysql.Libs...)
	}
	if a.Config.MyMigrations {
		a.log("lib: get myMigrations")
		lib = append(lib, myMigrations.Libs...)
	}
	if a.Config.RabbitConsumer {
		a.log("lib: get rmq_consumer")
		lib = append(lib, rmq_consumer.Libs...)
	}
	if a.Config.RabbitPublisher {
		a.log("lib: get rmq_publisher")
		lib = append(lib, rmq_publisher.Libs...)
	}
	if a.Config.Prometheus {
		a.log("lib: get prometheus")
		lib = append(lib, prometheus.Libs...)
	}
	if a.Config.Http {
		a.log("lib: get http")
		lib = append(lib, http.Libs...)
	}

	for _, l = range lib {
		if l.Alias != "" {
			parts[l.Name] = l.Alias + ` "` + l.Name + `"`
		} else {
			parts[l.Name] = `"` + l.Name + `"`
		}
	}
	return "\t" + joinMap(parts, "\n\t")
}

func (a *Action) getRunners() string {
	a.log("runners: start")
	var (
		parts []string
	)
	if a.Config.Babex {
		a.log("runners: add babex")
		parts = append(parts, babex.TemplateRunFunction)
	}
	if a.Config.Postgres {
		a.log("runners: add postgres")
		parts = append(parts, postgres.TemplateRunFunction)
	}
	if a.Config.PgMigrations {
		a.log("runners: add pgMigrations")
		parts = append(parts, pgMigrations.TemplateRunFunction)
	}
	if a.Config.MySQL {
		a.log("runners: add mysql")
		parts = append(parts, mysql.TemplateRunFunction)
	}
	if a.Config.MyMigrations {
		a.log("runners: add myMigrations")
		parts = append(parts, myMigrations.TemplateRunFunction)
	}
	if a.Config.RabbitConsumer {
		a.log("runners: add rmq_consumer")
		parts = append(parts, rmq_consumer.TemplateRunFunction)
	}
	if a.Config.RabbitPublisher {
		a.log("runners: add rmq_publisher")
		parts = append(parts, rmq_publisher.TemplateRunFunction)
	}
	if a.Config.Prometheus {
		a.log("runners: add prometheus")
		parts = append(parts, prometheus.TemplateRunFunction)
	}
	if a.Config.Http && !a.Config.Prometheus {
		a.log("runners: add http")
		parts = append(parts, http.TemplateRunFunction)
	}

	return join(parts, "\n")
}

func (a *Action) getSetter() string {
	a.log("setter: start")
	var (
		parts []string
	)
	if a.Config.Babex {
		a.log("setter: add babex")
		parts = append(parts, babex.TemplateSetter)
	}
	if a.Config.Postgres {
		a.log("setter: add postgres")
		parts = append(parts, postgres.TemplateSetter)
	}
	if a.Config.PgMigrations {
		a.log("setter: add pgMigrations")
		parts = append(parts, pgMigrations.TemplateSetter)
	}
	if a.Config.MySQL {
		a.log("setter: add mysql")
		parts = append(parts, mysql.TemplateSetter)
	}
	if a.Config.MyMigrations {
		a.log("setter: add myMigrations")
		parts = append(parts, myMigrations.TemplateSetter)
	}
	if a.Config.RabbitConsumer {
		a.log("setter: add rmq_consumer")
		parts = append(parts, rmq_consumer.TemplateSetter)
	}
	if a.Config.RabbitPublisher {
		a.log("setter: add rmq_publisher")
		parts = append(parts, rmq_publisher.TemplateSetter)
	}
	if a.Config.Prometheus {
		a.log("setter: add prometheus")
		parts = append(parts, prometheus.TemplateSetter)
	}
	if a.Config.Http {
		a.log("setter: add http")
		parts = append(parts, http.TemplateSetter)
	}

	return join(parts, "\n")
}

func (a *Action) getSetterFunction() string {
	a.log("setter-function: start")
	var (
		parts []string
	)
	if a.Config.Babex {
		a.log("setter-function: add babex")
		parts = append(parts, babex.TemplateSetterFunction)
	}
	if a.Config.Postgres {
		a.log("setter-function: add postgres")
		parts = append(parts, postgres.TemplateSetterFunction)
	}
	if a.Config.PgMigrations {
		a.log("setter-function: add pgMigrations")
		parts = append(parts, pgMigrations.TemplateSetterFunction)
	}
	if a.Config.MySQL {
		a.log("setter-function: add mysql")
		parts = append(parts, mysql.TemplateSetterFunction)
	}
	if a.Config.MyMigrations {
		a.log("setter-function: add myMigrations")
		parts = append(parts, myMigrations.TemplateSetterFunction)
	}
	if a.Config.RabbitConsumer {
		a.log("setter-function: add rmq_consumer")
		parts = append(parts, rmq_consumer.TemplateSetterFunction)
	}
	if a.Config.RabbitPublisher {
		a.log("setter-function: add rmq_publisher")
		parts = append(parts, rmq_publisher.TemplateSetterFunction)
	}
	if a.Config.Prometheus {
		a.log("setter-function: add prometheus")
		parts = append(parts, prometheus.TemplateSetterFunction)
	}
	if a.Config.Http {
		a.log("setter-function: add http")
		parts = append(parts, http.TemplateSetterFunction)
	}

	return join(parts, "\n")
}

func (a *Action) getModels() string {
	a.log("models: start")
	var (
		parts = make(map[string]string)
	)
	for name, data := range templates.Models {
		parts[name] = data
	}
	if a.Config.Babex {
		a.log("models: add babex")
		for name, data := range babex.Models {
			parts[name] = data
		}
	}
	if a.Config.Postgres {
		a.log("models: add postgres")
		for name, data := range postgres.Models {
			parts[name] = data
		}
	}
	if a.Config.PgMigrations {
		a.log("models: add pgMigrations")
		for name, data := range pgMigrations.Models {
			parts[name] = data
		}
	}
	if a.Config.MySQL {
		a.log("models: add mysql")
		for name, data := range mysql.Models {
			parts[name] = data
		}
	}
	if a.Config.MyMigrations {
		a.log("models: add myMigrations")
		for name, data := range myMigrations.Models {
			parts[name] = data
		}
	}
	if a.Config.RabbitConsumer {
		a.log("models: add rmq_consumer")
		for name, data := range rmq_consumer.Models {
			parts[name] = data
		}
	}
	if a.Config.RabbitPublisher {
		a.log("models: add rmq_publisher")
		for name, data := range rmq_publisher.Models {
			parts[name] = data
		}
	}
	if a.Config.Prometheus {
		a.log("models: add prometheus")
		for name, data := range prometheus.Models {
			parts[name] = data
		}
	}
	if a.Config.Http {
		a.log("models: add http")
		for name, data := range http.Models {
			parts[name] = data
		}
	}

	return joinMap(parts, "\n")
}

func (a *Action) getFiles() map[string]string {
	a.log("files: start")
	var (
		parts = make(map[string]string)
	)
	for name, data := range templates.Templates {
		parts[name] = data
	}
	if a.Config.Babex {
		a.log("files: add babex")
		for name, data := range babex.Templates {
			parts[name] = data
		}
	}
	if a.Config.Postgres {
		a.log("files: add postgres")
		for name, data := range postgres.Templates {
			parts[name] = data
		}
	}
	if a.Config.PgMigrations {
		a.log("files: add pgMigrations")
		for name, data := range pgMigrations.Templates {
			parts[name] = data
		}
	}
	if a.Config.MySQL {
		a.log("files: add mysql")
		for name, data := range mysql.Templates {
			parts[name] = data
		}
	}
	if a.Config.MyMigrations {
		a.log("files: add myMigrations")
		for name, data := range myMigrations.Templates {
			parts[name] = data
		}
	}
	if a.Config.RabbitConsumer {
		a.log("files: add rmq_consumer")
		for name, data := range rmq_consumer.Templates {
			parts[name] = data
		}
	}
	if a.Config.RabbitPublisher {
		a.log("files: add rmq_publisher")
		for name, data := range rmq_publisher.Templates {
			parts[name] = data
		}
	}
	if a.Config.Prometheus {
		a.log("files: add prometheus")
		for name, data := range prometheus.Templates {
			parts[name] = data
		}
	}
	if a.Config.Http {
		a.log("files: add http")
		for name, data := range http.Templates {
			parts[name] = data
		}
	}

	return parts
}
