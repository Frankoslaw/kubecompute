package domain

type HealthcheckProviderKind string

const (
	HealthcheckProviderKindPing HealthcheckProviderKind = "ping"
)

type HealthCheckSpec struct {
	Type HealthcheckProviderKind
}
