package domain

type LifecycleProviderKind string

const (
	LifecycleProviderKindDocker LifecycleProviderKind = "docker"
)

type LifecycleOperation string

const (
	LifecycleOperationStart  LifecycleOperation = "start"
	LifecycleOperationStop   LifecycleOperation = "stop"
	LifecycleOperationReboot LifecycleOperation = "reboot"
)

type PowerSpec struct {
	Type LifecycleProviderKind
	Ops  []LifecycleOperation
}
