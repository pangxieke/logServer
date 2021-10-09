package storage

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"reflect"
	"strings"
	"time"
)

var (
	RedisCluster *redis.ClusterClient
	CacheClient  *Cache
	CacheExpire  = time.Second * 125
	keyPrefix    = "logServer_"
)

func InitRedis(cfg config.Configer) error {
	addrs := cfg.String("redis::addrs")
	passwd := cfg.String("redis::password")
	poolsize := cfg.DefaultInt("redis::poolsize", 10)

	if addrs == "" {
		return fmt.Errorf("redis::addrs is not found")
	}

	opt := redis.ClusterOptions{
		Addrs:    strings.Split(addrs, ","),
		Password: passwd,
		PoolSize: poolsize,
	}
	RedisCluster = redis.NewClusterClient(&opt)

	logs.Info("redis, addrs = %s, passwd = %s, poolsize = %d", addrs, passwd, poolsize)

	if CacheClient == nil {
		CacheClient = NewCache()
	}
	return nil
}

type Cache struct {
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Set(key string, val interface{}, expire time.Duration) error {
	err := RedisCluster.Set(key, val, expire).Err()
	if err != nil {
		logs.Warn("redis set error, key=%s, err=%s,val=%v,", key, err, val)
	}
	return err
}

func (c *Cache) Get(key string) string {
	val, err := RedisCluster.Get(key).Result()
	if err != nil {
		//logs.Debug("redis get %s err=%s", key, err)
	}
	return val
}

func (c *Cache) GetObj(key string, obj interface{}) (err error) {
	val := CacheClient.Get(key)

	if val != "" {
		err := json.Unmarshal([]byte(val), &obj)
		if err != nil {
			logs.Info("json.Unmarshal, error=%s", err)
		}
	}
	return
}

func (c *Cache) SetObj(key string, obj interface{}) (err error) {
	str, err := json.Marshal(obj)
	if err != nil {
		logs.Info("json.Marshal error=%s, key=%s, obj=%s", err, key, obj)
	}
	err = CacheClient.Set(key, str, CacheExpire)
	return
}

// sprintf key
func (c *Cache) CreateKey(format string, a ...interface{}) string {
	str := keyPrefix + fmt.Sprintf(format, a...)
	return str
}

func (c *Cache) ParamsKey(params ...interface{}) string {
	key := ""
	for _, param := range params {
		t := reflect.TypeOf(param)
		switch t.Kind() {
		case reflect.Int:
			key += fmt.Sprintf("_%d", param)
		case reflect.String:
			key += fmt.Sprintf("_%s", param)
		case reflect.Struct:
			if param != nil {
				str, _ := json.Marshal(param)
				key += fmt.Sprintf("_%s", str)
			}
		case reflect.Func:
			//key += fmt.Sprintf("%s", param)
			//TODO
		case reflect.Array:
			key += fmt.Sprintf("_%s", param)
		case reflect.Slice:
			key += fmt.Sprintf("_%s", param)
		case reflect.Map:
			key += fmt.Sprintf("_%s", param)

		}
	}
	return key
}

func (c *Cache) HashKey(params ...interface{}) string {
	key := c.ParamsKey(params...)
	return fmt.Sprintf("hash_key_%s", key)
}

func (c *Cache) Push(key string, obj interface{}) {
	str, _ := json.Marshal(obj)

	RedisCluster.LPush(key, str)
}

func (c *Cache) Pop(key string) string {
	val, err := RedisCluster.LPop(key).Result()
	if err != nil {

	}
	return val
}

func (c *Cache) LLen(key string) int64 {
	val, err := RedisCluster.LLen(key).Result()
	if err != nil {

	}
	return val
}
