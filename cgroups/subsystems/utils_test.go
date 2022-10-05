package subsystems_test

import (
	"miniDocker/cgroups/subsystems"
	"testing"
)

func TestFindCgroupMountpoint(t *testing.T) {
	cpuPath := "/sys/fs/cgroup/cpu"
	memoryPath := "/sys/fs/cgroup/memory"
	cpusetPath := "/sys/fs/cgroup/cpuset"

	cpuPathActual := subsystems.FindCgroupMountpoint("cpu")
	memoryPathActual := subsystems.FindCgroupMountpoint("memory")
	cpusetPathActual := subsystems.FindCgroupMountpoint("cpuset")
	if cpuPath != cpuPathActual {
		t.Errorf("cpu path should be %v, actual %v", cpuPath, cpuPathActual)
	}
	if memoryPath != memoryPathActual {
		t.Errorf("memory path should be %v, actual %v", memoryPath, memoryPathActual)
	}
	if cpusetPath != cpusetPathActual {
		t.Errorf("cpuset path should be %v, actual %v", cpusetPath, cpusetPathActual)
	}
}
