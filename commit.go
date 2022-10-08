package main

import (
	"fmt"
	"miniDocker/container"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func commitContainer(containerName string, imageName string) {
	mntPath := fmt.Sprintf(container.MntPath, containerName)
	mntPath += "/"
	imageTar := container.RootPath + "/" + imageName + ".tar"
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		logrus.Errorf("Tar folder %s error %v", mntPath, err)
	}
}
