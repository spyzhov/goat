package http

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "http",
		Name:    "HTTP server",
		Package: "net/http",

		Environments: []*templates.Environment{},
		Properties: []*templates.Property{
			{Name: "Http", Type: "*http.ServeMux", Default: "http.NewServeMux()"},
		},
		Libraries: []*templates.Library{
			{Name: "net/http"},
		},
		Models: map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction: func(config *templates.Config) (s string) {
			s = `	// Run HTTP server
	if err := app.RunHttp(); err != nil {
		app.Logger.Panic("HTTP Server start error", zap.Error(err))
	}`
			return
		},
		TemplateClosers: templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/http.go": `package app

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// Declare all necessary HTTP methods
func (app *Application) registerRoutes() {
	app.Http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		if _, err := fmt.Fprint(w, "Not implemented"); err != nil {
			app.Logger.Warn("error on write response", zap.Error(err))
		}
	})
	app.Http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		info, status := app.healthCheck()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(info); err != nil {
			app.Logger.Warn("error on write response", zap.Error(err))
		}
	})
	app.Http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(app.Info); err != nil {
			app.Logger.Warn("error on write response", zap.Error(err))
		}
	})
}

// Start HTTP server
func (app *Application) RunHttp() error {
	app.registerRoutes()

	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()
		app.Logger.Info("http server started on [::]:" + strconv.Itoa(app.Config.Port))
		server := &http.Server{
			Addr:    ":" + strconv.Itoa(app.Config.Port),
			Handler: app.Http,
		}
		server.RegisterOnShutdown(app.ctxCancel)

		app.WaitGroup.Add(1)
		go func() {
			defer app.WaitGroup.Done()
			app.Error <- server.ListenAndServe()
			app.Logger.Debug("http server stops serve")
		}()

		<-app.Ctx.Done()

		if err := server.Close(); err != nil {
			app.Logger.Error("http server close error", zap.Error(err))
		}
		app.Logger.Debug("http stops")
	}()
	return nil
}
`,
			}
			return
		},
	}
}
