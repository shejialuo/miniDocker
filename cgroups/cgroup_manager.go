package cgroups

import (
	"miniDocker/cgroups/subsystems"

	"github.com/sirupsen/logrus"
)

// `CgroupManager` is a manager which manages
// all resources with the unique path.
type CgroupManager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subSystem := range subsystems.Subsystems {
		subSystem.Apply(c.Path, pid)
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSystem := range subsystems.Subsystems {
		subSystem.Set(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subSystem := range subsystems.Subsystems {
		if err := subSystem.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
