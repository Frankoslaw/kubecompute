package service

import (
	"context"
	"kubecompute/internal/core/port"
	"time"
)

type NodeReconciler struct {
	repository port.NodeRepository
	provider   port.NodeProvider
}

func NewNodeReconciler(repository port.NodeRepository, provider port.NodeProvider) port.NodeReconciler {
	return &NodeReconciler{
		repository: repository,
		provider:   provider,
	}
}

func (r *NodeReconciler) Reconcile(ctx context.Context, req port.ReconcileRequest) (port.ReconcileResult, error) {
	node, err := r.repository.GetNodeWithDeleted(ctx, req.Name)
	if err != nil || node == nil {
		return port.ReconcileResult{}, err
	}

	var ok bool
	if node.Deleted {
		ok, err = r.provider.DeleteNode(ctx, node.ObjectMeta.NamespacedName())
	} else {
		ok, err = r.provider.EnsureNode(ctx, node)
	}

	if err != nil {
		return port.ReconcileResult{Requeue: true, RequeueAfter: 5 * time.Second}, err
	}

	if ok {
		err = r.repository.UpdateNode(ctx, node)
		if err != nil {
			return port.ReconcileResult{Requeue: true, RequeueAfter: 5 * time.Second}, err
		}
	}

	return port.ReconcileResult{}, nil
}
