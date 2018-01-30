package main

import (
	"github.com/codegangsta/cli"
	"os"
	"github.com/serussell/logxi/v1"
	"github.com/jasoncorbett/goslick/slickconfig"
)

func main() {
	logger := log.New("slick")
	app := cli.NewApp()
	app.Name = "slick"
	app.Version = "5.0.0"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "config, c",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.IsSet("c") {
			logger.Debug("Loading config from %s", c.String("c"))
			slickconfig.Configuration.LoadFromLocation(c.String("c"))
		} else {
			logger.Debug("Loading config from standard locations.")
			slickconfig.Configuration.LoadFromStandardLocations()
		}
		return nil
	}

	app.Commands = []cli.Command{
		InitCommand,
	}

	app.Run(os.Args)
}
