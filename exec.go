package main

import (
	"encoding/json"
	"fmt"
	"miniDocker/container"

	// Here, in order to include the C code we
	// need to import this package
	_ "miniDocker/nsenter"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	ENV_EXEC_PID = "mydocker_pid"
	ENV_EXEC_CMD = "mydocker_cmd"
)

// We just create a new process to exec. It's an ugly way
// for calling C code. Really
func ExecContainer(containerName string, commandArray []string) {
	pid, err := GetContainerPidByName(containerName)
	if err != nil {
		logrus.Errorf("exec container getContainerPidByName %s error %v", containerName, err)
		return
	}
	cmdString := strings.Join(commandArray, " ")
	logrus.Infof("container pid %s", pid)
	logrus.Infof("command %s", cmdString)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdString)

	if err := cmd.Run(); err != nil {
		logrus.Errorf("exec container %s error %v", containerName, err)
	}
}

func GetContainerPidByName(containerName string) (string, error) {
	configPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := configPath + container.ConfigName
	contentBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}
