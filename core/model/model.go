package model

import (
	"fmt"
	"regexp"
)

type BaseObject struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type BaseLambdaM struct {
	BaseObject
	Runtime  string `json:"runtime"`
	Endpoint string `json:"endpoint"`
}

type DockerM struct {
	Image       *string `json:"image,omitempty"`
	Container   *string `json:"container,omitempty"`
	ContainerID *string `json:"container_id,omitempty"`
	Status      string  `json:"status"`
}

type LambdaM struct {
	BaseLambdaM
	Docker DockerM `json:"docker"`
}

type CreateLambdaM struct {
	BaseLambdaM
	Archive string `json:"archive"`
}

func ValidateCreateLambdaM(lambda *CreateLambdaM) error {
	if lambda.Name == "" {
		return fmt.Errorf("'name' is required")
	}

	if lambda.Runtime == "" {
		return fmt.Errorf("'runtime' is required")
	}

	if lambda.Endpoint == "" {
		return fmt.Errorf("'endpoint' is required")
	}

	if err := ValidateEndpoint(lambda.Endpoint); err != nil {
		return err
	}

	if lambda.Archive == "" {
		return fmt.Errorf("'archive' is required")
	}

	return nil
}

var EndpointRegex = regexp.MustCompile("^(/[0-9a-zA-Z-_]+)+$")

func ValidateEndpoint(path string) error {
	if !EndpointRegex.MatchString(path) {
		return fmt.Errorf("'endpoint' doesn't conform regex: %s", EndpointRegex.String())
	}

	return nil
}

type UpgradeLambdaM struct {
	Archive string `json:"archive"`
}

func ValidateUpgradeLambdaM(lambda *UpgradeLambdaM) error {
	if lambda.Archive == "" {
		return fmt.Errorf("'archive' is required")
	}

	return nil
}

type RuntimeM struct {
	BaseObject
}

type CreateRuntimeM struct {
	RuntimeM
	Dockerfile string `json:"dockerfile"`
}

func ValidateCreateRuntimeM(runtime *CreateRuntimeM) error {
	if runtime.Name == "" {
		return fmt.Errorf("'name' is required")
	}

	if runtime.Dockerfile == "" {
		return fmt.Errorf("'dockerfile' is required")
	}

	return nil
}
