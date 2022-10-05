package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	// Get the name of the subsystem, such as cpu, memory
	Name() string
	// Set the resource limit in one cgroup of this subsystem
	Set(path string, res *ResourceConfig) error
	// Add process into the cgroup
	Apply(path string, pid int) error
	// Remove cgroup
	Remove(path string) error
}

var (
	Subsystems = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
