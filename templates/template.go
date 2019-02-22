package templates

import (
	"fmt"
	"github.com/spyzhov/goat/console"
	"sort"
)

type Config struct {
	Templates []*Template
	Install   []*Template
}

type Template struct {
	ID           string
	Name         string
	Package      string
	Dependencies []string
	Select       []*Template

	Environments []*Environment
	Properties   []*Property
	Libraries    []*Library
	Models       map[string]string

	TemplateSetter         func(*Config) string
	TemplateSetterFunction func(*Config) string
	TemplateRunFunction    func(*Config) string
	TemplateClosers        func(*Config) string
	Templates              func(*Config) map[string]string
}

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
	Name    string
	Alias   string
	Repo    string
	Version string
	Branch  string
}

type (
	Environments            []*Environment
	Properties              []*Property
	Libraries               []*Library
	Models                  []string
	TemplateRunFunctions    []string
	TemplateClosers         []string
	TemplateSetters         []string
	TemplateSetterFunctions []string
	TemplateFiles           map[string]string
)

func BlankFunction(*Config) string {
	return ""
}

func BlankFunctionMap(*Config) map[string]string {
	return map[string]string{}
}

//region Template
func (t *Template) Prompt() string {
	return fmt.Sprintf("Use %s (%s)?", t.Name, t.Package)
}

//endregion
//region Config
func (c *Config) Init(console *console.Console) {
	for _, tpl := range c.Templates {
		if c.canPrompt(tpl) && console.Prompt(tpl.Prompt()) {
			c.Install = append(c.Install, tpl)
		}
	}
}

func (c *Config) canPrompt(tpl *Template) bool {
	for _, tID := range tpl.Dependencies {
		if !c.IsEnabled(tID) {
			return false
		}
	}
	return true
}

func (c *Config) IsEnabled(tID string) bool {
	for _, tpl := range c.Install {
		if tpl.ID == tID {
			return true
		}
	}
	return false
}

func (c *Config) Environments() (result Environments) {
	result = make(Environments, 0)
	setup := make(map[string]bool)
	for _, tpl := range c.Install {
		for _, env := range tpl.Environments {
			if !setup[env.Name] {
				result = append(result, env)
				setup[env.Name] = true
			}
		}
	}
	return
}

