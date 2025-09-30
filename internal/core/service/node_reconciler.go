package service

import (
	"context"
	"kubecompute/internal/core/port"
	"log/slog"
	"time"
)

type NodeReconciler struct {
	repository port.NodeRepository
	providers  map[string]port.NodeProvider
}

func NewNodeReconciler(repository port.NodeRepository) port.NodeReconciler {
	return &NodeReconciler{
		repository: repository,
		providers:  make(map[string]port.NodeProvider),
	}
}

func (r *NodeReconciler) RegisterProvider(name string, provider port.NodeProvider) {
	r.providers[name] = provider
}

func (r *NodeReconciler) Reconcile(ctx context.Context, req port.ReconcileRequest) (port.ReconcileResult, error) {
	node, err := r.repository.GetNodeWithDeleted(ctx, req.Name)
	if err != nil || node == nil {
		return port.ReconcileResult{}, err
	}

	provider, exists := r.providers[node.Spec.Backend]
	if !exists {
		slog.Warn("no provider registered for backend", "backend", node.Spec.Backend)
		return port.ReconcileResult{}, nil
	}

	var ok bool
	if node.Deleted {
		ok, err = provider.DeleteNode(ctx, node.ObjectMeta.NamespacedName())
		if err != nil {
			return port.ReconcileResult{Requeue: true, RequeueAfter: 5 * time.Second}, err
		}
		if !ok {
			// nothing changed
			return port.ReconcileResult{}, nil
		}

		_, err = r.repository.HardDeleteNode(ctx, node)
		if err != nil {
			return port.ReconcileResult{Requeue: true, RequeueAfter: 5 * time.Second}, err
		}
	} else {
		ok, err = provider.EnsureNode(ctx, node)
		if err != nil {
			return port.ReconcileResult{Requeue: true, RequeueAfter: 5 * time.Second}, err
		}
		if !ok {
			// nothing changed
			return port.ReconcileResult{}, nil
		}

		// TODO: Create separate method for updating status for safety
		_, err = r.repository.UpdateNode(ctx, node)
		if err != nil {
			return port.ReconcileResult{Requeue: true, RequeueAfter: 5 * time.Second}, err
		}
	}

	return port.ReconcileResult{}, nil
}
