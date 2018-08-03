package action

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/spyzhov/goat/console"
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/babex"
	"github.com/spyzhov/goat/templates/http"
	"github.com/spyzhov/goat/templates/migrations"
	"github.com/spyzhov/goat/templates/postgres"
	"github.com/spyzhov/goat/templates/prometheus"
	"github.com/spyzhov/goat/templates/rmq_consumer"
	"github.com/spyzhov/goat/templates/rmq_publisher"
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
}

type Config struct {
	Babex           bool
	Postgres        bool
	Migrations      bool
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

func New(c *cli.Context) *Action {
	var (
		err    error
		gopath = os.Getenv("GOPATH") + "/src/"
		a      = &Action{
			Console: console.New(),
		}
	)
	// region Debug
	if c.Bool("debug") {
		a.log = func(format string, v ...interface{}) {
			log.Printf(format+"\n", v...)
		}
	} else {
		a.log = func(format string, v ...interface{}) {}
	}
	// endregion
	// region Path
	if a.Path, err = filepath.Abs(c.String("path")); err != nil {
		log.Fatal(err)
	}
	a.log("found path: %s", a.Path)
	if !a.Console.PromptY("Project path [%s]?", a.Path) {
		if a.Path, err = a.Console.Scanln("Enter correct path: "); err != nil {
			log.Fatal(err)
		}
		if a.Path, err = filepath.Abs(a.Path); err != nil {
			log.Fatal(err)
		}
		if a.Path == "" {
			log.Fatal("Project path was not set")
		}
		// TODO: validate
	}
	a.log("use path: %s", a.Path)
	// endregion
	// region Name
	a.Name = filepath.Base(a.Path)
	a.log("found name: %s", a.Name)
	if !a.Console.PromptY("Project name [%s]?", a.Name) {
		if a.Name, err = a.Console.Scanln("Enter correct name: "); err != nil {
			log.Fatal(err)
		}
		if a.Name == "" {
			log.Fatal("Project name was not set")
		}
		// TODO: validate
	}
	a.log("use name: %s", a.Name)
	// endregion
	// region Repository name
	if strings.HasPrefix(a.Path, gopath) {
		a.Repo = strings.TrimPrefix(a.Path, gopath)
	}
	a.log("found repository name: %s", a.Repo)
	if a.Repo != "" && !a.Console.PromptY("Repository name [%s]?", a.Repo) {
		if a.Repo, err = a.Console.Scanln("Enter correct repository name: "); err != nil {
			log.Fatal(err)
		}
		if a.Repo == "" {
			log.Fatal("Repository name was not set")
		}
		// TODO: validate
	}
	a.log("use repository name: %s", a.Repo)
	// endregion
	return a
}

func (a *Action) Invoke() (err error) {
	a.log("start to generate")
	a.Config = a.getConfig()
	context := &Context{
		Env:            a.getEnv(),
		Name:           a.Name,
		Repo:           a.Repo,
		Repos:          a.getLibs(),
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
		file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	return
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

func (a *Action) getConfig() *Config {
	a.log("generate config")
	conf := &Config{}
	conf.Postgres = a.Console.Prompt("Use Postgres connection (github.com/go-pg)?")
	if conf.Postgres {
		conf.Migrations = a.Console.PromptY("With migrations (github.com/go-pg/migrations)?")
	}
	conf.Http = a.Console.Prompt("Use HTTP server (github.com/labstack/echo)?")
	conf.Prometheus = a.Console.Prompt("Use Prometheus (github.com/prometheus/client_golang)?")
	conf.Babex = a.Console.Prompt("Use Babex-service (github.com/matroskin13/babex)?")
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
	if a.Config.Migrations {
		a.log("env: get migrations")
		env = append(env, migrations.Env...)
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
	tpld := fmt.Sprintf("\t%%-%ds %%-%ds `env:\"%%s\" envDefault:%%s`", length[0], length[1])
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
	if a.Config.Migrations {
		a.log("props: get migrations")
		props = append(props, migrations.Props...)
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
	if a.Config.Migrations {
		a.log("props: get migrations")
		props = append(props, migrations.Props...)
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
	if a.Config.Migrations {
		a.log("lib: get migrations")
		lib = append(lib, migrations.Libs...)
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
		parts[l.Name] = l.Name
	}
	return "\t\"" + joinMap(parts, "\"\n\t\"") + "\""
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
	if a.Config.Migrations {
		a.log("runners: add migrations")
		parts = append(parts, migrations.TemplateRunFunction)
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
	if a.Config.Migrations {
		a.log("setter: add migrations")
		parts = append(parts, migrations.TemplateSetter)
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
	if a.Config.Migrations {
		a.log("setter-function: add migrations")
		parts = append(parts, migrations.TemplateSetterFunction)
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
	if a.Config.Migrations {
		a.log("models: add migrations")
		for name, data := range migrations.Models {
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
	if a.Config.Migrations {
		a.log("files: add migrations")
		for name, data := range migrations.Templates {
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
