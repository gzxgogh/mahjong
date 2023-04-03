package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

var DB *redis.Client

func InItRedisCoon() {
	rdb := redis.NewClient(&redis.Options{
		// 需要修改成你的配置，本地无需修改
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	DB = rdb
}

func SetValue(key, value string, expiration time.Duration) {
	DB.Del(key)
	DB.Set(key, value, 24*time.Hour)
}

func DelKey(key string) {
	DB.Del(key)
}

func GetValue(key string) string {
	res := DB.Get(key)
	str := strings.ReplaceAll(fmt.Sprint(res), " ", "")
	str = strings.ReplaceAll(str, fmt.Sprintf(`get%s:`, key), "")
	if str == "redis:nil" {
		str = ""
	}
	return str
}
