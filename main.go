package main

import (
	"github.com/spyzhov/goat/action"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Goat"
	app.HelpName = "goat"
	app.Usage = "golang: application template"
	app.Description = "Goat is golang application template generator"
	app.Version = "0.1.5"
	app.Authors = []cli.Author{
		{Name: "S.Pyzhov", Email: "turin.tomsk@gmail.com"},
	}
	app.Before = func(context *cli.Context) error {
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
