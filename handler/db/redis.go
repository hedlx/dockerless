package db

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

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
		if err := rdb.Ping(context.Background()).Err(); err == nil {
			break
		}

		time.Sleep(time.Second)
	}
}

func scanValues[T any](ctx context.Context, prefix string, handler func(x *T) bool) error {
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

func GetValues[T any](ctx context.Context, prefix string) ([]*T, error) {
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

func getValueByKey[T any](ctx context.Context, key string) (*T, error) {
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

func GetValue[T any](ctx context.Context, prefix string, id string) (*T, error) {
	return getValueByKey[T](ctx, prefix+":"+id)
}

type SetNotification[T any] struct {
	Value *T
}

type DelNotification struct {
	Key string
}

func Subscribe[T any](ctx context.Context, prefix string) <-chan interface{} {
	topic := "__keyspace@0__:" + prefix + "*"
	logger.L.Info("Subscribe to topic", zap.String("topic", topic))
	pubsub := rdb.PSubscribe(ctx, topic)
	notificationsC := make(chan interface{})

	go func() {
		defer close(notificationsC)
		logger.L.Info("Wait for notifications")

		for msg := range pubsub.Channel() {
			logger.L.Info("New notification", zap.String("channel", msg.Channel))
			tSlice := strings.SplitN(msg.Channel, ":", 2)
			if len(tSlice) < 2 {
				logger.L.Error(
					"Unexpected notification",
					zap.String("msg", msg.Channel),
				)
				continue
			}

			t := msg.Payload
			if t == "del" {
				notificationsC <- &DelNotification{msg.Payload}
				continue
			}

			val, err := getValueByKey[T](ctx, tSlice[1])
			if err != nil {
				logger.L.Error(
					"Failed to get redis value",
					zap.Error(err),
					zap.String("key", msg.Payload),
				)
				continue
			}

			notificationsC <- &SetNotification[T]{val}
		}
	}()

	return notificationsC
}
