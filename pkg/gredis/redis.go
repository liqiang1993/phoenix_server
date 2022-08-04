package gredis

import (
	"github.com/lucky-cheerful-man/phoenix_server/pkg/log"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/setting"

	"time"

	"github.com/gomodule/redigo/redis"
)

type CacheInterface interface {
	Set(key string, data []byte, time int) error
	Get(key string) ([]byte, error)
	Delete(key string) (bool, error)
}

// Setup Initialize the Redis instance
func Setup() *RedisOperate {
	op := new(RedisOperate)
	Conn := &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	op.Conn = Conn

	return op
}

type RedisOperate struct {
	Conn *redis.Pool
}

// Set 设置键值对
func (r *RedisOperate) Set(key string, value []byte, time int) error {
	conn := r.Conn.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Warnf("close redis conn failed:%s", err)
		}
	}()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Get 根据键查询对应的值
func (r *RedisOperate) Get(key string) ([]byte, error) {
	conn := r.Conn.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Warnf("close redis conn failed:%s", err)
		}
	}()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete 删除指定键
func (r *RedisOperate) Delete(key string) (bool, error) {
	conn := r.Conn.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Warnf("close redis conn failed:%s", err)
		}
	}()

	return redis.Bool(conn.Do("DEL", key))
}
