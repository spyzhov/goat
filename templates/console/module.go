package console

import (
	"github.com/spyzhov/goat/templates"
	"github.com/spyzhov/goat/templates/console/blank"
	"github.com/spyzhov/goat/templates/console/cobra"
)

func New() *templates.Template {
	return &templates.Template{
		ID:   "console",
		Name: "Console",
		Select: []*templates.Template{
			blank.New(),
			cobra.New(),
		},
		Conflicts: []string{"aws_lambda", "webserver"},

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries:    []*templates.Library{},
		Models:       map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction: func(config *templates.Config) (s string) {
			s = `	// Run Action
	if err := app.RunAction(); err != nil {
		app.Logger.Panic("Action start error", zap.Error(err))
	}`
			return
		},
		TemplateClosers: templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/console.go": `package app

import (
	"go.uber.org/zap"
	"time"
)

// Start of action
func (app *Application) RunAction() error {
	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()
		defer func(start time.Time) {
			app.Logger.Debug("Action: done", zap.Duration("duration", time.Since(start)))
		}(time.Now())
		app.Logger.Debug("Action: start")
		if err := app.action(); err != nil {
			app.Logger.Error("Action: error", zap.Error(err))
			app.Error <- err
		} else {
			app.ctxCancel()
		}
	}()
	return nil
}
`,
			}
			return
		},
	}
}
