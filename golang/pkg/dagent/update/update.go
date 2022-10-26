package update

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/dyrector-io/dyrectorio/golang/pkg/dagent/config"
	"github.com/dyrector-io/dyrectorio/golang/pkg/dagent/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func tryParseCGroupFile() (string, error) {
	cgroupFile, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return "", err
	}

	defer cgroupFile.Close()

	scanner := bufio.NewScanner(cgroupFile)
	if scanner.Scan() {
		group := scanner.Text()
		lastSlash := strings.LastIndex(group, "/")
		if lastSlash < 0 {
			return group, nil
		}
		return group[lastSlash+1:], nil
	}

	return "", scanner.Err()
}

func getSelfID() string {
	cgroup, err := tryParseCGroupFile()
	if err != nil {
		return os.Getenv("HOSTNAME")
	}

	return cgroup
}

func InitUpdater(cfg *config.Configuration) {
	switch cfg.UpdateMethod {
	case "poll":
		log.Print("Update mode: polling")
		log.Print("Remote DAgent image: " + cfg.DagentImage + ":" + cfg.DagentTag)
		if err := utils.ExecWatchtowerPoll(context.Background(), cfg); err != nil {
			log.Error().Stack().Err(err).Msg("Error starting watchtower")
		}
	case "off":
	default:
		log.Print("No update was set up")
	}
}

func SelfUpdate(ctx context.Context) error {
	containerID := getSelfID()
	if containerID == "" {
		return errors.New("unable to get self container ID")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "id", Value: containerID}),
	})
	if err != nil {
		return err
	}

	if len(containers) != 1 {
		return errors.New("unable to find self")
	}

	container := containers[0]

	fmt.Println(container.Image)

	return nil
}
