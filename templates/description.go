package templates

type Environment struct {
	Name    string
	Type    string
	Env     string
	Default string
}
type Property struct {
	Name    string
	Type    string
	Default string
}
type Library struct {
	Name string
}

var Env = []Environment{
	{Name: "Level", Type: "string", Env: "LOG_LEVEL", Default: "info"},
	{Name: "Debug", Type: "bool", Env: "DEBUG"},
}
var Props = []Property{
	{Name: "Logger", Type: "*zap.Logger", Default: "logger"},
	{Name: "Config", Type: "*Config", Default: "config"},
	{Name: "Error", Type: "chan error", Default: "make(chan error)"},
}
var Libs = []Library{
	{Name: "go.uber.org/zap"},
}
var Models map[string]string

var Templates = map[string]string{
	"main.go": `package main

import (
	"os"
	"go.uber.org/zap"
	"{{.Repo}}/app"
	"{{.Repo}}/signals"
)

func main() {
	var (
		application *app.Application
		err         error
	)
	if application, err = app.NewApp(); err != nil {
		application.Logger.Fatal("service init error", zap.Error(err))
	}

{{.Runners}}

	select {
	case err = <-application.Error:
		application.Logger.Fatal("service crashed", zap.Error(err))
	case sig := <-signals.WaitExit():
		application.Logger.Info("service stop", zap.Stringer("signal", sig))
		os.Exit(0)
	}
}

`,
	"app/config.go": `package app

import (
	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

type Config struct {
{{.Env}}
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

func NewLogger(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	atom := zap.NewAtomicLevel()
	err := atom.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	cfg.Level = atom

	return cfg.Build()
}

`,
	"app/app.go": `package app

import (
{{.Repos}}
)

type Application struct {
{{.Props}}
}
{{.Models}}

func NewApp() (*Application, error) {
	config, err := NewConfig()
	logger, _ := NewLogger(config.Level)
	if err != nil {
		logger.Fatal("cannot parse config", zap.Error(err))
		return nil, err
	}
	logger.Debug("debug mode on")

	app := &Application{
{{.PropsValue}}
	}
{{.Setter}}

	return app, nil
}
{{.SetterFunction}}
`,
	"signals/signals.go": `package signals

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitExit waits while user don't press Ctrl+C
func WaitExit() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return sigs
}

`,
	"Dockerfile": `FROM golang:1.10 as builder

WORKDIR /go/src/{{.Repo}}

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{.Name}} .

FROM scratch

EXPOSE 4000

WORKDIR /root/
COPY --from=builder /go/src/{{.Repo}} .
CMD ["./{{.Name}}"]
`,
	"README.md": `# About

TODO

## Config [ENV]
{{.MdCode}}go
type Config struct {
{{.Env}}
}
{{.MdCode}}
`,
}
