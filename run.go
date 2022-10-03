package main

import (
	"miniDocker/container"
	"os"

	"github.com/sirupsen/logrus"
)

// Run is the interface for the `run` command.
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	parent.Wait()
	os.Exit(0)
}
