package blank

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:   "console_blank",
		Name: "Blank",

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries:    []*templates.Library{},
		Models:       map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction:    templates.BlankFunction,
		TemplateClosers:        templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/console_action.go": `package app

import (
	"errors"
)

// TODO: Implement action
func (app *Application) action() (err error) {
	app.Logger.Info("Action: run")
	return errors.New("not implemented")
}
`,
			}
			return
		},
	}
}
