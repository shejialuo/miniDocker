package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// For Cgroup-v1, the hierarchy provided by kernel
// is already mounted for us, so we need to find the
// corresponding path from `/proc/self/mountinfo`. Use
// `cat /proc/self/mountinfo` to see the fields to
// understand the logic of this function
//
// @param subsystem: the name of the subsystem
//
// @return the corresponding path
func FindCgroupMountpoint(subsystem string) string {
	const mountInfo = "/proc/self/mountinfo"
	f, err := os.Open(mountInfo)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// For each hierarchy, it is a directory in the file system.
// This function is used to get the absolute path also we could
// create a new hierarchy.
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	// When the file exists or the file not exists with autoCreate is true
	// we should go on
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
