package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func commitContainer(imageName string) {
	mntPath := "/root/mnt/"
	imageTar := "/root/" + imageName + ".tar"
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		logrus.Errorf("Tar folder %s error %v", mntPath, err)
	}
}
