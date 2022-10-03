package container

import (
	"os"
	"os/exec"
	"syscall"
)

// This function creates the new parent process, it should
// be only invoked once. And wait for the child process (container)
// terminate, when the container terminate, parent should terminate
// too.
func NewParentProcess(tty bool, command string) *exec.Cmd {

	// Here, `/proc/self/exe` means execute itself, for example
	// if we are in `bash`, we could use `/proc/self/exe` to spawn
	// a child bash shell. In this example, we actually
	// execute the process `miniDocker` again, here, we need
	// to call initCommand to initialize the container so we
	// pass "init" to the arguments.
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
