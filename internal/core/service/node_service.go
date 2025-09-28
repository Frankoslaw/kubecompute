package service

import (
	"context"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
)

type NodeService struct {
	repository port.NodeRepository
	controller port.NodeController
}

func NewNodeService(repository port.NodeRepository, controller port.NodeController) *NodeService {
	return &NodeService{
		repository: repository,
		controller: controller,
	}
}

func (ns *NodeService) CreateNode(ctx context.Context, node *domain.Node) error {
	if err := ns.repository.CreateNode(ctx, node); err != nil {
		return err
	}

	ns.controller.Enqueue(node.ObjectMeta.NamespacedName())
	return nil
}

func (ns *NodeService) UpdateNode(ctx context.Context, node *domain.Node) error {
	if err := ns.repository.UpdateNode(ctx, node); err != nil {
		return err
	}

	ns.controller.Enqueue(node.ObjectMeta.NamespacedName())
	return nil
}

func (ns *NodeService) DeleteNode(ctx context.Context, name domain.NamespacedName) error {
	err := ns.repository.SoftDeleteNode(ctx, name)
	if err != nil {
		return err
	}

	ns.controller.Enqueue(name)
	return nil
}
