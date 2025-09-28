package provision

import (
	"context"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerNodeProvider struct {
	cli *client.Client
}

func NewDockerNodeProvider() (port.NodeProvider, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerNodeProvider{
		cli: cli,
	}, nil
}

func (d *DockerNodeProvider) EnsureNode(ctx context.Context, desired *domain.Node) (bool, error) {
	if desired.Deleted {
		return false, nil
	}

	observed, err := d.observe(ctx, desired.ObjectMeta.NamespacedName())
	if err != nil {
		return false, err
	}
	if observed == nil {
		return true, d.create(ctx, desired)
	}

	var drifted bool = false
	if observed.Spec.Image != desired.Spec.Image {
		drifted = true
	}

	// simple cmd comparison, only checks suffix to allow for entrypoint overrides
	if !strings.HasSuffix(strings.Join(observed.Spec.Cmd, " "), strings.Join(desired.Spec.Cmd, " ")) {
		drifted = true
	}

	if drifted {
		err := d.remove(ctx, observed)
		if err != nil {
			return false, err
		}

		return true, d.create(ctx, desired)
	}

	return false, nil
}

func (d *DockerNodeProvider) DeleteNode(ctx context.Context, name domain.NamespacedName) (bool, error) {
	node, err := d.observe(ctx, name)
	if err != nil {
		return false, err
	}

	if node == nil {
		return false, nil
	}

	return true, d.remove(ctx, node)
}

func (d *DockerNodeProvider) observe(ctx context.Context, name domain.NamespacedName) (*domain.Node, error) {
	containers, err := d.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "/"+name.String() {
				cmd := strings.Split(c.Command, " ")

				return &domain.Node{
					ObjectMeta: name.ObjectMeta(),
					Spec: domain.NodeSpec{
						Image: c.Image,
						Cmd:   cmd,
					},
					Status: domain.NodeStatus{
						ContainerID: c.ID,
					},
				}, nil
			}
		}
	}

	return nil, nil
}

func (d *DockerNodeProvider) create(ctx context.Context, desired *domain.Node) error {
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: desired.Spec.Image,
		Cmd:   desired.Spec.Cmd,
	}, nil, nil, nil, desired.ObjectMeta.NamespacedName().String())
	if err != nil {
		return err
	}

	desired.Status.ContainerID = resp.ID
	return d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

func (d *DockerNodeProvider) remove(ctx context.Context, desired *domain.Node) error {
	return d.cli.ContainerRemove(ctx, desired.Status.ContainerID, container.RemoveOptions{Force: true})
}
