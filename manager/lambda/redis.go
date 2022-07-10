package lambda

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/hedlx/doless/manager/logger"
	"github.com/hedlx/doless/manager/model"
	"github.com/hedlx/doless/manager/util"
)

var rdb *redis.Client
var DolessID = ""

func init() {
	redisEndpoint := util.GetStrVar("REDIS_ENDPOINT")
	rdb = redis.NewClient(&redis.Options{
		Addr: redisEndpoint,
	})

	for {
		res := rdb.Get(context.Background(), "doless-id")
		err := res.Err()
		if err == nil {
			DolessID = res.Val()
			break
		}

		if err == redis.Nil {
			DolessID = util.UUID()
			if err := rdb.Set(context.Background(), "doless-id", DolessID, 0).Err(); err != nil {
				panic(err)
			}
			break
		}

		time.Sleep(time.Second)
	}
}

func setValue(ctx context.Context, key string, val interface{}) error {
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

func SetLambda(ctx context.Context, lambda *model.LambdaM) error {
	return setValue(ctx, "lambda:"+lambda.ID, lambda)
}

func SetRuntime(ctx context.Context, runtime *model.RuntimeM) error {
	return setValue(ctx, "runtime:"+runtime.ID, runtime)
}

func scanValues[T model.LambdaM | model.RuntimeM](ctx context.Context, prefix string, handler func(x *T) bool) error {
	var cursor uint64
	traversed := map[string]bool{}

	for {
		var keys []string
		var err error

		keys, cursor, err = rdb.Scan(ctx, cursor, prefix+":*", 0).Result()

		if err != nil {
			logger.L.Error(
				"Failed to scan redis",
				zap.Error(err),
			)
			return err
		}

		for _, key := range keys {
			if traversed[key] {
				continue
			}

			traversed[key] = true
			val, err := getValueByKey[T](ctx, key)
			if err != nil {
				continue
			}

			if !handler(val) {
				return nil
			}
		}

		if cursor == 0 {
			return nil
		}
	}
}

func getValues[T model.LambdaM | model.RuntimeM](ctx context.Context, prefix string) ([]*T, error) {
	res := []*T{}

	err := scanValues(ctx, prefix, func(val *T) bool {
		res = append(res, val)
		return true
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetLambdas(ctx context.Context) ([]*model.LambdaM, error) {
	return getValues[model.LambdaM](ctx, "lambda")
}

func GetRuntimes(ctx context.Context) ([]*model.RuntimeM, error) {
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

func findValue[T model.LambdaM | model.RuntimeM](ctx context.Context, prefix string, predicate func(x *T) bool) (*T, error) {
	var res *T
	err := scanValues(ctx, prefix, func(val *T) bool {
		if predicate(val) {
			res = val
			return false
		}

		return true
	})

	return res, err
}

func FindLambda(ctx context.Context, predicate func(val *model.LambdaM) bool) (*model.LambdaM, error) {
	return findValue(ctx, "lambda", predicate)
}
