package docker

import (
	"context"
	"fmt"
	"main/internal/storage"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func GetMinioConfig(ctx context.Context, cli *client.Client, filterKey, filterValue, networkName string) ([]storage.MinioConfig, error) {
	args := filters.NewArgs()
	args.Add("label", fmt.Sprintf("%s=%s", filterKey, filterValue))

	containers, err := cli.ContainerList(ctx, container.ListOptions{Filters: args})

	if err != nil {
		return nil, err
	}

	var configs []storage.MinioConfig

	for _, containerInfo := range containers {
		_, ok := containerInfo.NetworkSettings.Networks[networkName]
		if !ok {
			continue
		}

		containerData, containerErr := cli.ContainerInspect(ctx, containerInfo.ID)

		if containerErr != nil {
			return nil, err
		}

		envMap := make(map[string]string)
		for _, env := range containerData.Config.Env {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}

		cfg := storage.MinioConfig{
			Id:        containerData.Name,
			Endpoint:  containerData.NetworkSettings.Networks[networkName].IPAddress + ":9000",
			AccessKey: envMap["MINIO_ACCESS_KEY"],
			SecretKey: envMap["MINIO_SECRET_KEY"],
		}

		configs = append(configs, cfg)
	}

	return configs, err
}
