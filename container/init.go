package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	commandArray := readUserCommand()
	if len(commandArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	setUpMount()

	path, err := exec.LookPath(commandArray[0])
	if err != nil {
		logrus.Errorf("exec loop path error %v", err)
		return err
	}
	logrus.Infof("Find path %s", path)

	if err := syscall.Exec(path, commandArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	message := string(msg)
	return strings.Split(message, " ")
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Get current location error %v", err)
		return
	}
	logrus.Infof("Current location is %s", pwd)
	pivotRoot(pwd)

	// Now we need to mount /proc /dev some system information
	// MS_NOEXEC: Do not allow programs to be executed from this filesystem
	// MS_NOSUID: Do not honor set-user-ID and set-group-ID bits
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	// Here, we use mount bind, because pivot_root needs two different
	// file system.
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	// Now we call `pivot_root` to change the old root to mount into
	// `.pivot_root`.
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// Change the current working directory
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")

	// MNT_DETACH: Perform  a lazy unmount: make the mount unavailable
	// for new processes, immediately disconnect the filesystem and all
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	return os.Remove(pivotDir)
}
