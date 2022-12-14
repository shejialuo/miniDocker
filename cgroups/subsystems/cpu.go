package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct{}

func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuShare != "" {
			cpuShare := "cpu.shares"
			cpuSharePath := path.Join(subSystemCgroupPath, cpuShare)
			if err := os.WriteFile(cpuSharePath, []byte(res.CpuShare), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu share fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subSystemCgroupPath)
	} else {
		return err
	}
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int) error {
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

func (s *CpuSubSystem) Name() string {
	return "cpu"
}
