package lambda

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	model "github.com/hedlx/doless/core/model"
)

var dockerCli *client.Client

func init() {
	var err error
	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
}

// DeployLambda is just a snippet and POC
// TODO: it should be complete docker service that can deploy, inspect, destroy and report about lambda issues
func DeployLambda(ctx context.Context, lambda io.Reader, meta model.LambdaMeta) (*model.LambdaInternal, error) {
	out, err := dockerCli.ImageBuild(ctx, lambda, types.ImageBuildOptions{
		Tags:   []string{meta.Name},
		Labels: map[string]string{"lambda": "true"},
	})
	if err != nil {
		return nil, err
	}

	defer out.Body.Close()
	// Just wait for the end
	// TODO: need to parse logs and report errors
	io.ReadAll(out.Body)

	container, err := dockerCli.ContainerCreate(ctx, &container.Config{
		Image: meta.Name,
	}, nil, nil, nil, "doless_"+meta.Name)
	if err != nil {
		dockerCli.ImageRemove(ctx, meta.Name, types.ImageRemoveOptions{})
		return nil, err
	}

	if err := dockerCli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		info, err := dockerCli.ContainerInspect(ctx, container.ID)

		if err != nil {
			dockerCli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
			dockerCli.ImageRemove(ctx, meta.Name, types.ImageRemoveOptions{})

			return nil, err
		}

		if info.State.Running {
			if info.State.Health == nil || info.State.Health.Status == "none" || info.State.Health.Status != "starting" {
				break
			}
		}

		time.Sleep(time.Second)
	}

	return &model.LambdaInternal{Meta: meta, ContainerID: container.ID}, nil
}
