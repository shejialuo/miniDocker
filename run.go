package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"miniDocker/cgroups"
	"miniDocker/cgroups/subsystems"
	"miniDocker/container"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Run is the interface for the `run` command.
func Run(tty bool, commandArray []string, res *subsystems.ResourceConfig, volume string,
	containerName string, imageName string) {

	containerID := randStringBytes(10)
	if containerName == "" {
		containerName = containerID
	}

	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	// We should record the container information into persistent storage
	containerName, err := recordContainerInfo(parent.Process.Pid, commandArray, containerName, containerID, volume)
	if err != nil {
		logrus.Errorf("Record container info error %v", err)
		return
	}

	// Initialize the Cgroup manager

	// TODO: fix the error
	// Now, if we run the daemon container, our parent will exit immediately,
	// However, there will be child process in the cgroup we created.
	// Thus we can not delete the cgroup.
	cgroupManager := cgroups.NewCgroupManager("miniDocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(commandArray, writePipe)
	if tty {
		parent.Wait()
		deleteContainerInfo(containerName)
		container.DeleteWorkSpace(volume, containerName)
	}
}

// Use pipe to send the message to the child
func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

// Generate random id for each container
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Record container information in the `/var/run/miniDocker/
// <container name>/config`.
func recordContainerInfo(containerPID int, commandArray []string, containerName, id, volume string) (string, error) {

	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")

	containerInfo := &container.ContainerInfo{
		ID:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.Running,
		Name:        containerName,
		Volume:      volume,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	recordPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(recordPath, 0622); err != nil {
		logrus.Errorf("Mkdir error %s error %v", recordPath, err)
		return "", err
	}
	filePath := recordPath + "/" + container.ConfigName
	file, err := os.Create(filePath)
	if err != nil {
		logrus.Errorf("create file %s error %v", filePath, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		logrus.Errorf("file write string error %v", err)
		return "", err
	}
	return containerName, nil
}

func deleteContainerInfo(containerID string) {
	recordPath := fmt.Sprintf(container.DefaultInfoLocation, containerID)
	if err := os.RemoveAll(recordPath); err != nil {
		logrus.Errorf("Remove dir %s error %v", recordPath, err)
	}
}
