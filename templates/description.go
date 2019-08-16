package templates

func New() *Template {
	return &Template{
		ID:      "main",
		Name:    "Main",
		Package: "",

		Environments: []*Environment{
			{Name: "Level", Type: "string", Env: "LOG_LEVEL", Default: "info"},
			{Name: "Debug", Type: "bool", Env: "DEBUG"},
		},
		Properties: []*Property{
			{Name: "Logger", Type: "*zap.Logger", Default: "logger"},
			{Name: "Config", Type: "*Config", Default: "config"},
			{Name: "Error", Type: "chan error", Default: "make(chan error, math.MaxUint8)"},
			{Name: "Ctx", Type: "context.Context"},
			{Name: "ctxCancel", Type: "context.CancelFunc"},
			{Name: "WaitGroup", Type: "*sync.WaitGroup", Default: "new(sync.WaitGroup)"},
		},
		Libraries: []*Library{
			{Name: "go.uber.org/zap", Version: "v1.10.0"},
			{Name: "{{.Repo}}/signals"},
			{Name: "math"},
			{Name: "io"},
			{Name: "context"},
			{Name: "sync"},
			{Name: "time"},
		},
		Models: map[string]string{},

		TemplateSetter:         BlankFunction,
		TemplateSetterFunction: BlankFunction,
		TemplateRunFunction:    BlankFunction,
		TemplateClosers:        BlankFunction,

		Templates: func(config *Config) (strings map[string]string) {
			strings = map[string]string{
				"main.go": `package main

import (
	"{{.Repo}}/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		panic(err)
	}
	defer application.Close()
	application.Start()
}
`,
				"app/config.go": `package app

import (` + Str(config.IsEnabled("console"), `
	"flag"`, "") + `
	"github.com/caarlos0/env/v6"
)

type Config struct {
{{.Env}}
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)` + Str(config.IsEnabled("console"), `
	if err = env.Parse(cfg); err != nil {
		return nil, err
	}
	return FlagsConfig(cfg)`, `
	return cfg, env.Parse(cfg)`) + `
}
` + Str(config.IsEnabled("console"), `
func FlagsConfig(cfg *Config) (*Config, error) {
{{.Flags}}
	flag.Parse()
{{.FlagsEnv}}
	return cfg, nil
}
`, ""),
				"app/logger.go": `package app

import (
	"fmt"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(level string) (logger *zap.Logger, err error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	atom := zap.NewAtomicLevel()
	err = atom.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	cfg.Level = atom

	logger, err = cfg.Build()

	zap.ReplaceGlobals(logger)

	return logger, err
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Println(v ...interface{}) {
	l.logger.Warn(fmt.Sprint(v...))
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
	if err != nil {
		return nil, err
	}
	logger, err := NewLogger(config.Level)
	if err != nil {
		return nil, err
	}
	logger.Debug("debug mode on")

	app := &Application{
{{.PropsValue}}
	}
	app.Ctx, app.ctxCancel = context.WithCancel(context.Background())
	defer func() {
		if err != nil {
			app.Close()
		}
	}()
{{.Setter}}

	return app, nil
}

func (app *Application) Close() {
	app.Logger.Debug("Application stops")
{{.Closers}}
}

func (app *Application) Start() {
	defer app.Stop()
{{.Runners -}}
{{- if eq .ServiceType "lambda"}}
	lambda.Start(app.Lambda)
{{- else}}

	select {
	case err := <-app.Error:
		app.Logger.Panic("service crashed", zap.Error(err))
	case <-app.Ctx.Done():
		app.Logger.Error("service stops via context")
	case sig := <-signals.WaitExit():
		app.Logger.Info("service stop", zap.Stringer("signal", sig))
	} {{- end}}
}

func (app *Application) Stop() {
	app.Logger.Info("service stopping...")
	app.ctxCancel()
{{- if eq .ServiceType "lambda"}}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
{{- else}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
{{- end}}

	go func() {
		defer cancel()
		app.WaitGroup.Wait()
	}()

	<-ctx.Done()

	if ctx.Err() != context.Canceled {
		app.Logger.Panic("service stopped with timeout")
	} else {
		app.Logger.Info("service stopped with success")
	}
}

func (app *Application) Closer(closer io.Closer, scope string) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			app.Logger.Warn("closer error", zap.String("scope", scope), zap.Error(err))
		}
	}
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
# -0 means no compression. Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM golang:1.12 AS builder
ENV GO111MODULE=on
# build via packr hard way https://github.com/gobuffalo/packr#building-a-binary-the-hard-way
RUN go get -u github.com/gobuffalo/packr/v2/packr2
WORKDIR /go/src/{{.Repo}}
COPY . .
RUN go mod download && \
	packr2 && \
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/{{.Name}} . && \
	packr2 clean

FROM busybox:latest
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
				"go.mod": `module {{.Repo}}

go 1.12

require (
	{{.GoMods}}
)
`,
			}
			return
		},
	}
}
