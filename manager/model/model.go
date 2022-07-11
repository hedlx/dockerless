package model

import (
	"fmt"
	"regexp"

	api "github.com/hedlx/doless/client"
)

func ValidateCreateLambda(lambda *api.CreateLambda) error {
	if lambda.Name == "" {
		return fmt.Errorf("'name' is required")
	}

	if lambda.Runtime == "" {
		return fmt.Errorf("'runtime' is required")
	}

	if lambda.LambdaType == "" {
		return fmt.Errorf("'lambda_type' is required")
	}

	if lambda.LambdaType != "ENDPOINT" && lambda.LambdaType != "INTERNAL" {
		return fmt.Errorf("invalid 'lambda_type' value: %s", lambda.LambdaType)
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

func ValidateCreateRuntime(runtime *api.CreateRuntime) error {
	if runtime.Name == "" {
		return fmt.Errorf("'name' is required")
	}

	if runtime.Dockerfile == "" {
		return fmt.Errorf("'dockerfile' is required")
	}

	return nil
}
