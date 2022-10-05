package subsystems_test

import (
	"miniDocker/cgroups/subsystems"
	"os"
	"path"
	"strconv"
	"testing"
)

func TestCpuCgroup(t *testing.T) {
	cpuSubsystem := subsystems.CpuSubSystem{}
	resourceConfig := subsystems.ResourceConfig{
		CpuShare: "30",
	}

	testGroupName := "testForCpu"

	if err := cpuSubsystem.Set(testGroupName, &resourceConfig); err != nil {
		t.Fatalf("Set %s cgroup fail %v", testGroupName, err)
	}

	stat, _ := os.Stat(path.Join(subsystems.FindCgroupMountpoint("cpu"), testGroupName))

	if !stat.IsDir() {
		t.Fatalf("%s is not a directory, it should be a directory", testGroupName)
	}

	if stat.Name() != testGroupName {
		t.Fatalf("The test group name should be %s, but it is %s", testGroupName, stat.Name())
	}

	cpuShare, _ := os.ReadFile(path.Join(subsystems.FindCgroupMountpoint("cpu"), testGroupName, "cpu.shares"))
	cpuShareString := string(cpuShare)
	// Here we need to drop the last element.
	if cpuShareString[0:len(cpuShareString)-1] != "30" {
		t.Fatalf("The cpu share should be 30, but it is %s", cpuShareString)
	}

	pid := os.Getpid()
	if err := cpuSubsystem.Apply(testGroupName, pid); err != nil {
		t.Fatalf("Apply %d to %s cgroup fail %v", pid, testGroupName, err)
	}
	pidActualBytes, _ := os.ReadFile(path.Join(subsystems.FindCgroupMountpoint("cpu"), testGroupName, "tasks"))
	pidActual := string(pidActualBytes)
	if strconv.Itoa(pid) != pidActual[0:len(pidActual)-1] {
		t.Fatalf("The pid should be %d, but it is %s", pid, string(pidActual))
	}

	// Here, we must remove the pid to the root
	// Because we add itself to the cgroup we created
	// Thus we can safely remove the directory
	cpuSubsystem.Apply("", os.Getpid())

	if err := cpuSubsystem.Remove(testGroupName); err != nil {
		t.Fatalf("Remove %s cgroup fail %v", testGroupName, err)
	}

}
