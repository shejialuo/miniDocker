package container

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// In order to get the status of the detached container
// information, we need a data structure to describe it
type ContainerInfo struct {
	Pid         string `json:"pid"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"createTime"`
	Status      string `json:"status"`
}

const (
	Running             string = "running"
	Stop                string = "stop"
	Exit                string = "exited"
	DefaultInfoLocation string = "/var/run/miniDocker/%s/"
	ConfigName          string = "config.json"
)

// This function creates the new parent process, it should
// be only invoked once. And wait for the child process (container)
// terminate, when the container terminate, parent should terminate
// too.
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {

	// Here, we create a new pipe
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}

	// Here, `/proc/self/exe` means execute itself, for example
	// if we are in `bash`, we could use `/proc/self/exe` to spawn
	// a child bash shell. In this example, we actually
	// execute the process `miniDocker` again, here, we need
	// to call initCommand to initialize the container. In order to
	// pass the arguments here we use pipe for IPC
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Here, we let the readPipe descriptor for the child to be 3.
	cmd.ExtraFiles = []*os.File{readPipe}

	mntPath := "/root/mnt/"
	rootPath := "/root/"
	NewWorkSpace(rootPath, mntPath, volume)
	cmd.Dir = mntPath

	// We return the `writePipe` to the parent process, make it to write.
	// There is no data race, because when child reads if parent
	// doesn't write, it will automatically be blocked.
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkSpace(rootPath string, mntPath string, volume string) {
	CreateReadOnlyLayer(rootPath)
	CreateWriteLayer(rootPath)
	CreateMountPoint(rootPath, mntPath)
	// Here, we need to handle mount the volume to the container
	if volume != "" {
		volumePaths := volumePathExtract(volume)
		length := len(volumePaths)
		if length == 2 && volumePaths[0] != "" && volumePaths[1] != "" {
			MountVolume(rootPath, mntPath, volumePaths)
			logrus.Infof("%q", volumePaths)
		} else {
			logrus.Infof("Volume parameter input is not correct.")
		}
	}
}

// The read only layer should be in the `/root/busybox`, if
// there is no directory, we should unTar the file `/root/busybox.tar`.
// We need to mount it as the read-only layer.
func CreateReadOnlyLayer(rootPath string) {
	busyboxPath := rootPath + "busybox"
	busyboxTarPath := rootPath + "busybox.tar"
	exist, err := PathExists(busyboxPath)
	if err != nil {
		logrus.Infof("Fail to judge whether dir %s exists. %v", busyboxPath, err)
	}
	if !exist {
		if err := os.Mkdir(busyboxPath, 0777); err != nil {
			logrus.Errorf("Mkdir dir %s error. %v", busyboxPath, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarPath, "-C", busyboxPath).CombinedOutput(); err != nil {
			logrus.Errorf("unTar dir %s error %v", busyboxTarPath, err)
		}
	}
}

// We mount the `/root/writeLayer` to the write layer.
func CreateWriteLayer(rootPath string) {
	writePath := rootPath + "writeLayer"
	if err := os.Mkdir(writePath, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error. %v", writePath, err)
	}

}

// What we actually mount is the `mntPath`. And kernel has provided
// us this abstraction
func CreateMountPoint(rootPath string, mntPath string) {
	if err := os.Mkdir(mntPath, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error. %v", mntPath, err)
	}
	dirs := "dirs=" + rootPath + "writeLayer:" + rootPath + "busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
}

// We should mount the volume specified by the user
func MountVolume(rootPath string, mntPath string, volumePaths []string) {
	parentPath, containerPath := volumePaths[0], volumePaths[1]
	// I don't think it's good idea just create a new directory.
	// It is really messy here.
	if err := os.Mkdir(parentPath, 0777); err != nil {
		logrus.Infof("Mkdir parent dir %s error. %v", parentPath, err)
	}
	containerVolumePath := mntPath + containerPath
	if err := os.Mkdir(containerVolumePath, 0777); err != nil {
		logrus.Infof("Mkdir container dir %s error. %v", containerVolumePath, err)
	}

	dirs := "dirs=" + parentPath
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Mount volume failed. %v", err)
	}
}

func DeleteWorkSpace(rootPath string, mntPath string, volume string) {
	DeleteMountPoint(rootPath, mntPath, volume)
	DeleteWriteLayer(rootPath)
}

// We should umount the path and delete the path
// Also, we need to umount the volume
func DeleteMountPoint(rootPath string, mntPath string, volume string) {
	if volume != "" {
		volumePaths := volumePathExtract(volume)
		length := len(volumePaths)
		if length == 2 && volumePaths[0] != "" && volumePaths[1] != "" {
			containerPath := mntPath + volumePaths[1]
			cmd := exec.Command("umount", containerPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				logrus.Errorf("umount volume failed. %v", err)
			}
		}
	}
	cmd := exec.Command("umount", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
	if err := os.RemoveAll(mntPath); err != nil {
		logrus.Errorf("Remove dir %s error %v", mntPath, err)
	}
}

// We should delete the write layer
func DeleteWriteLayer(rootPath string) {
	writePath := rootPath + "writeLayer"
	if err := os.RemoveAll(writePath); err != nil {
		logrus.Errorf("Remove dir %s error %v", writePath, err)
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

func volumePathExtract(volume string) []string {
	volumePaths := strings.Split(volume, ":")
	return volumePaths
}
