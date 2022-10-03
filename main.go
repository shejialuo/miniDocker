package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `miniDocker is a simple container runtime implementation.
               The purpose of this project is to learn how docker works
							 and how to write a docker by ourselves. Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "miniDocker"
	app.Usage = usage

	// Here, we register the two commands into the
	// global command.
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
