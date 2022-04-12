package lambda

import (
	"context"
	"encoding/json"

	redis "github.com/go-redis/redis/v8"
	zap "go.uber.org/zap"

	logger "github.com/hedlx/doless/core/logger"
	model "github.com/hedlx/doless/core/model"
	util "github.com/hedlx/doless/core/util"
)

var rdb *redis.Client

func init() {
	redisEndpoint := util.GetStrVar("REDIS_ENDPOINT")
	rdb = redis.NewClient(&redis.Options{
		Addr: redisEndpoint,
	})
}

func addValue(ctx context.Context, key string, val interface{}) error {
	// Notice: datarace is possible here, but due to specific of addValue usage
	// it is not addValue responsibility to manage this at all.

	obj, err := json.Marshal(val)
	if err != nil {
		return err
	}

	status := rdb.Set(ctx, key, string(obj), 0)
	if err := status.Err(); err != nil {
		return err
	}

	return nil
}

func AddLambda(ctx context.Context, lambda *model.LambdaM) error {
	return addValue(ctx, "lambda:"+lambda.ID, lambda)
}

func AddRuntime(ctx context.Context, runtime *model.RuntimeM) error {
	return addValue(ctx, "runtime:"+runtime.ID, runtime)
}

func getValues[T model.LambdaM | model.RuntimeM](ctx context.Context, prefix string) (<-chan *T, <-chan error) {
	resC := make(chan *T)
	errC := make(chan error)

	go func() {
		var cursor uint64

		for {
			var keys []string
			var err error

			keys, cursor, err = rdb.Scan(ctx, cursor, prefix+":*", 0).Result()

			if err != nil {
				close(resC)
				errC <- err

				logger.L.Error(
					"Failed to scan redis",
					zap.Error(err),
				)
				return
			}

			for _, key := range keys {
				val, err := getValueByKey[T](ctx, key)
				if err != nil {
					continue
				}

				resC <- val
			}

			if cursor == 0 {
				close(resC)
				return
			}
		}
	}()

	return resC, errC
}

func GetLambdas(ctx context.Context) (<-chan *model.LambdaM, <-chan error) {
	return getValues[model.LambdaM](ctx, "lambda")
}

func GetRuntimes(ctx context.Context) (<-chan *model.RuntimeM, <-chan error) {
	return getValues[model.RuntimeM](ctx, "runtime")
}

func getValueByKey[T model.LambdaM | model.RuntimeM](ctx context.Context, key string) (*T, error) {
	rawVal, err := rdb.Get(ctx, key).Result()
	if err != nil {
		logger.L.Error(
			"Failed to get redis members",
			zap.Error(err),
			zap.String("key", key),
		)
		return nil, err
	}

	var val T

	if err := json.Unmarshal([]byte(rawVal), &val); err != nil {
		logger.L.Error(
			"Failed to parse redis value",
			zap.Error(err),
			zap.String("key", key),
			zap.String("value", rawVal),
		)
		return nil, err
	}

	return &val, err
}

func getValue[T model.LambdaM | model.RuntimeM](ctx context.Context, prefix string, id string) (*T, error) {
	key := prefix + ":" + id
	return getValueByKey[T](ctx, key)
}

func GetLambda(ctx context.Context, id string) (*model.LambdaM, error) {
	return getValue[model.LambdaM](ctx, "lambda", id)
}

func GetRuntime(ctx context.Context, id string) (*model.RuntimeM, error) {
	return getValue[model.RuntimeM](ctx, "runtime", id)
}
