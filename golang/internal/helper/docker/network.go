package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func GetNetworks() ([]types.NetworkResource, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return cli.NetworkList(ctx, types.NetworkListOptions{})
}

func CreateNetwork(ctx context.Context, name, driver string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	filter := filters.NewArgs()
	filter.Add("name", name)

	networkListOption := types.NetworkListOptions{
		Filters: filter,
	}

	networks, err := cli.NetworkList(ctx, networkListOption)
	if err != nil {
		return fmt.Errorf("error list existing networks: %w", err)
	}

	if len(networks) > 0 {
		log.Printf("Provided network name: %s is exists. Skip to create new network.", name)
		return nil
	}

	networkCreateOptions := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         driver,
	}

	_, err = cli.NetworkCreate(ctx, name, networkCreateOptions)

	if err != nil {
		return err
	}

	return nil
}

func DeleteNetworkByID(ctx context.Context, networkID string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return cli.NetworkRemove(ctx, networkID)
}
