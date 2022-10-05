package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type CpusetSubSystem struct{}

func (s *CpusetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuSet != "" {
			cpuSet := "cpuset.cpus"
			cpuSetPath := path.Join(subSystemCgroupPath, cpuSet)
			if err := os.WriteFile(cpuSetPath, []byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("set cgroup cpuset fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subSystemCgroupPath)
	} else {
		return err
	}
}

func (s *CpusetSubSystem) Apply(cgroupPath string, pid int) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		taskPath := path.Join(subSystemCgroupPath, "tasks")
		if err := os.WriteFile(taskPath, []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return err
	}
}

func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}
