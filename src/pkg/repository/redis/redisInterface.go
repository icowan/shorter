/**
 * @Time : 2019-11-15 15:36
 * @Author : solacowa@gmail.com
 * @File : redisInterface
 * @Software: GoLand
 */

package redis

import (
	"strings"
	"time"
)
import "github.com/go-redis/redis"

type RedisInterface interface {
	Set(k string, v interface{}, expir ...time.Duration) (err error)
	Get(k string) (v string, err error)
	Del(k string) (err error)
	HSet(k string, field string, v interface{}) (err error)
	HGet(k string, field string) (res string, err error)
	HGetAll(key string) (map[string]string, error)
	HDelAll(k string) (err error)
	HDel(k string, field string) (err error)
	Close() error
	Subscribe(channels ...string) *redis.PubSub
	Publish(channel string, message interface{}) error
}

type RedisDrive string

const (
	RedisCluster RedisDrive = "cluster"
	RedisSingle  RedisDrive = "single"
	expiration              = 600 * time.Second
)

func NewRedisClient(drive RedisDrive, hosts, password, prefix string, database int) RedisInterface {
	if drive == RedisCluster {
		return NewRedisCluster(strings.Split(hosts, ";"), password, prefix)
	} else {
		return NewRedisSingle(hosts, password, prefix, database)
	}
}
