// These functions are dagent specific, yet it's more closely related to dockerhelper itself.
// Latter implementation is not straightforward at this moment as it would introduce circular dependencies

package utils

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	dockerhelper "github.com/dyrector-io/dyrectorio/golang/internal/helper/docker"
	"github.com/dyrector-io/dyrectorio/golang/internal/mapper"
	"github.com/dyrector-io/dyrectorio/protobuf/go/common"
)

func DeleteContainerByNameGrpc(ctx context.Context, prefix, name string) error {
	return dockerhelper.DeleteContainerByName(context.Background(), nil, fmt.Sprintf("%s-%s", prefix, name))
}

func GetContainersByNameCrux(ctx context.Context, nameFilter string) []*common.ContainerStateItem {
	containers, err := dockerhelper.GetContainersByName(ctx, nil, nameFilter)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	return mapper.MapContainerState(&containers)
}
