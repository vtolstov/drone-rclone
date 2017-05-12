package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version string

func main() {
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "rclone plugin"
	app.Usage = "rclone plugin"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "flags",
			Usage:  "flags to pass",
			EnvVar: "PLUGIN_FLAGS",
		},
		cli.StringFlag{
			Name:   "action",
			Usage:  "action to run",
			EnvVar: "PLUGIN_ACTION",
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "upload files from source folder",
			EnvVar: "PLUGIN_SOURCE",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "upload files to target folder",
			EnvVar: "PLUGIN_TARGET",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {

	plugin := &Plugin{
		Flags:  c.StringSlice("flags"),
		Action: c.String("action"),
		Source: c.String("source"),
		Target: c.String("target"),
	}

	return plugin.Exec()
}
