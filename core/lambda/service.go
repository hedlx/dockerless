package lambda

import (
	"context"
	"fmt"
	"time"

	"github.com/hedlx/doless/core/common"
	"github.com/hedlx/doless/core/docker"
	"github.com/hedlx/doless/core/logger"
	"github.com/hedlx/doless/core/model"
	"github.com/hedlx/doless/core/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type service struct {
	dockerSvc     docker.DockerService
	bootstrapping common.ConcurrentSet[string]
	starting      common.ConcurrentSet[string]
	lambdas       common.ConcurrentMap[string, model.LambdaM]
	inspect       common.ConcurrentMap[string, func()]
}

type LambdaService interface {
	Init() error
	Stop(ctx context.Context)
	BootstrapRuntime(ctx context.Context, runtime *model.CreateRuntimeM) (*model.CreateRuntimeM, error)
	BootstrapLambda(ctx context.Context, lambda *model.CreateLambdaM) (*model.CreateLambdaM, error)
	Start(ctx context.Context, id string) error
	Destroy(ctx context.Context, id string) error
}

func CreateLambdaService() (LambdaService, error) {
	dockerSvc, err := docker.NewDockerService(DolessID)

	if err != nil {
		return nil, err
	}

	svc := &service{
		dockerSvc:     dockerSvc,
		bootstrapping: common.CreateConcurrentSet[string](),
		starting:      common.CreateConcurrentSet[string](),
		lambdas:       common.CreateConcurrentMap[string, model.LambdaM](),
		inspect:       common.CreateConcurrentMap[string, func()](),
	}

	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func (s service) Init() error {
	ctx := context.Background()
	lambdaC, errC := GetLambdas(ctx)

	hasLambdas := true
	for hasLambdas {
		select {
		case err := <-errC:
			return err
		case lambda, ok := <-lambdaC:
			if !ok {
				hasLambdas = false
				break
			}

			s.lambdas.Set(lambda.ID, *lambda)
		}
	}

	var lambdaInitErr error
	lo.ForEach(s.lambdas.Values(), func(lambda model.LambdaM, _ int) {
		if lambda.Docker.ContainerID == nil {
			return
		}

		container, err := s.dockerSvc.Inspect(ctx, *lambda.Docker.ContainerID)

		// TODO: check error more precisely and handle correctly
		if err != nil {
			id, err := s.dockerSvc.CreateContainer(ctx, &lambda)
			if err != nil {
				id, lambdaInitErr = s.start(ctx, &lambda)
			}

			if id != "" {
				lambda.Docker.ContainerID = &id
				s.updateLambda(ctx, lambda)
			}
		}

		if err != nil || (!container.State.Running && !container.State.Restarting) {
			s.dockerSvc.Start(ctx, &lambda)
		}
	})

	if lambdaInitErr != nil {
		return lambdaInitErr
	}

	s.lambdas.ForEach(func(_ string, lambda model.LambdaM) {
		if lambda.Docker.ContainerID == nil {
			return
		}

		lctx, cancel := context.WithCancel(context.Background())
		s.inspect.Set(lambda.ID, cancel)
		go s.inspectRoutine(lctx, lambda)
	})

	return nil
}

func (s *service) Stop(ctx context.Context) {
	s.inspect.ForEach(func(_ string, stop func()) {
		stop()
	})

	s.lambdas.ForEach(func(_ string, lambda model.LambdaM) {
		if lambda.Docker.ContainerID == nil {
			return
		}

		s.dockerSvc.Stop(ctx, &lambda)
	})
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

	if err := SetRuntime(ctx, &model.RuntimeM{
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

	lambda := model.LambdaM{
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
	}

	if err := SetLambda(ctx, &lambda); err != nil {
		return nil, err
	}

	s.lambdas.Set(lambda.ID, lambda)

	return cLambda, nil
}

func (s service) start(ctx context.Context, lambda *model.LambdaM) (string, error) {
	tar, err := TarLambda(ctx, lambda.ID, lambda.Runtime)
	if err != nil {
		return "", err
	}

	image := lambda.Name
	container := "doless-" + lambda.Name
	lambda.Docker.Image = &image
	lambda.Docker.Container = &container

	containerID, err := s.dockerSvc.Create(ctx, lambda, tar)
	if err != nil {
		return "", err
	}

	lambda.Docker.ContainerID = &containerID

	return containerID, s.dockerSvc.Start(ctx, lambda)
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

	if _, err = s.start(ctx, lambda); err != nil {
		return err
	}

	if err := s.updateLambda(ctx, *lambda); err != nil {
		return err
	}

	lctx, cancel := context.WithCancel(context.Background())
	s.inspect.Set(lambda.ID, cancel)
	go s.inspectRoutine(lctx, *lambda)

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

	s.inspect.Get(id, func() {})()
	s.inspect.Delete(id)

	if err := s.dockerSvc.Remove(ctx, lambda); err != nil {
		return err
	}

	lambda.Docker = model.DockerM{}

	if err := s.updateLambda(ctx, *lambda); err != nil {
		return err
	}

	return nil
}

func (s service) updateLambda(ctx context.Context, lambda model.LambdaM) error {
	var updateErr error
	s.lambdas.Update(lambda.ID, func(prev model.LambdaM) model.LambdaM {
		if err := SetLambda(ctx, &lambda); err != nil {
			updateErr = err
			return prev
		}

		return lambda
	})

	return updateErr
}

func (s service) inspectRoutine(ctx context.Context, lambda model.LambdaM) {
	id := *lambda.Docker.ContainerID
	for {
		container, err := s.dockerSvc.Inspect(ctx, id)
		actual, rErr := GetLambda(ctx, lambda.ID)

		if rErr == nil {
			if err != nil {
				// TODO: Handle external delete case and stop
				logger.L.Error(
					"Failed to inspect container",
					zap.Error(err),
					zap.String("container_id", id),
				)
				lambda.Docker.Status = "error"
			} else {
				lambda.Docker.Status = container.State.Health.Status
			}

			if actual.Docker.Status != lambda.Docker.Status {
				if err := s.updateLambda(ctx, lambda); err != nil {
					logger.L.Error(
						"Failed to update lambda",
						zap.Error(err),
						zap.String("id", lambda.ID),
					)
				}
			}
		} else {
			logger.L.Error(
				"Failed to gather lambda",
				zap.Error(err),
				zap.String("id", lambda.ID),
			)
		}

		select {
		case <-time.After(10 * time.Second):
			continue
		case <-ctx.Done():
			return
		}
	}
}
