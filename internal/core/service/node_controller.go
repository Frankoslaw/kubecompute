package service

import (
	"context"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
	"kubecompute/internal/core/util"
	"log/slog"
	"time"
)

type NodeController struct {
	repository port.NodeRepository
	reconciler port.NodeReconciler
	workQueue  *util.WorkQueue[port.ReconcileRequest]
}

func NewNodeController(repository port.NodeRepository, reconciler port.NodeReconciler, workQueue *util.WorkQueue[port.ReconcileRequest]) port.NodeController {
	return &NodeController{
		repository: repository,
		reconciler: reconciler,
		workQueue:  workQueue,
	}
}

func (nc *NodeController) Start(ctx context.Context) {
	go nc.periodicSync(ctx)
	go nc.worker(ctx)
}

func (nc *NodeController) Enqueue(name domain.NamespacedName) {
	nc.workQueue.Add(port.ReconcileRequest{Name: name})
}

func (nc *NodeController) EnqueueAll() {
	nodes, err := nc.repository.ListNodesWithDeleted(context.Background())
	if err != nil {
		slog.Error("error listing nodes for enqueue all", "error", err)
		return
	}
	for _, node := range nodes {
		nc.Enqueue(node.ObjectMeta.NamespacedName())
	}
}

func (nc *NodeController) periodicSync(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	nc.EnqueueAll()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			nc.EnqueueAll()
		}
	}
}

func (nc *NodeController) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			nc.workQueue.Close()
			return
		default:
			req, ok := nc.workQueue.Get()
			if !ok {
				return
			}
			nc.processNextWorkItem(ctx, req)
		}
	}
}

func (nc *NodeController) processNextWorkItem(ctx context.Context, req port.ReconcileRequest) {
	res, err := nc.reconciler.Reconcile(ctx, req)
	if err != nil {
		slog.Error("error reconciling node", "error", err, "node", req.Name)
		go nc.requeue(req, 5*time.Second)
		return
	}

	if res.Requeue {
		go nc.requeue(req, res.RequeueAfter)
	}
}

func (nc *NodeController) requeue(req port.ReconcileRequest, after time.Duration) {
	time.Sleep(after)
	nc.Enqueue(req.Name)
}
