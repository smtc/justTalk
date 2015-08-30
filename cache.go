package main

// 2015-08-22
// 直接使用redis

import (
	"encoding/json"

	"github.com/fzzy/radix/extra/pool"
	"github.com/guotie/deferinit"
)

const (
	defaultAddr     = "localhost:6379"
	defaultNetwork  = "tcp"
	defaultPoolSize = 10
)

var (
	_pool *pool.Pool
)

func init() {
	deferinit.AddInit(func() {
		createRedisPool("")
	}, nil, 100)
}

func getFromCache(key string) (reply []byte) {
	reply, _ = GET(key)
	return
}

func getFromCacheString(key string) string {
	var (
		c   cachedUserData
		err error
	)

	reply := getFromCache(key)
	err = json.Unmarshal(reply, &c)
	if err != nil {
		return ""
	}
	return c.Id
}

// options sample:
//   `{  "addr": "127.0.0.1:6389",
//       "network":"tcp",
//       "db": 0,
//       "password": "",
//       "pools": 5
//    }`
func createRedisPool(options string) {
	var (
		err    error
		config struct {
			Addr     string
			Db       int
			Network  string
			Password string
			Pools    int
		}
	)

	if options != "" {
		err = json.Unmarshal([]byte(options), &config)
		if err != nil {
			//println("unmarshal failed:", err.Error())
			config.Addr = defaultAddr
			config.Network = defaultNetwork
			config.Pools = defaultPoolSize
		}
	}

	if config.Pools <= 0 {
		config.Pools = defaultPoolSize
	}
	if config.Addr == "" {
		config.Addr = defaultAddr
	}
	if config.Network == "" {
		config.Network = "tcp"
	}

	_pool, err = pool.NewPool(config.Network, config.Addr, config.Pools)
	if err != nil {
		panic(err.Error())
	}

}

// redis command:
// set key val
func SET(key string, val []byte) error {
	c, err := _pool.Get()
	if err != nil {
		return err
	}
	err = c.Cmd("SET", key, val).Err
	_pool.Put(c)
	return err
}

// redis command:
// setex key seconds val
func SETEX(key string, secs int, val []byte) error {
	c, err := _pool.Get()
	if err != nil {
		return err
	}
	err = c.Cmd("SETEX", key, secs, val).Err
	_pool.Put(c)
	return err
}

// redis command:
// get key
func GET(key string) (val []byte, err error) {
	c, err := _pool.Get()
	if err != nil {
		return
	}
	val, err = c.Cmd("GET", key).Bytes()
	_pool.Put(c)
	return
}

func DEL(key string) {
	c, err := _pool.Get()
	if err != nil {
		return
	}

	c.Cmd("DEL", key)
	_pool.Put(c)
}

// redis command:
// hset key field val
func HSET(key, field string, val []byte) error {
	c, err := _pool.Get()
	if err != nil {
		return err
	}
	err = c.Cmd("HSET", key, field, val).Err
	_pool.Put(c)
	return err
}

// redis command:
// HSETNX key field val timeout
func HSETNX(key, field string, val []byte) (err error) {
	c, err := _pool.Get()
	if err != nil {
		return
	}
	err = c.Cmd("HSETNX", key, field, val).Err
	_pool.Put(c)
	return
}

// redis command:
// HGET key field
func HGET(key, field string) (val []byte, err error) {
	c, err := _pool.Get()
	if err != nil {
		return
	}
	val, err = c.Cmd("HGET", key, field).Bytes()
	_pool.Put(c)
	return
}

// redis command:
// HDEL key field
func HDEL(key, field string) {
	c, err := _pool.Get()
	if err != nil {
		return
	}
	_ = c.Cmd("HDEL", key, field).Err
	_pool.Put(c)
}

func FlushAll() {
	c, err := _pool.Get()
	if err != nil {
		return
	}
	_ = c.Cmd("FLUSHALL").Err
	_pool.Put(c)
}
