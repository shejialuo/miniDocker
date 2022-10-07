package main

import (
	"miniDocker/cgroups"
	"miniDocker/cgroups/subsystems"
	"miniDocker/container"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Run is the interface for the `run` command.
func Run(tty bool, commandArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	// Initialize the Cgroup manager
	cgroupManager := cgroups.NewCgroupManager("miniDocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(commandArray, writePipe)
	parent.Wait()
	mntPath := "/root/mnt/"
	rootPath := "/root/"
	container.DeleteWorkSpace(rootPath, mntPath)
}

// Use pipe to send the message to the child
func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
