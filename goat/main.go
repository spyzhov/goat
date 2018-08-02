package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/spyzhov/goat/action"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Goat"
	app.Description = "golang: application template"
	app.Version = "0.0.1"
	app.Before = func(context *cli.Context) error {
		fmt.Println("before")
		return nil
	}
	app.Action = func(c *cli.Context) error {
		a := action.New(c)

		return a.Invoke()
	}
	app.EnableBashCompletion = true
	app.BashComplete = func(c *cli.Context) {
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug mode",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "Path to the output directory",
			Value: getPath(),
		},
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
