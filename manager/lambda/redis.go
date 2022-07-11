package lambda

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	api "github.com/hedlx/doless/client"
	"github.com/hedlx/doless/manager/logger"
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

func SetLambda(ctx context.Context, lambda *api.Lambda) error {
	return setValue(ctx, "lambda:"+lambda.Id, lambda)
}

func SetRuntime(ctx context.Context, runtime *api.Runtime) error {
	return setValue(ctx, "runtime:"+runtime.Id, runtime)
}

func scanValues[T api.Lambda | api.Runtime](ctx context.Context, prefix string, handler func(x *T) bool) error {
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

func getValues[T api.Lambda | api.Runtime](ctx context.Context, prefix string) ([]*T, error) {
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

func GetLambdas(ctx context.Context) ([]*api.Lambda, error) {
	return getValues[api.Lambda](ctx, "lambda")
}

func GetRuntimes(ctx context.Context) ([]*api.Runtime, error) {
	return getValues[api.Runtime](ctx, "runtime")
}

func getValueByKey[T api.Lambda | api.Runtime](ctx context.Context, key string) (*T, error) {
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

func getValue[T api.Lambda | api.Runtime](ctx context.Context, prefix string, id string) (*T, error) {
	key := prefix + ":" + id
	return getValueByKey[T](ctx, key)
}

func GetLambda(ctx context.Context, id string) (*api.Lambda, error) {
	return getValue[api.Lambda](ctx, "lambda", id)
}

func GetRuntime(ctx context.Context, id string) (*api.Runtime, error) {
	return getValue[api.Runtime](ctx, "runtime", id)
}

func findValue[T api.Lambda | api.Runtime](ctx context.Context, prefix string, predicate func(x *T) bool) (*T, error) {
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

func FindLambda(ctx context.Context, predicate func(val *api.Lambda) bool) (*api.Lambda, error) {
	return findValue(ctx, "lambda", predicate)
}
