package main

import (
	"encoding/json"
	"fmt"
	"miniDocker/container"
	"os"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
)

// From persistent files to get the container information
func ListContainers() {
	recordPath := fmt.Sprintf(container.DefaultInfoLocation, "")
	recordPath = recordPath[:len(recordPath)-1]
	files, err := os.ReadDir(recordPath)
	if err != nil {
		logrus.Errorf("Read directory %s error %v", recordPath, err)
	}

	var containers []*container.ContainerInfo
	for _, file := range files {
		tmpContainer, err := getContainerInfo(file)
		if err != nil {
			logrus.Errorf("get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}
	if err := w.Flush(); err != nil {
		logrus.Errorf("Flush error %v", err)
		return
	}
}

// Read json file and return the containerInfo pointer
func getContainerInfo(file os.DirEntry) (*container.ContainerInfo, error) {
	containerName := file.Name()
	configFilePath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath = configFilePath + container.ConfigName
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		logrus.Errorf("Read file %s error %v", configFilePath, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		logrus.Errorf("Json unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil

}
