package fasthttp

import (
	"github.com/spyzhov/goat/templates"
)

func New() *templates.Template {
	return &templates.Template{
		ID:      "fasthttp",
		Name:    "FastHTTP server",
		Package: "github.com/valyala/fasthttp",

		Environments: []*templates.Environment{},
		Properties: []*templates.Property{
			{Name: "Http", Type: "*fasthttp.Server", Default: `&fasthttp.Server{
			DisableKeepalive: true,
			LogAllErrors:     true,
			Logger:           &Logger{logger: logger.Named("fasthttp")},
		}`},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/valyala/fasthttp", Version: "v1.2.0"},
		},
		Models: map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction: func(config *templates.Config) (s string) {
			s = `	// Run FastHTTP server
	if err := app.RunFastHttp(); err != nil {
		app.Logger.Panic("FastHTTP Server start error", zap.Error(err))
	}`
			return
		},
		TemplateClosers: templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			prom := config.IsEnabled("prometheus")
			strings = map[string]string{
				"app/http.go": `package app

import (
	"encoding/json"
	"fmt"` + templates.Str(prom, `
	"github.com/prometheus/client_golang/prometheus/promhttp"`, "") + `
	"github.com/valyala/fasthttp"` + templates.Str(prom, `
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"{{.Repo}}/metrics"`, "") + `
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// Start HTTP server
func (app *Application) RunFastHttp() error {
	app.Logger.Info("Start FastHTTP server")
	app.Http.Handler = app.fastHttpHandler

	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()
		app.Logger.Info("fasthttp server started on [::]:" + strconv.Itoa(app.Config.Port))

		app.WaitGroup.Add(1)
		go func() {
			defer app.WaitGroup.Done()
			app.Error <- app.Http.ListenAndServe(":" + strconv.Itoa(app.Config.Port))
			app.Logger.Debug("fasthttp server stops serve")
		}()

		<-app.Ctx.Done()

		if err := app.Http.Shutdown(); err != nil {
			app.Logger.Error("fasthttp server close error", zap.Error(err))
		}
		app.Logger.Debug("fasthttp stops")
	}()
	return nil
}

// Base FastHTTP handler: resolve routes
func (app *Application) fastHttpHandler(ctx *fasthttp.RequestCtx) {
	if app.Config.Debug {
		defer func(start time.Time) {
			app.Logger.Debug("request",
				zap.Duration("duration", time.Now().Sub(start)),
				zap.ByteString("method", ctx.Method()),
				zap.ByteString("path", ctx.Path()))
		}(time.Now())
	}` + templates.Str(prom, `
	metrics.RequestInc()`, "") + `
	switch string(ctx.Path()) {
	case "/favicon.ico":
		setStatusCode(ctx, http.StatusNoContent)
	case "/healthcheck":
		app.healthCheck(ctx)` + templates.Str(prom, `
	case "/metrics":
		handler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
		handler(ctx)
		metrics.RequestStatusInc(string(ctx.Method()), ctx.Response.StatusCode())`, "") + `
	default:
		app.defaultHttpHandler(ctx)
	}
}

func (app *Application) defaultHttpHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	if _, err := fmt.Fprint(ctx, "Not implemented"); err != nil {
		app.Logger.Warn("error on write response", zap.Error(err))
	}
	setStatusCode(ctx, http.StatusNotImplemented)
}

func (app *Application) healthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	info, status := app.healthCheck()
	err := json.NewEncoder(ctx).Encode(info)
	if err != nil {
		app.Logger.Warn("error on write response", zap.Error(err))
		setStatusCode(ctx, http.StatusInternalServerError)
		return
	}
	setStatusCode(ctx, status)
}

func setStatusCode(ctx *fasthttp.RequestCtx, status int) {` + templates.Str(prom, `
	metrics.RequestStatusInc(string(ctx.Method()), status)`, "") + `
	ctx.SetStatusCode(status)
}
`,
			}
			return
		},
	}
}
