package domain

type NodeSpec struct {
	Backend string   `json:"backend" yaml:"backend"`
	Image   string   `json:"image" yaml:"image"`
	Cmd     []string `json:"cmd" yaml:"cmd"`
	// TODO: Health and lifecycle specs
}

type NodeStatus struct {
	ContainerID string `json:"containerID" yaml:"containerID"`
}

type Node struct {
	ObjectMeta ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       NodeSpec   `json:"spec" yaml:"spec"`
	Status     NodeStatus `json:"status" yaml:"status"`

	// soft delete derived from DeletedAt being non-null
	Deleted bool `json:"-" yaml:"-"`
}
