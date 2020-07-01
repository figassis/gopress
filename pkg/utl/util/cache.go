package util

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"
	"time"

	redis "gopkg.in/redis.v3"
)

const cacheTTL = time.Hour

var (
	cache    *redis.Client
	cacheKey string
)

func loadRedis(namespace, host, port string, db int64) (err error) {
	cacheKey = namespace
	cache = redis.NewClient(&redis.Options{Addr: host + ":" + port, Password: "", DB: db})
	return cache.Ping().Err()
}

func GetCache(key string, result interface{}) (err error) {
	if cache == nil {
		return errors.New("Redis not ready")
	}
	if key == "" {
		return errors.New("Cache key cannot be empty")
	}
	key = cacheKey + "_" + key
	//logging.BLog(1, fmt.Errorf("KEY: %s", key), tid)
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return errors.New("Please provide a pointer")
	}

	obj, err := cache.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.New("Record not found")
		}
		return
	}

	switch string(obj) {
	case "[]", "{}", "":
		return errors.New("Cache miss")
	}

	if err = json.Unmarshal(obj, result); err != nil {
		return
	}

	go cache.Set(key, obj, cacheTTL)
	return
}

func DeleteCache(key string) (err error) {
	if cache == nil {
		return errors.New("Redis not ready")
	}

	key = cacheKey + "_" + key
	_, err = cache.Del(key).Result()
	if err == redis.Nil {
		return errors.New("Key not found")
	} else if err != nil {
		return errors.New("Could not remove key")
	}

	return
}

func Cache(key string, obj interface{}) (err error) {
	if cache == nil {
		return errors.New("Redis not ready")
	}

	if key == "" {
		return errors.New("Key cannot be empty")
	}
	key = cacheKey + "_" + key

	if obj == nil {
		return
	}

	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return cache.Set(key, jsonObj, cacheTTL).Err()
}

func CacheTTL(key string, obj interface{}, ttl time.Duration) (err error) {
	if cache == nil {
		return errors.New("Redis not ready")
	}

	key = cacheKey + "_" + key

	if obj == nil {
		return
	}

	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return cache.Set(key, jsonObj, ttl).Err()
}

func New() (err error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	if err = loadRedis(os.Getenv("FQDN"), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), int64(redisDB)); err != nil {
		return
	}

	return
}
