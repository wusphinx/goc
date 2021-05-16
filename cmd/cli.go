package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	appName    string
	moduleName string
	goVersion  string
	parentDir  string
	debugMode  bool
)

func App() {
	app := cli.NewApp()
	app.Name = "goc-cli"
	app.Version = "v0.0.1"

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version: %s\n", c.App.Version)
	}

	app.Usage = "An interactive cli to create go app"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "new",
			Aliases:     []string{"n"},
			Usage:       "create new go app",
			Destination: &appName,
			Required:    true,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"D"},
			Usage:       "Enable debug mode",
			Destination: &debugMode,
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		if debugMode {
			logrus.SetLevel(logrus.DebugLevel)
		}

		goVersion, err = promptGoVersion()
		if err != nil {
			return err
		}

		moduleName, err = promptModule()
		if err != nil {
			return err
		}

		if err := Generator(); err != nil {
			return err
		}

		fmt.Printf("Congratulation！👏  Go ahead with %s 👊\n", appName)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
