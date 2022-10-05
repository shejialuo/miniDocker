package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	commandArray := readUserCommand()
	if len(commandArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	// MS_NOEXEC: Do not allow programs to be executed from this filesystem
	// MS_NOSUID: Do not honor set-user-ID and set-group-ID bits

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	// We need to mount the `/proc` for the child process
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	path, err := exec.LookPath(commandArray[0])
	if err != nil {
		logrus.Errorf("exec loop path error %v", err)
		return err
	}
	logrus.Infof("Find path %s", path)

	if err := syscall.Exec(path, commandArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	message := string(msg)
	return strings.Split(message, " ")
}
