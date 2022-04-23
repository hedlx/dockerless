package ops

import (
	"context"
	"fmt"

	api "github.com/hedlx/doless/client"
)

func CreateRuntime(ctx context.Context, name string, path string) (*api.Runtime, error) {
	uploadID, err := upload(ctx, path, false)
	if err != nil {
		return nil, err
	}

	createResp, _, err := client.RuntimeApi.
		CreateRuntime(ctx).
		CreateRuntime(api.CreateRuntime{
			Name:       name,
			Dockerfile: uploadID,
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error when calling `RuntimeApi.CreateRuntime``: %v", err)
	}

	return createResp, nil
}

func ListRuntimes(ctx context.Context) ([]api.Runtime, error) {
	listResp, _, err := client.RuntimeApi.
		ListRuntimes(ctx).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error when calling `RuntimeApi.ListRuntimes``: %v", err)
	}

	return listResp, nil
}
