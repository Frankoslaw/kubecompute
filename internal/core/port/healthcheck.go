package port

import (
	"context"
	"kubecompute/internal/core/domain"
)

type LifecycleProvider interface {
	Supports(op domain.LifecycleOperation) bool
	Start(ctx context.Context, node *domain.Node) error
	Stop(ctx context.Context, node *domain.Node) error
	Reboot(ctx context.Context, node *domain.Node) error
}
