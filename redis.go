package husky

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Prefix   string `toml:"prefix"`
}

type _Redis struct {
	*redis.Client
	Prefix string
}

var redisIns map[string]*_Redis

func init() {
	redisIns = make(map[string]*_Redis)
}

func InitRedis(config *RedisConfig, key ...string) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", config.Addr, config.Port),
		Password: config.Password,
	})
	r := _Redis{
		Client: client,
		Prefix: config.Prefix,
	}
	if len(key) == 0 {
		redisIns[""] = &r
	} else {
		redisIns[key[0]] = &r
	}
}

func Redis(key ...string) *_Redis {
	if len(key) == 0 {
		return redisIns[""]
	} else {
		return redisIns[key[0]]
	}
}
