package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"miniDocker/container"
	"os"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
)

func stopContainer(containerName string) {
	pid, err := GetContainerPidByName(containerName)
	if err != nil {
		logrus.Errorf("get contaienr pid by name %s error %v", containerName, err)
		return
	}
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		logrus.Errorf("convert pid from string to int error %v", err)
		return
	}
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		logrus.Errorf("stop container %s error %v", containerName, err)
		return
	}
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		logrus.Errorf("get container %s info error %v", containerName, err)
		return
	}
	containerInfo.Status = container.Stop
	containerInfo.Pid = " "
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("Json marshal %s error %v", containerName, err)
		return
	}
	configPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := configPath + container.ConfigName
	if err := ioutil.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		logrus.Errorf("Write file %s error", configFilePath, err)
	}
}

func getContainerInfoByName(containerName string) (*container.ContainerInfo, error) {
	configPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := configPath + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		logrus.Errorf("read file %s error %v", configFilePath, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		logrus.Errorf("getContainerInfoByName unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}

func removeContainer(containerName string) {
	containerInfo, err := getContainerInfoByName(containerName)

	if err != nil {
		logrus.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	if containerInfo.Status != container.Stop {
		logrus.Errorf("Couldn't remove running container")
		return
	}
	configPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(configPath); err != nil {
		logrus.Errorf("Remove file %s error %v", configPath, err)
		return
	}
}
