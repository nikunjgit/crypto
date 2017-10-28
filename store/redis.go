package store

import (
	"github.com/go-redis/redis"
	"encoding/json"
	"time"
	"github.com/nikunjgit/crypto/event"
)

type RedisClient struct {
	Client        *redis.Client
	ttl time.Duration
}

func NewRedisClient(ttl time.Duration) (*RedisClient, error) {
	rClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rClient.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{rClient, ttl}, nil
}

func (m *RedisClient) Set(key string, messages event.Messages) error {
	b, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	err = m.Client.Set(key, string(b), m.ttl).Err()
	return err
}

func (m *RedisClient) Get(keys []string) (event.Messages, error) {
	vals, err := m.Client.MGet(keys...).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	msgArray := make(event.Messages, 0, 10)
	for _, val := range vals {
		str, ok := val.(string)
		if !ok {
			continue
		}
		var messages event.Messages
		if err := json.Unmarshal([]byte(str), &messages); err != nil {
			return nil, err
		}
		msgArray = append(msgArray, messages...)
	}

	return msgArray, nil
}
