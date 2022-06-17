package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"strconv"
)

var redisConfigMap map[string]interface{}

var (
	Redis  *redis.Client
	ARedis *redis.Client
)

func InitRedis() {
	initRedis("master", &Redis)
	initRedis("activity", &ARedis)
}

func initRedis(node string, rdb **redis.Client) {
	redisConfigMap = make(map[string]interface{}, 4)
	redisConfigMap = viper.GetStringMap("redis." + node)
	*rdb = redis.NewClient(&redis.Options{
		Addr:     redisConfigMap["host"].(string) + ":" + strconv.FormatInt(redisConfigMap["port"].(int64), 10),
		Password: redisConfigMap["password"].(string),
		DB:       int(redisConfigMap["database"].(int64)),
	})
}

func getRedisUrl() string {
	return mysqlConfigMap["user"].(string) + ":" + mysqlConfigMap["password"].(string) + "@tcp(" + mysqlConfigMap["host"].(string) + ":" + strconv.FormatInt(mysqlConfigMap["port"].(int64), 10) + ")/" + mysqlConfigMap["name"].(string)
}
