package lambda

import (
	"context"

	api "github.com/hedlx/doless/client"
	"github.com/hedlx/doless/manager/db"
)

func GetLambda(ctx context.Context, id string) (*api.Lambda, error) {
	return db.GetValue[api.Lambda](ctx, "lambda", id)
}

func GetRuntime(ctx context.Context, id string) (*api.Runtime, error) {
	return db.GetValue[api.Runtime](ctx, "runtime", id)
}

func GetLambdas(ctx context.Context) ([]*api.Lambda, error) {
	return db.GetValues[api.Lambda](ctx, "lambda")
}

func GetRuntimes(ctx context.Context) ([]*api.Runtime, error) {
	return db.GetValues[api.Runtime](ctx, "runtime")
}

func SetLambda(ctx context.Context, lambda *api.Lambda) error {
	return db.SetValue(ctx, "lambda:"+lambda.Id, lambda)
}

func SetRuntime(ctx context.Context, runtime *api.Runtime) error {
	return db.SetValue(ctx, "runtime:"+runtime.Id, runtime)
}

func FindLambda(ctx context.Context, predicate func(val *api.Lambda) bool) (*api.Lambda, error) {
	return db.FindValue(ctx, "lambda", predicate)
}
