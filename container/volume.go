package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// Create a AUFS filesystem as container root workspace
func NewWorkSpace(volume, imageName, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName, imageName)
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs, containerName)
			logrus.Infof("NewWorkSpace volume urls %q", volumeURLs)
		} else {
			logrus.Infof("Volume parameter input is not correct.")
		}
	}
}

// Decompression tar image
func CreateReadOnlyLayer(imageName string) error {
	unTarFolderUrl := RootPath + "/" + imageName + "/"
	imageUrl := RootPath + "/" + imageName + ".tar"
	exist, err := PathExists(unTarFolderUrl)
	if err != nil {
		logrus.Infof("Fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
		return err
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			logrus.Errorf("Mkdir %s error %v", unTarFolderUrl, err)
			return err
		}

		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			logrus.Errorf("Untar dir %s error %v", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerPath, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		logrus.Infof("Mkdir write layer dir %s error. %v", writeURL, err)
	}
}

func MountVolume(volumeURLs []string, containerName string) error {
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		logrus.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
	}
	containerUrl := volumeURLs[1]
	mntURL := fmt.Sprintf(MntPath, containerName)
	containerVolumeURL := mntURL + "/" + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		logrus.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}
	dirs := "dirs=" + parentUrl
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL).CombinedOutput()
	if err != nil {
		logrus.Errorf("Mount volume failed. %v", err)
		return err
	}
	return nil
}

func CreateMountPoint(containerName, imageName string) error {
	mntUrl := fmt.Sprintf(MntPath, containerName)
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		logrus.Errorf("Mkdir mountpoint dir %s error. %v", mntUrl, err)
		return err
	}
	tmpWriteLayer := fmt.Sprintf(WriteLayerPath, containerName)
	tmpImageLocation := RootPath + "/" + imageName
	mntURL := fmt.Sprintf(MntPath, containerName)
	dirs := "dirs=" + tmpWriteLayer + ":" + tmpImageLocation
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL).CombinedOutput()
	if err != nil {
		logrus.Errorf("Run command for creating mount point failed %v", err)
		return err
	}
	return nil
}

// Delete the AUFS filesystem while container exit
func DeleteWorkSpace(volume, containerName string) {
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(volumeURLs, containerName)
		} else {
			DeleteMountPoint(containerName)
		}
	} else {
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) error {
	mntURL := fmt.Sprintf(MntPath, containerName)
	_, err := exec.Command("umount", mntURL).CombinedOutput()
	if err != nil {
		logrus.Errorf("Unmount %s error %v", mntURL, err)
		return err
	}
	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}

func DeleteMountPointWithVolume(volumeURLs []string, containerName string) error {
	mntURL := fmt.Sprintf(MntPath, containerName)
	containerUrl := mntURL + "/" + volumeURLs[1]
	if _, err := exec.Command("umount", containerUrl).CombinedOutput(); err != nil {
		logrus.Errorf("Umount volume %s failed. %v", containerUrl, err)
		return err
	}

	if _, err := exec.Command("umount", mntURL).CombinedOutput(); err != nil {
		logrus.Errorf("Umount mountpoint %s failed. %v", mntURL, err)
		return err
	}

	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("Remove mountpoint dir %s error %v", mntURL, err)
	}

	return nil
}

func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerPath, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		logrus.Infof("Remove writeLayer dir %s error %v", writeURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
