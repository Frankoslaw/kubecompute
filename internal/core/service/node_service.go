package service

import (
	"context"
	"errors"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
)

var ErrConflict = errors.New("resource version conflict")

type NodeService struct {
	repository port.NodeRepository
	controller port.NodeController
}

func NewNodeService(repository port.NodeRepository, controller port.NodeController) port.NodeService {
	return &NodeService{
		repository: repository,
		controller: controller,
	}
}

func (ns *NodeService) CreateNode(ctx context.Context, node *domain.Node) (*domain.Node, error) {
	updated, err := ns.repository.CreateNode(ctx, node)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrConflict
	}

	ns.controller.Enqueue(updated.ObjectMeta.NamespacedName())
	return updated, nil
}

func (ns *NodeService) GetNode(ctx context.Context, name domain.NamespacedName) (*domain.Node, error) {
	return ns.repository.GetNode(ctx, name)
}

func (ns *NodeService) ListNodes(ctx context.Context) ([]*domain.Node, error) {
	return ns.repository.ListNodes(ctx)
}

func (ns *NodeService) UpdateNode(ctx context.Context, node *domain.Node) (*domain.Node, error) {
	updated, err := ns.repository.UpdateNode(ctx, node)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrConflict
	}

	ns.controller.Enqueue(updated.ObjectMeta.NamespacedName())
	return updated, nil
}

func (ns *NodeService) DeleteNode(ctx context.Context, node *domain.Node) (*domain.Node, error) {
	updated, err := ns.repository.SoftDeleteNode(ctx, node)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrConflict
	}

	ns.controller.Enqueue(updated.ObjectMeta.NamespacedName())
	return updated, nil
}
