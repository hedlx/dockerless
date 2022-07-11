package endpoint

import (
	"context"

	api "github.com/hedlx/doless/client"
	"github.com/hedlx/doless/manager/db"
)

func GetEndpoint(ctx context.Context, id string) (*api.Endpoint, error) {
	return db.GetValue[api.Endpoint](ctx, "endpoint", id)
}

func GetEndpoints(ctx context.Context) ([]*api.Endpoint, error) {
	return db.GetValues[api.Endpoint](ctx, "endpoint")
}

func SetEndpoint(ctx context.Context, endpoint *api.Endpoint) error {
	return db.SetValue(ctx, "endpoint:"+endpoint.Id, endpoint)
}

func FindEndpoint(ctx context.Context, predicate func(val *api.Endpoint) bool) (*api.Endpoint, error) {
	return db.FindValue(ctx, "endpoint", predicate)
}
