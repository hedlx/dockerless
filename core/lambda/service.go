package lambda

import (
	"context"
	"fmt"
	"time"

	model "github.com/hedlx/doless/core/model"
	util "github.com/hedlx/doless/core/util"
)

type service struct {
	bootstrapping ConcurrentSet[string]
	starting      ConcurrentSet[string]
}

type LambdaService interface {
	BootstrapRuntime(ctx context.Context, runtime *model.CreateRuntimeM) (*model.CreateRuntimeM, error)
	BootstrapLambda(ctx context.Context, lambda *model.CreateLambdaM) (*model.CreateLambdaM, error)
	Start(ctx context.Context, id string) error
}

func CreateLambdaService() LambdaService {
	return &service{
		bootstrapping: CreateConcurrentSet[string](),
		starting:      CreateConcurrentSet[string](),
	}
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
		return nil, fmt.Errorf("Lambda with '%s' archive is already in progress", cLambda.Archive)
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
		BaseObject: model.BaseObject{
			ID:        cLambda.ID,
			Name:      cLambda.Name,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
		Runtime:  cLambda.Runtime,
		Endpoint: cLambda.Endpoint,
	}); err != nil {
		return nil, err
	}

	return cLambda, nil
}

func (s service) Start(ctx context.Context, id string) error {
	if succ := s.starting.AddUniq(id); !succ {
		return fmt.Errorf("Lambda '%s' is already being started", id)
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

	_, err = DeployLambda(ctx, tar, model.LambdaMeta{
		Name:    lambda.Name,
		Runtime: lambda.Runtime,
	})

	if err != nil {
		return err
	}

	return nil
}
