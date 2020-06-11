package redis

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const Nil = redis.Nil

const RetryCount = 5

type RedisPool struct {
	client *redis.Client
	conf   *domain.DBConf
}

func NewPool(conf *domain.DBConf) *RedisPool {
	if conf.Host == "" {
		panic(errors.New("redis config error"))
	}
	if conf.Port == 0 {
		conf.Port = 6379
	}
	var (
		client *redis.Client
		err    error
	)
	r := &RedisPool{conf: conf}
	client = r.CreateConn()
	if nil == client {
		panic(err)
	}
	r.client = client
	return r
}

func (r *RedisPool) NewConn() *redis.Client {
	return r.client
}

func (r *RedisPool) CreateConn() *redis.Client {
	db := 0
	if r.conf.DBName != "" {
		var err error
		db, err = strconv.Atoi(r.conf.DBName)
		if err != nil {
			fmt.Printf("Invalid DBName:%s", r.conf.DBName)
			panic(err)
		}
	}
	for i := 0; i < RetryCount; i++ {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", r.conf.Host, r.conf.Port),
			Password: "",
			DB:       db,
		})
		_, err := client.Ping().Result()
		if err != nil {
			loggers.Error.Printf("Failed to connect Redis Server: %v error:%s", r.conf, err.Error())
			time.Sleep(2 * time.Second)
			loggers.Warn.Printf("Retrying to connect")
		} else {
			return client
		}
	}
	return nil
}

func (r *RedisPool) IsNil(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}
