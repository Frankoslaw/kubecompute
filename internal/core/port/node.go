package port

import (
	"context"
	"kubecompute/internal/core/domain"
	"time"
)

type NodeRepository interface {
	CreateNode(ctx context.Context, node *domain.Node) (*domain.Node, error)
	GetNode(ctx context.Context, name domain.NamespacedName) (*domain.Node, error)
	ListNodes(ctx context.Context, namespace string) ([]*domain.Node, error)
	UpdateNode(ctx context.Context, node *domain.Node) (*domain.Node, error)
	SoftDeleteNode(ctx context.Context, node *domain.Node) (*domain.Node, error)
	// Internal use only
	GetNodeWithDeleted(ctx context.Context, name domain.NamespacedName) (*domain.Node, error)
	ListNodesWithDeleted(ctx context.Context) ([]*domain.Node, error)
}

type NodeService interface {
	CreateNode(ctx context.Context, node *domain.Node) (*domain.Node, error)
	GetNode(ctx context.Context, name domain.NamespacedName) (*domain.Node, error)
	ListNodes(ctx context.Context, namespace *string) ([]*domain.Node, error)
	UpdateNode(ctx context.Context, node *domain.Node) (*domain.Node, error)
	DeleteNode(ctx context.Context, node *domain.Node) (*domain.Node, error)
}

type NodeController interface {
	Start(ctx context.Context)
	Enqueue(name domain.NamespacedName)
	EnqueueAll()
}

type ReconcileRequest struct {
	Name domain.NamespacedName
}

type ReconcileResult struct {
	Requeue      bool
	RequeueAfter time.Duration
}

type NodeReconciler interface {
	Reconcile(ctx context.Context, req ReconcileRequest) (ReconcileResult, error)
}

type NodeProvider interface {
	EnsureNode(ctx context.Context, node *domain.Node) (bool, error)
	DeleteNode(ctx context.Context, name domain.NamespacedName) (bool, error)
}
