package util

import (
	"os"
	"strconv"

	"github.com/hedlx/doless/manager/logger"
	"go.uber.org/zap"
)

func GetStrVar(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logger.L.Fatal("env var is missing", zap.String("env_var", name))
	}

	return v
}

func GetIntVar(name string) int {
	s := GetStrVar(name)
	v, err := strconv.Atoi(s)

	if err != nil {
		panic(err)
	}

	return v
}
