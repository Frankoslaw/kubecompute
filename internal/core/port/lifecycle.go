package port

import (
	"context"
	"kubecompute/internal/core/domain"
)

type HealthcheckProvider interface {
	Check(ctx context.Context, node *domain.Node) (bool, error)
}