func (c *Config) Properties() (result Properties) {
	result = make(Properties, 0)
	setup := make(map[string]bool)
	for _, tpl := range c.Install {
		for _, prop := range tpl.Properties {
			if !setup[prop.Name] {
				result = append(result, prop)
				setup[prop.Name] = true
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return
}

func (c *Config) Libraries() (result Libraries) {
	result = make(Libraries, 0)
	setup := make(map[string]bool)
	for _, tpl := range c.Install {
		for _, lib := range tpl.Libraries {
			if !setup[lib.Name] {
				result = append(result, lib)
				setup[lib.Name] = true
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return
}

func (c *Config) Models() (result Models) {
	result = make(Models, 0)
	setup := make(map[string]bool)
	for _, tpl := range c.Install {
		for name, model := range tpl.Models {
			if !setup[name] {
				result = append(result, model)
				setup[name] = true
			}
		}
	}
	return
}

func (c *Config) TemplateRunFunctions() (result TemplateRunFunctions) {
	result = make(TemplateRunFunctions, 0)
	for _, tpl := range c.Install {
		result = appendIf(result, tpl.TemplateRunFunction(c))
	}
	return
}

func (c *Config) TemplateClosers() (result TemplateClosers) {
	result = make(TemplateClosers, 0)
	for _, tpl := range c.Install {
		result = appendIf(result, tpl.TemplateClosers(c))
	}
	return
}

func (c *Config) TemplateSetters() (result TemplateSetters) {
	result = make(TemplateSetters, 0)
	for _, tpl := range c.Install {
		result = appendIf(result, tpl.TemplateSetter(c))
	}
	return
}

func (c *Config) TemplateSetterFunctions() (result TemplateSetterFunctions) {
	result = make(TemplateSetterFunctions, 0)
	for _, tpl := range c.Install {
		result = appendIf(result, tpl.TemplateSetterFunction(c))
	}
	return
}

func (c *Config) TemplateFiles() (result TemplateFiles) {
	result = make(TemplateFiles)
	for _, tpl := range c.Install {
		for name, value := range tpl.Templates(c) {
			result[name] = value
		}
	}
	return
}

//endregion
//region Environments
func (env Environments) String() string {
	var (
		parts  []string
		length [2]int
	)
	for _, e := range env {
		l := len(e.Name)
		if l > length[0] {
			length[0] = l
		}
		l = len(e.Type)
		if l > length[1] {
			length[1] = l
		}
	}
	tpl := fmt.Sprintf("\t%%-%ds %%-%ds `env:\"%%s\"`", length[0], length[1])
	tpld := fmt.Sprintf("\t%%-%ds %%-%ds `env:\"%%s\" envDefault:\"%%s\"`", length[0], length[1])
	for _, e := range env {
		if e.Default != "" {
			parts = append(parts, fmt.Sprintf(tpld, e.Name, e.Type, e.Env, e.Default))
		} else {
			parts = append(parts, fmt.Sprintf(tpl, e.Name, e.Type, e.Env))
		}
	}
	return join(parts, "\n")
}

//endregion
//region Properties
func (props Properties) String() string {
	length := 0
	parts := make([]string, 0, len(props))
	for _, p := range props {
		l := len(p.Name)
		if l > length {
			length = l
		}
	}
	tpl := fmt.Sprintf("\t%%-%ds %%s", length)
	for _, p := range props {
		parts = append(parts, fmt.Sprintf(tpl, p.Name, p.Type))
	}
	return join(parts, "\n")
}

func (props Properties) Values() string {
	length := 0
	parts := make([]string, 0, len(props))
	for _, p := range props {
		l := len(p.Name)
		if l > length && p.Default != "" {
			length = l
		}
	}
	tpl := fmt.Sprintf("\t\t%%-%ds %%s", length+1)
	for _, p := range props {
		if p.Default != "" {
			parts = append(parts, fmt.Sprintf(tpl, p.Name+":", p.Default+","))
		}
	}
	return join(parts, "\n")
}

//endregion
//region Libraries
func (libs Libraries) String() string {
	parts := make([]string, 0, len(libs))
	for _, l := range libs {
		if l.Alias != "" {
			parts = append(parts, fmt.Sprintf("\t%s \"%s\"", l.Alias, l.Name))
		} else {
			parts = append(parts, fmt.Sprintf("\t\"%s\"", l.Name))
		}
	}
	return join(parts, "\n")
}

func (libs Libraries) Dep() string {
	parts := make([]string, 0, len(libs))
	for _, l := range libs {
		if l.Version != "" || l.Branch != "" {
			repo := l.Repo
			bound, version := "version", l.Version
			if repo == "" {
				repo = l.Name
			}
			if version == "" {
				bound, version = "branch", l.Branch
			}
			parts = append(parts, `[[constraint]]
  name = "`+repo+`"
  `+bound+` = "`+version+`"
`)
		}
	}
	return join(parts, "\n")
}

//endregion
//region Models
func (models Models) String() string {
	return join(models, "\n")
}

//endregion
//region TemplateRunFunctions
func (functions TemplateRunFunctions) String() string {
	return join(functions, "\n")
}

//endregion
//region TemplateClosers
func (functions TemplateClosers) String() string {
	return join(functions, "\n")
}

//endregion
//region TemplateSetters
func (functions TemplateSetters) String() string {
	return join(functions, "\n")
}

//endregion
//region TemplateSetterFunctions
func (functions TemplateSetterFunctions) String() string {
	return join(functions, "\n")
}

//endregion
//region Functions
func join(parts []string, sep string) (result string) {
	for _, part := range parts {
		if part != "" {
			result += sep + part
		}
	}
	l := len(sep)
	if len(result) >= l {
		return result[l:]
	}
	return
}

func appendIf(array []string, value string) []string {
	if value != "" {
		return append(array, value)
	}
	return array
}

//endregion
