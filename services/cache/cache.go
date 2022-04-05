package cache

import (
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type CacheClient struct {
	cache *redis.Client
}

func NewClient() (*CacheClient, error) {

	address := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")

	cache := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	// check if the redis is available
	_, err := cache.Ping().Result()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not connect to cache")
		return nil, err
	}

	return &CacheClient{cache: cache}, nil
}

func (c *CacheClient) Write(key, val string) error {

	if _, err := c.cache.Set(key, val, 0).Result(); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not write to cache")
		return err
	}

	return nil
}

func (c *CacheClient) Read(key string) (string, error) {

	value, err := c.cache.Get(key).Result()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not read from cache")
		return "", err
	}
	return value, nil

}
