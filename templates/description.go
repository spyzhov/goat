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
	Name  string
	Alias string
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
	"go.uber.org/zap"
	"{{.Repo}}/app"
	"{{.Repo}}/signals"
)

func main() {
	var (
		application *app.Application
		err         error
	)
	if application, err = app.New(); err != nil {
		panic(err)
	}

{{.Runners}}

	select {
	case err = <-application.Error:
		application.Logger.Fatal("service crashed", zap.Error(err))
	case sig := <-signals.WaitExit():
		application.Logger.Info("service stop", zap.Stringer("signal", sig))
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

func New() (*Application, error) {
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
	"Dockerfile": `FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM golang:1.10 AS builder
# build via packr hard way https://github.com/gobuffalo/packr#building-a-binary-the-hard-way
RUN go get -u github.com/gobuffalo/packr/...
WORKDIR /go/src/{{.Repo}}
ADD . .
RUN packr
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/{{.Name}} .
RUN packr clean

FROM scratch
# configurations
EXPOSE 4000
WORKDIR /root
# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# the main program:
COPY --from=builder /go/bin/{{.Name}} ./{{.Name}}
CMD ["./{{.Name}}"]
`,
	"README.md": `# About

TODO

## Config [ENV]
{{.MdCode}}go
package app
type Config struct {
{{.Env}}
}
{{.MdCode}}
`,
}
