package lambda

import (
	"context"
	"fmt"
	"time"

	"github.com/hedlx/doless/core/docker"
	model "github.com/hedlx/doless/core/model"
	util "github.com/hedlx/doless/core/util"
)

type service struct {
	dockerSvc     docker.DockerService
	bootstrapping ConcurrentSet[string]
	starting      ConcurrentSet[string]
}

type LambdaService interface {
	BootstrapRuntime(ctx context.Context, runtime *model.CreateRuntimeM) (*model.CreateRuntimeM, error)
	BootstrapLambda(ctx context.Context, lambda *model.CreateLambdaM) (*model.CreateLambdaM, error)
	Start(ctx context.Context, id string) error
	Destroy(ctx context.Context, id string) error
}

func CreateLambdaService() (LambdaService, error) {
	dockerSvc, err := docker.NewDockerService()

	if err != nil {
		return nil, err
	}

	return &service{
		dockerSvc:     dockerSvc,
		bootstrapping: CreateConcurrentSet[string](),
		starting:      CreateConcurrentSet[string](),
	}, nil
}

func (s *service) BootstrapRuntime(ctx context.Context, cRuntime *model.CreateRuntimeM) (*model.CreateRuntimeM, error) {
	if succ := s.bootstrapping.AddUniq(cRuntime.Dockerfile); !succ {
		return nil, fmt.Errorf("Lambda with '%s' archive is already in progress", cRuntime.Dockerfile)
	}
	defer s.bootstrapping.Remove(cRuntime.Dockerfile)

	cRuntime.ID = util.UUID()

	if err := BootstrapRuntime(ctx, cRuntime); err != nil {
		return nil, err
	}

	createdAt := time.Now().UnixMilli()
	cRuntime.CreatedAt = createdAt
	cRuntime.UpdatedAt = createdAt

	if err := AddRuntime(ctx, &model.RuntimeM{
		BaseObject: model.BaseObject{
			ID:        cRuntime.ID,
			Name:      cRuntime.Name,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
	}); err != nil {
		return nil, err
	}

	return cRuntime, nil
}

func (s *service) BootstrapLambda(ctx context.Context, cLambda *model.CreateLambdaM) (*model.CreateLambdaM, error) {
	if succ := s.bootstrapping.AddUniq(cLambda.Archive); !succ {
		return nil, fmt.Errorf("Lambda with '%s' archive is already being bootstrapped", cLambda.Archive)
	}
	defer s.bootstrapping.Remove(cLambda.Archive)

	if _, err := GetRuntime(ctx, cLambda.Runtime); err != nil {
		return nil, err
	}

	cLambda.ID = util.UUID()

	if err := BootstrapLambda(ctx, cLambda); err != nil {
		return nil, err
	}

	createdAt := time.Now().UnixMilli()
	cLambda.CreatedAt = createdAt
	cLambda.UpdatedAt = createdAt

	if err := AddLambda(ctx, &model.LambdaM{
		BaseLambdaM: model.BaseLambdaM{
			BaseObject: model.BaseObject{
				ID:        cLambda.ID,
				Name:      cLambda.Name,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
			Runtime:  cLambda.Runtime,
			Endpoint: cLambda.Endpoint,
		},
	}); err != nil {
		return nil, err
	}

	return cLambda, nil
}

func (s service) Start(ctx context.Context, id string) error {
	if succ := s.starting.AddUniq(id); !succ {
		return fmt.Errorf("Lambda '%s' is already being processed", id)
	}
	defer s.starting.Remove(id)

	lambda, err := GetLambda(ctx, id)
	if err != nil {
		return err
	}

	tar, err := TarLambda(ctx, lambda.ID, lambda.Runtime)
	if err != nil {
		return err
	}

	image := lambda.Name
	container := "doless-" + lambda.Name
	lambda.Docker.Image = &image
	lambda.Docker.Container = &container

	containerID, err := s.dockerSvc.Create(ctx, *lambda, tar)
	if err != nil {
		return err
	}

	lambda.Docker.ContainerID = &containerID

	if err := AddLambda(ctx, lambda); err != nil {
		return err
	}

	if err := s.dockerSvc.Start(ctx, *lambda); err != nil {
		return err
	}

	return nil
}

func (s service) Destroy(ctx context.Context, id string) error {
	if succ := s.starting.AddUniq(id); !succ {
		return fmt.Errorf("Lambda '%s' is already being processed", id)
	}
	defer s.starting.Remove(id)

	lambda, err := GetLambda(ctx, id)
	if err != nil {
		return err
	}

	if err := s.dockerSvc.Remove(ctx, *lambda); err != nil {
		return err
	}

	lambda.Docker = model.DockerM{}

	if err := AddLambda(ctx, lambda); err != nil {
		return err
	}

	return nil
}
