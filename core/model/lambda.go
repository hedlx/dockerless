package model

type LambdaMeta struct {
	Name    string `json:"name"`
	Runtime string `json:"runtime"`
}

type LambdaInternal struct {
	Meta        LambdaMeta `json:"meta"`
	ContainerID string     `json:"containerID"`
}
