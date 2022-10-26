package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/dyrector-io/dyrectorio/golang/internal/dogger"
)

const dockerClientTimeoutSeconds = 30

func DeleteContainerByName(ctx context.Context, dog *dogger.DeploymentLogger, nameFilter string) error {
	matchedContainer, err := GetOneContainerByName(ctx, nil, nameFilter)
	if err != nil {
		return fmt.Errorf("builder could not get container (%s) to remove: %s", nameFilter, err.Error())
	}

	switch matchedContainer.State {
	case "running", "paused", "restarting":
		if err = StopContainerByName(ctx, nil, exactMatch(nameFilter)); err != nil {
			return fmt.Errorf("builder could not stop container (%s): %s", nameFilter, err.Error())
		}
		fallthrough
	case "exited", "dead", "created":
		if err = RemoveContainerByName(ctx, nil, exactMatch(nameFilter)); err != nil {
			return fmt.Errorf("builder could not remove container (%s): %s", nameFilter, err.Error())
		}
		return nil
	case "":
		// when there's no container we just skip it
		return nil
	default:
		return fmt.Errorf("builder could not determine the state (%s) of the container (%s) for deletion: %s",
			matchedContainer.State,
			nameFilter,
			err.Error())
	}
}

func StopContainerByName(ctx context.Context, dog *dogger.DeploymentLogger, containerName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	timeoutValue := (time.Duration(dockerClientTimeoutSeconds) * time.Second)
	if err := cli.ContainerStop(ctx, containerName, &timeoutValue); err != nil {
		return err
	}

	return nil
}

// Matches one
func RemoveContainerByName(ctx context.Context, dog *dogger.DeploymentLogger, nameFilter string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerRemove(ctx, exactMatch(nameFilter), types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}

// Check the existence of containers, then return it
func GetContainersByName(ctx context.Context, dog *dogger.DeploymentLogger, nameFilter string) ([]types.Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	filter := filters.NewArgs()
	filter.Add("name", nameFilter)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil
}

func GetAllContainers(ctx context.Context, dog *dogger.DeploymentLogger) ([]types.Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil
}

// Using exact match!
func GetOneContainerByName(ctx context.Context, dog *dogger.DeploymentLogger, nameFilter string) (types.Container, error) {
	containers, err := GetContainersByName(ctx, nil, exactMatch(nameFilter))
	if err != nil {
		return types.Container{}, err
	}

	switch len(containers) {
	case 1:
		return containers[0], nil
	case 0:
		return types.Container{}, nil
	default:
		return types.Container{}, fmt.Errorf("more than one matching container")
	}
}

func exactMatch(name string) string {
	return fmt.Sprintf("^%s$", name)
}
