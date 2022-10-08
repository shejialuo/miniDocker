package container

import (
	"fmt"
	"os"
	"os/exec"
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
	Volume      string `json"volume"`
}

const (
	Running             string = "running"
	Stop                string = "stop"
	Exit                string = "exited"
	DefaultInfoLocation string = "/var/run/miniDocker/%s/"
	ConfigName          string = "config.json"
	ContainerLogFile    string = "container.log"
	RootPath            string = "/root"
	MntPath             string = "/root/mnt/%s"
	WriteLayerPath      string = "/root/writeLayer/%s"
)

// This function creates the new parent process, it should
// be only invoked once. And wait for the child process (container)
// terminate, when the container terminate, parent should terminate
// too.
func NewParentProcess(tty bool, volume string, containerName string, imageName string) (*exec.Cmd, *os.File) {

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
	} else {
		// Now we need to redirect the standard output
		logPath := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err := os.MkdirAll(logPath, 0622); err != nil {
			logrus.Errorf("NewParentProcess mkdir %s error %v", logPath, err)
		}
		logFilePath := logPath + ContainerLogFile
		logFile, err := os.Create(logFilePath)
		if err != nil {
			logrus.Errorf("NewParentProcess create file %s error %v", logFilePath, err)
			return nil, nil
		}
		cmd.Stdout = logFile
	}

	// Here, we let the readPipe descriptor for the child to be 3.
	cmd.ExtraFiles = []*os.File{readPipe}

	// Here, we give each container each workspace
	NewWorkSpace(volume, imageName, containerName)
	cmd.Dir = fmt.Sprintf(MntPath, containerName)

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
