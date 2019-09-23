package templates

import (
	"fmt"
	"github.com/spyzhov/goat/console"
	"sort"
	"strings"
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
	Conflicts    []string
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
	Name        string
	Type        string
	Env         string
	Default     string
	Flag        string
	Description string
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

func BlankFunction(_ *Config) string {
	return ""
}

func BlankFunctionMap(_ *Config) map[string]string {
	return map[string]string{}
}

//region Template
func (t *Template) Prompt() string {
	verb := "Use"
	if t.Select != nil {
		verb = "Select"
	}
	if t.Package != "" {
		return fmt.Sprintf("%s %s (%s)?", verb, t.Name, console.Wrap(t.Package, console.Underline))
	}
	return fmt.Sprintf("%s %s?", verb, t.Name)
}

func (t *Template) Variants() (result []interface{}) {
	result = make([]interface{}, len(t.Select))
	for i, t := range t.Select {
		result[i] = t.Prompt()
	}
	return result
}

//endregion
//region Config
func (c *Config) Init(console *console.Console) {
	for _, tpl := range c.Templates {
		if c.canPrompt(tpl) {
			if tpl.Select == nil {
				if console.Prompt(tpl.Prompt()) {
					c.Install = append(c.Install, tpl)
				}
			} else {
				if n := console.Select(tpl.Prompt(), tpl.Variants()...); n != 0 {
					c.Install = append(c.Install, tpl)
					c.Install = append(c.Install, tpl.Select[n-1])
				}
			}
		}
	}
}

func (c *Config) canPrompt(tpl *Template) bool {
	for _, tID := range tpl.Dependencies {
		if !c.IsEnabled(tID) {
			return false
		}
	}
	for _, tID := range tpl.Conflicts {
		if c.IsEnabled(tID) {
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

func (env Environments) Flags() string {
	parts := make([]string, 0, len(env))
	tpl := `	var %s = flag.%s("%s", cfg.%s, "%s")`
	for _, e := range env {
		parts = append(parts, fmt.Sprintf(tpl, e.FlagVar(), e.FlagType(), e.FlagName(), e.Name, e.FlagDescription()))
	}
	return join(parts, "\n")
}

func (env Environments) CobraFlags() string {
	parts := make([]string, 0, len(env))
	tpl := `	cmd.PersistentFlags().%sVarP(&cfg.%s, "%s", "", cfg.%s, "%s")`
	for _, e := range env {
		if e.Env != "LOG_LEVEL" {
			parts = append(parts, fmt.Sprintf(tpl, e.FlagType(), e.Name, e.FlagName(), e.Name, e.FlagDescription()))
		}
	}
	return join(parts, "\n")
}

func (env Environments) FlagsEnv() string {
	parts := make([]string, 0, len(env))
	tpl := `	cfg.%s = *%s`
	for _, e := range env {
		parts = append(parts, fmt.Sprintf(tpl, e.Name, e.FlagVar()))
	}
	return join(parts, "\n")
}

//endregion
//region Environment

func (e Environment) FlagName() string {
	result := e.Flag
	if result == "" {
		result = strings.ToLower(e.Env)
		result = strings.ReplaceAll(result, "_", "-")
	}
	return result
}

func (e Environment) FlagType() string {
	return strings.ToUpper(e.Type[0:1]) + e.Type[1:]
}

func (e Environment) FlagVar() string {
	return strings.ToLower(e.Name[0:1]) + e.Name[1:]
}

func (e Environment) FlagDescription() string {
	if e.Description == "" {
		return "flag for ENV:" + e.Env
	}
	return e.Description
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
func (libs Libraries) Interface() []interface{} {
	result := make([]interface{}, len(libs))
	for i, l := range libs {
		result[i] = l
	}
	return result
}

func (libs Libraries) GoMod() string {
	parts := make([]string, 0, len(libs))
	for _, l := range libs {
		if l.Version != "" {
			repo := l.Repo
			if repo == "" {
				repo = l.Name
			}
			parts = append(parts, repo+" "+l.Version)
		}
	}
	return join(parts, "\n\t")
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

func Str(isSuccess bool, success, not string) string {
	if isSuccess {
		return success
	}
	return not
}

//endregion
