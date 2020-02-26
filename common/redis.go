package common

import (
	"github.com/go-redis/redis"
	"io/ioutil"
	"path"
	"seckilling-practice-project/configs"
)

var luaDecrby *redis.Script
var redisClient *redis.Client

func GetDecrbyScr() (*redis.Script, error) {
	if luaDecrby == nil {
		rootPath := configs.GetProjectPath()
		scr, err := ioutil.ReadFile(path.Join(rootPath, "scripts", "decrby.lua"))
		if err != nil {
			return &redis.Script{}, nil
		}
		luaDecrby = redis.NewScript(string(scr))
	}
	return luaDecrby, nil
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		config := configs.Cfg
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Address,
			Password: "", // no password set
			DB:       0,  // use default DB
			PoolSize: 1000,
		})
	}
	return redisClient
}
