package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

func initEtcd() (err error) {
	return
}

func initRedis() (err error) {
	_ = &redis.Pool{
		MaxIdle:     secKillConf.redisConf.redisMaxIdle,
		MaxActive:   secKillConf.redisConf.redisMaxActive,
		IdleTimeout: time.Duration(secKillConf.redisConf.redisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.redisConf.redisAddr)
		},
	}
	return
}

func initSec() (err error) {
	//etcd
	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed,err:%v", err)
		return
	}

	//redis
	err = initRedis()
	if err != nil {
		logs.Error("init redis failed,err:%v", err)
		return
	}

	logs.Info("init sec Succ")
	return
}
