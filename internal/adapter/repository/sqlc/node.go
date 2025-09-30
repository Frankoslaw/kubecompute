// TODO: Allow for bypasing resource version checks if not passed
package sqlc

import (
	"context"
	"database/sql"
	db "kubecompute/db/sqlc"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
	"strings"
)

func dbNodeToDomain(dbNode db.Node) *domain.Node {
	return &domain.Node{
		// metadata
		ObjectMeta: domain.ObjectMeta{
			Namespace:       dbNode.Namespace,
			Name:            dbNode.Name,
			ResourceVersion: int(dbNode.ResourceVersion),
		},
		// spec
		Spec: domain.NodeSpec{
			Backend: dbNode.Backend,
			Image:   dbNode.Image,
			Cmd:     strings.Split(dbNode.Cmd, " "),
		},
		// status
		Status: domain.NodeStatus{
			ContainerID: dbNode.ContainerID.String,
		},
		// misc
		Deleted: dbNode.DeletedAt.Valid,
	}
}

type SqliteNodeRepository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewSqlcNodeRepository(dbInstance *sql.DB) port.NodeRepository {
	return &SqliteNodeRepository{
		db:      dbInstance,
		queries: db.New(dbInstance),
	}
}

func (s *SqliteNodeRepository) CreateNode(ctx context.Context, node *domain.Node) (*domain.Node, error) {
	row, err := s.queries.CreateNode(ctx, db.CreateNodeParams{
		// metadata
		Namespace: node.ObjectMeta.Namespace,
		Name:      node.ObjectMeta.Name,
		// spec
		Backend: node.Spec.Backend,
		Image:   node.Spec.Image,
		Cmd:     strings.Join(node.Spec.Cmd, " "),
		// status
		ContainerID: sql.NullString{String: node.Status.ContainerID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return dbNodeToDomain(row), nil
}

func (s *SqliteNodeRepository) SoftDeleteNode(ctx context.Context, node *domain.Node) (*domain.Node, error) {
	row, err := s.queries.SoftDeleteNode(ctx, db.SoftDeleteNodeParams{
		Namespace:       node.ObjectMeta.Namespace,
		Name:            node.ObjectMeta.Name,
		ResourceVersion: int64(node.ObjectMeta.ResourceVersion),
		// TODO: Fix this doubled parameter due to sqlc limitations
		Column3: int64(node.ObjectMeta.ResourceVersion),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: Signal conflict errors
			return nil, nil
		}
		return nil, err
	}
	return dbNodeToDomain(row), nil
}

func (s *SqliteNodeRepository) HardDeleteNode(ctx context.Context, node *domain.Node) (*domain.Node, error) {
	row, err := s.queries.HardDeleteNode(ctx, db.HardDeleteNodeParams{
		Namespace:       node.ObjectMeta.Namespace,
		Name:            node.ObjectMeta.Name,
		ResourceVersion: int64(node.ObjectMeta.ResourceVersion),
		// TODO: Fix this doubled parameter due to sqlc limitations
		Column3: int64(node.ObjectMeta.ResourceVersion),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: Signal conflict errors
			return nil, nil
		}
		return nil, err
	}
	return dbNodeToDomain(row), nil
}

func (s *SqliteNodeRepository) GetNode(ctx context.Context, name domain.NamespacedName) (*domain.Node, error) {
	row, err := s.queries.GetNode(ctx, db.GetNodeParams{
		Namespace: name.Namespace,
		Name:      name.Name,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dbNodeToDomain(row), nil
}

func (s *SqliteNodeRepository) GetNodeWithDeleted(ctx context.Context, name domain.NamespacedName) (*domain.Node, error) {
	row, err := s.queries.GetNodeWithDeleted(ctx, db.GetNodeWithDeletedParams{
		Namespace: name.Namespace,
		Name:      name.Name,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dbNodeToDomain(row), nil
}

func (s *SqliteNodeRepository) ListNodes(ctx context.Context, namespace string) ([]*domain.Node, error) {
	rows, err := s.queries.ListNodes(ctx, db.ListNodesParams{
		Namespace: namespace,
		// TODO: Fix this doubled parameter due to sqlc limitations
		Column1: namespace,
	})
	if err != nil {
		return nil, err
	}

	var nodes []*domain.Node
	for _, row := range rows {
		nodes = append(nodes, dbNodeToDomain(row))
	}
	return nodes, nil
}

func (s *SqliteNodeRepository) ListNodesWithDeleted(ctx context.Context) ([]*domain.Node, error) {
	rows, err := s.queries.ListNodesWithDeleted(ctx)
	if err != nil {
		return nil, err
	}

	var nodes []*domain.Node
	for _, row := range rows {
		nodes = append(nodes, dbNodeToDomain(row))
	}
	return nodes, nil
}

func (s *SqliteNodeRepository) UpdateNode(ctx context.Context, node *domain.Node) (*domain.Node, error) {
	row, err := s.queries.UpdateNode(ctx, db.UpdateNodeParams{
		// metadata
		Namespace: node.ObjectMeta.Namespace,
		Name:      node.ObjectMeta.Name,
		// TODO: Fix this doubled parameter due to sqlc limitations
		Column6:         int64(node.ObjectMeta.ResourceVersion),
		ResourceVersion: int64(node.ObjectMeta.ResourceVersion),
		// spec
		Image: node.Spec.Image,
		Cmd:   strings.Join(node.Spec.Cmd, " "),
		// status
		ContainerID: sql.NullString{String: node.Status.ContainerID, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: Signal conflict errors
			return nil, nil
		}
		return nil, err
	}
	return dbNodeToDomain(row), nil
}
