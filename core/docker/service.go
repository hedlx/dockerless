package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/hedlx/doless/core/model"
	lo "github.com/samber/lo"
)

type service struct {
	client *client.Client
}

type DockerService interface {
	Create(ctx context.Context, lambda model.LambdaM, tar io.Reader) (string, error)
	Start(ctx context.Context, lambda model.LambdaM) error
	Stop(ctx context.Context, lambda model.LambdaM) error
	// Inspect(ctx context.Context, lambda model.LambdaM) (<-chan *types.ContainerJSON, error)
	Remove(ctx context.Context, lambda model.LambdaM) error
}

func NewDockerService() (DockerService, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}

	return &service{client: client}, nil
}

func (s *service) Create(ctx context.Context, lambda model.LambdaM, tar io.Reader) (string, error) {
	if lambda.Docker.Container == nil || lambda.Docker.Image == nil {
		return "", fmt.Errorf("lambda model is not complete")
	}

	images, err := s.client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return "", err
	}

	_, exists := lo.Find(images, func(image types.ImageSummary) bool {
		return len(image.RepoTags) > 0 && strings.Split(image.RepoTags[0], ":")[0] == *lambda.Docker.Image
	})
	if exists {
		return "", fmt.Errorf("Image already exists: %s", *lambda.Docker.Image)
	}

	out, err := s.client.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Tags:   []string{*lambda.Docker.Image},
		Labels: map[string]string{"lambda": "true"},
	})
	if err != nil {
		return "", err
	}

	defer out.Body.Close()

	errorMsg := ""
	scanner := bufio.NewScanner(out.Body)
	for scanner.Scan() {
		e := struct {
			Err *string `json:"error"`
		}{}

		if err := json.Unmarshal([]byte(scanner.Text()), &e); err != nil {
			return "", err
		}

		if e.Err != nil {
			errorMsg += *e.Err
		}
	}

	if errorMsg != "" {
		return "", fmt.Errorf(errorMsg)
	}

	container, err := s.client.ContainerCreate(ctx, &container.Config{
		Image: *lambda.Docker.Image,
	}, nil, nil, nil, *lambda.Docker.Container)
	if err != nil {
		s.client.ImageRemove(ctx, *lambda.Docker.Image, types.ImageRemoveOptions{})
		return "", err
	}

	return container.ID, nil
}

func (s *service) Start(ctx context.Context, lambda model.LambdaM) error {
	if lambda.Docker.ContainerID == nil {
		return fmt.Errorf("lambda model is not complete")
	}

	info, err := s.client.ContainerInspect(ctx, *lambda.Docker.ContainerID)
	if err != nil {
		return err
	}

	if info.State.Running || info.State.Restarting {
		return nil
	}

	if err := s.client.ContainerStart(ctx, *lambda.Docker.ContainerID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (s *service) Stop(ctx context.Context, lambda model.LambdaM) error {
	if lambda.Docker.ContainerID == nil {
		return fmt.Errorf("lambda model is not complete")
	}

	info, err := s.client.ContainerInspect(ctx, *lambda.Docker.ContainerID)
	if err != nil {
		return err
	}

	if !info.State.Running && !info.State.Restarting {
		return nil
	}

	if err := s.client.ContainerStop(ctx, *lambda.Docker.ContainerID, nil); err != nil {
		return err
	}

	return nil
}

func (s *service) Remove(ctx context.Context, lambda model.LambdaM) error {
	if lambda.Docker.Container == nil || lambda.Docker.Image == nil {
		return fmt.Errorf("lambda model is not complete")
	}

	if err := s.Stop(ctx, lambda); err != nil {
		return err
	}

	if err := s.client.ContainerRemove(ctx, *lambda.Docker.ContainerID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	if _, err := s.client.ImageRemove(ctx, *lambda.Docker.Image, types.ImageRemoveOptions{}); err != nil {
		return err
	}

	return nil
}
