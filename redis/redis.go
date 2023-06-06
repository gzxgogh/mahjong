package redis

import (
	"fmt"
	"github.com/gzxgogh/ggin/db"
	"strings"
	"time"
)

func SetValue(key, value string, expiration time.Duration) {
	conn := db.DBObj.GetRedisConn()
	conn.Del(key)
	conn.Set(key, value, 24*time.Hour)
}

func DelKey(key string) {
	conn := db.DBObj.GetRedisConn()
	conn.Del(key)
}

func GetValue(key string) string {
	conn := db.DBObj.GetRedisConn()
	res := conn.Get(key)
	str := strings.ReplaceAll(fmt.Sprint(res), " ", "")
	str = strings.ReplaceAll(str, fmt.Sprintf(`get%s:`, key), "")
	if str == "redis:nil" {
		str = ""
	}
	return str
}
