package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

// This function creates the new parent process, it should
// be only invoked once. And wait for the child process (container)
// terminate, when the container terminate, parent should terminate
// too.
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {

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
