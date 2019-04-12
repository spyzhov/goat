package action

import (
	"errors"
	"github.com/spyzhov/goat/console"
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/aws_lambda"
	"github.com/spyzhov/goat/templates/clickhouse"
	chMigrations "github.com/spyzhov/goat/templates/clickhouse/migrations"
	"github.com/spyzhov/goat/templates/mysql"
	myMigrations "github.com/spyzhov/goat/templates/mysql/migrations"
	"github.com/spyzhov/goat/templates/postgres"
	pgMigrations "github.com/spyzhov/goat/templates/postgres/migrations"
	"github.com/spyzhov/goat/templates/prometheus"
	"github.com/spyzhov/goat/templates/rmq_consumer"
	"github.com/spyzhov/goat/templates/rmq_publisher"
	"github.com/spyzhov/goat/templates/webserver"
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
	Context *cli.Context
	Config  *templates.Config
}

type Context struct {
	Env            string
	Name           string
	Repo           string
	Repos          string
	Runners        string
	Closers        string
	Setter         string
	SetterFunction string
	Props          string
	PropsValue     string
	Models         string
	DepLibs        string
	MdCode         string
}

type TemplateResponse string

func New(c *cli.Context) *Action {
	var err error
	a := &Action{
		Context: c,
		Console: console.New(),
	}
	if _, err = a.Console.Print(console.Wrap("Select environment:", console.OkGreen)); err != nil {
		log.Fatal(err)
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
	if _, err = a.Console.Print(console.Wrap("Select packages:", console.OkGreen)); err != nil {
		log.Fatal(err)
	}
	a.Config = a.getConfig()
	context := &Context{
		Env:            a.getEnvironments(),
		Name:           a.Name,
		Repo:           a.Repo,
		Repos:          render("repos", a.getLibraries(), a),
		Runners:        a.getTemplateRunFunctions(),
		Closers:        a.getTemplateClosers(),
		Setter:         a.getTemplateSetters(),
		SetterFunction: a.getTemplateSetterFunction(),
		Props:          a.getProperties(),
		PropsValue:     a.getPropertiesValue(),
		Models:         a.getModels(),
		DepLibs:        a.getDepLibraries(),
		MdCode:         "```",
	}
	files := make(map[string]*template.Template)
	for name, content := range a.getTemplateFiles() {
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
			log.Fatal(err1)
		}
	}

	if _, err = a.Console.Print(console.Wrap("Done!", console.OkGreen)); err != nil {
		log.Fatal(err)
	}
	if _, err = a.Console.Print("Don't forget to run: %s", console.Wrap("dep ensure", console.Bold)); err != nil {
		log.Fatal(err)
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

func (t *TemplateResponse) Write(p []byte) (n int, err error) {
	*t += TemplateResponse(p)
	return len(*t), nil
}

func (a *Action) getConfig() *templates.Config {
	a.log("generate config")
	conf := &templates.Config{
		Templates: []*templates.Template{
			postgres.New(),
			pgMigrations.New(),
			mysql.New(),
			myMigrations.New(),
			clickhouse.New(),
			chMigrations.New(),
			webserver.New(),
			prometheus.New(),
			rmq_consumer.New(),
			rmq_publisher.New(),
			aws_lambda.New(),
		},
		Install: []*templates.Template{templates.New()},
	}
	conf.Init(a.Console)
	return conf
}

func (a *Action) getEnvironments() string {
	a.log("env: start")
	return a.Config.Environments().String()
}

func (a *Action) getProperties() string {
	a.log("props: start")
	return a.Config.Properties().String()
}

func (a *Action) getPropertiesValue() string {
	a.log("props-value: start")
	return a.Config.Properties().Values()
}

func (a *Action) getLibraries() string {
	a.log("lib: start")
	libs := a.Config.Libraries()
	if a.Config.IsEnabled("fasthttp") { // FIXME: chose correct way to do it
		index := find(func(i interface{}) bool {
			return i.(*templates.Library).Name == "github.com/prometheus/client_golang/prometheus/promhttp"
		}, libs.Interface()...)
		if index != -1 {
			libs = append(libs[:index], libs[index+1:]...)
		}
	}
	return libs.String()
}

func (a *Action) getDepLibraries() string {
	a.log("dep-lib: start")
	return a.Config.Libraries().Dep()
}

func (a *Action) getTemplateRunFunctions() string {
	a.log("runners: start")
	return a.Config.TemplateRunFunctions().String()
}

func (a *Action) getTemplateClosers() string {
	a.log("closers: start")
	return a.Config.TemplateClosers().String()
}

func (a *Action) getTemplateSetters() string {
	a.log("setter: start")
	return a.Config.TemplateSetters().String()
}

func (a *Action) getTemplateSetterFunction() string {
	a.log("setter-function: start")
	return a.Config.TemplateSetterFunctions().String()
}

func (a *Action) getModels() string {
	a.log("models: start")
	return a.Config.Models().String()
}

func (a *Action) getTemplateFiles() map[string]string {
	a.log("files: start")
	return a.Config.TemplateFiles()
}

func find(check func(interface{}) bool, args ...interface{}) int {
	for i, s := range args {
		if check(s) {
			return i
		}
	}
	return -1
}
