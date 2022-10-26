package docker_test

import (
	"context"

	dockerhelper "github.com/dyrector-io/dyrectorio/golang/internal/helper/docker"
)

func TestDeleteContainer(ctx context.Context, container string) error {
	return dockerhelper.DeleteContainerByName(ctx, nil, container)
}

func TestDeleteNetwork(ctx context.Context, networkID string) error {
	return dockerhelper.DeleteNetworkByID(ctx, networkID)
}
