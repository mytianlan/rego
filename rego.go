package rego

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

const timeout = 30
const MaxIdle = 20
const MaxActive = 50
const MaxIdleTimeout = 180

func Connect(conf map[string]interface{}) (redisConn redis.Conn, err error) {
	//check conf
	if _, ok := conf["Host"]; ok == false {
		return nil, errors.New("conf error")
	}
	if _, ok := conf["Password"]; ok == false {
		return nil, errors.New("conf error")
	}
	if _, ok := conf["Db"]; ok == false {
		return nil, errors.New("conf error")
	}
	con, err := redis.Dial("tcp", conf["Host"].(string),
		redis.DialPassword(conf["Password"].(string)),
		redis.DialDatabase(int(conf["Db"].(int64))),
		redis.DialConnectTimeout(timeout*time.Second),
		redis.DialReadTimeout(timeout*time.Second),
		redis.DialWriteTimeout(timeout*time.Second))
	if err != nil {
		return nil, err
	}

	return con, nil
}

func SaddIntSlice(redisConn redis.Conn, key string, slice interface{}) error {
	stringSlice := slice.([]uint64)
	for _, v := range stringSlice {
		redisConn.Do("SADD", key, v)
	}
	return nil
}

func SaddStringSlice(redisConn redis.Conn, key string, slice interface{}) error {
	stringSlice := slice.([]string)
	for _, v := range stringSlice {
		redisConn.Do("SADD", key, v)
	}
	return nil
}

func GetConnectionPool(conf map[string]interface{}) (*redis.Pool, error) {
	maxIdle := MaxIdle
	if v, ok := conf["MaxIdle"]; ok {
		maxIdle = int(v.(int64))
	}
	maxActive := MaxActive
	if v, ok := conf["MaxActive"]; ok {
		maxActive = int(v.(int64))
	}

	// 建立连接池
	redisClient := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: MaxIdleTimeout * time.Second,
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", conf["Host"].(string),
				redis.DialPassword(conf["Password"].(string)),
				redis.DialDatabase(int(conf["Db"].(int64))),
				redis.DialConnectTimeout(timeout*time.Second),
				redis.DialReadTimeout(timeout*time.Second),
				redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
	return redisClient, nil
}
