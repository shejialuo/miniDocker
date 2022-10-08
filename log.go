package main

import (
	"fmt"
	"io"
	"miniDocker/container"
	"os"

	"github.com/sirupsen/logrus"
)

func logContainer(containerName string) {
	logPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	logFilePath := logPath + container.ContainerLogFile
	file, err := os.Open(logFilePath)
	defer file.Close()
	if err != nil {
		logrus.Errorf("log container open file %s error %v", logFilePath, err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		logrus.Errorf("log container read file %s error %v", logFilePath, err)
		return
	}
	fmt.Fprint(os.Stdout, string(content))
}
