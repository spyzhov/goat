package simple

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:   "simple",
		Name: "Simple",

		Environments: []*templates.Environment{},
		Properties:   []*templates.Property{},
		Libraries:    []*templates.Library{},
		Models:       map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction:    templates.BlankFunction,
		TemplateClosers:        templates.BlankFunction,

		Templates: templates.BlankFunctionMap,
	}
}
