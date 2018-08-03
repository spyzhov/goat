package http

import "github.com/spyzhov/goat/templates"

var Env []templates.Environment
var Props = []templates.Property{
	{Name: "Echo", Type: "*echo.Echo", Default: "echo.New()"},
}
var Libs = []templates.Library{
	{Name: "github.com/labstack/echo"},
}
var Models = map[string]string{}

var TemplateSetter = ""
var TemplateSetterFunction = ""
var TemplateRunFunction = `	// Run HTTP server
	if err = application.RunHttp(); err != nil {
		application.Logger.Fatal("Echo start error", zap.Error(err))
	}`
var Templates = map[string]string{
	"app/http.go": `package app

import (
	"net/http"
	"github.com/labstack/echo"
)

func (a *Application) RunHttp() error {
	a.Echo.GET("/", a.httpGetMain)

	go func() {
		a.Error <- a.Echo.Start(":4000")
	}()
	return nil
}

func (a *Application) httpGetMain(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}
`,
}
