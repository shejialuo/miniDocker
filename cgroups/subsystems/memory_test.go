package subsystems_test

import (
	"miniDocker/cgroups/subsystems"
	"os"
	"path"
	"strconv"
	"testing"
)

func TestMemoryCgroup(t *testing.T) {
	memorySubsystem := subsystems.MemorySubSystem{}
	resourceConfig := subsystems.ResourceConfig{
		MemoryLimit: "100m",
	}

	testGroupName := "testForMemory"

	if err := memorySubsystem.Set(testGroupName, &resourceConfig); err != nil {
		t.Fatalf("Set %s cgroup fail %v", testGroupName, err)
	}

	stat, _ := os.Stat(path.Join(subsystems.FindCgroupMountpoint("memory"), testGroupName))

	if !stat.IsDir() {
		t.Fatalf("%s is not a directory, it should be a directory", testGroupName)
	}
	if stat.Name() != testGroupName {
		t.Fatalf("The test group name should be %s, but it is %s", testGroupName, stat.Name())
	}

	memoryLimit, _ := os.ReadFile(path.Join(subsystems.FindCgroupMountpoint("memory"), testGroupName, "memory.limit_in_bytes"))
	memoryLimitString := string(memoryLimit)
	// Here we need to drop the last element.
	if memoryLimitString[0:len(memoryLimitString)-1] != "104857600" {
		t.Fatalf("The memory limit should be 104857600, but it is %s", memoryLimitString)
	}

	pid := os.Getpid()
	if err := memorySubsystem.Apply(testGroupName, pid); err != nil {
		t.Fatalf("Apply %d to %s cgroup fail %v", pid, testGroupName, err)
	}
	pidActualBytes, _ := os.ReadFile(path.Join(subsystems.FindCgroupMountpoint("memory"), testGroupName, "tasks"))
	pidActual := string(pidActualBytes)
	if strconv.Itoa(pid) != pidActual[0:len(pidActual)-1] {
		t.Fatalf("The pid should be %d, but it is %s", pid, string(pidActual))
	}

	// Here, we must remove the pid to the root
	// Because we add itself to the cgroup we created
	// Thus we can safely remove the directory
	memorySubsystem.Apply("", os.Getpid())

	if err := memorySubsystem.Remove(testGroupName); err != nil {
		t.Fatalf("Remove %s cgroup fail %v", testGroupName, err)
	}

}
