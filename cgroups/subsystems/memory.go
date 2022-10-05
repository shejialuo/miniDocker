package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct{}

func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			memoryLimit := "memory.limit_in_bytes"
			memoryLimitPath := path.Join(subSystemCgroupPath, memoryLimit)
			if err := os.WriteFile(memoryLimitPath, []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subSystemCgroupPath)
	} else {
		return err
	}
}

func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
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

func (s *MemorySubSystem) Name() string {
	return "memory"
}
