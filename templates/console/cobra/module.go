package cobra

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "cobra",
		Name:    "Cobra",
		Package: "github.com/spf13/cobra",

		Environments: []*templates.Environment{},
		Properties: []*templates.Property{
			{Name: "Command", Type: "*cobra.Command"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/spf13/cobra", Version: "v0.0.5"},
		},
		Models: map[string]string{},

		TemplateSetter:         templates.BlankFunction,
		TemplateSetterFunction: templates.BlankFunction,
		TemplateRunFunction:    templates.BlankFunction,
		TemplateClosers:        templates.BlankFunction,

		Templates: func(config *templates.Config) (strings map[string]string) {
			strings = map[string]string{
				"app/console_action.go": `package app

func (app *Application) action() (err error) {
	return app.Command.Execute()
}
`,
				"app/init_commands.go": `package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (app *Application) InitCommands() error {
	app.Command = &cobra.Command{
		Use:     "{{.Name}}",
		Version: app.Info.Version,
		Short:   "{{.Name}} is a CLI application",
		Long:    "{{.Name}} is a CLI application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO:Command at app/init_commands.go")
		},
	}
	app.Command.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of {{.Name}}",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("{{.Name}} %s -- %s\nBuild at: %s", app.Info.Version, app.Info.Commit, app.Info.Created)
		},
	})
	return nil
}
`,
			}
			return
		},
	}
}
