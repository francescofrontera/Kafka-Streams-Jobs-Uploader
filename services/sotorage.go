package services

import (
	"github.com/go-redis/redis"
	"log"
)

type RedisConf struct {
	cli *redis.Client
	log *log.Logger
}

func RedisClient(addr string, pwd string, db int, logger *log.Logger) *RedisConf {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})

	_, err := cli.Ping().Result(); if err != nil {
		logger.Fatalf("An error occured during init of DB: %v", err)
	}

	return &RedisConf{
		cli: cli,
		log: logger,
	}
}
