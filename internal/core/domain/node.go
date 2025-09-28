package domain

type NodeSpec struct {
	Image string   `yaml:"image"`
	Cmd   []string `yaml:"cmd"`
	// TODO
	// Health    HealthCheckSpec
	// PowerSpec PowerSpec
}

type NodeStatus struct {
	ContainerID string `yaml:"containerID"`
}

type Node struct {
	ObjectMeta ObjectMeta `yaml:"metadata"`
	Spec       NodeSpec   `yaml:"spec"`
	Status     NodeStatus `yaml:"status"`

	// soft delete derived from DeletedAt being non-null
	Deleted bool `yaml:"-"`
}
