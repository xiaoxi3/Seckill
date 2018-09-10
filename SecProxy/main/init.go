package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	etcd_client "go.etcd.io/etcd/clientv3"
)

func initEtcd() (err error) {
	cli, err := etcd_client.New(etcd_client.Config{

		Endpoints:            []string{secKillConf.etcdConf.etcdAddr},
		DialKeepAliveTimeout: time.Duration(secKillConf.etcdConf.etcdTimeout) * time.Second,
	})

	if err != nil {
		logs.Error("connect etcd failed,err", err)
		return
	}

	etcdClient = cli

	return
}

var (
	redisPool  *redis.Pool
	etcdClient *etcd_client.Client
)

func initRedis() (err error) {
	redisPool = &redis.Pool{
		MaxIdle:     secKillConf.redisConf.redisMaxIdle,
		MaxActive:   secKillConf.redisConf.redisMaxActive,
		IdleTimeout: time.Duration(secKillConf.redisConf.redisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.redisConf.redisAddr)
		},
	}

	conn := redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failer %v", err)
	}
	return
}

//log级别转换

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = secKillConf.logPath
	config["level"] = convertLogLevel(secKillConf.logLevel)

	configStr,err := json.Marshal(config)
	if err != nil{
		fmt.Println("Marshal failed,err",err)
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

func loadSecConf() (err error) {
	key := fmt.Sprintf("%s/product", secKillConf.etcdConf.etcdSecKey)
	resp, err := etcdClient.Get(context.Background(), key)
	if err != nil {
		logs.Error("get [%s] from etcd failed,err:%v", key, err)
		return
	}
	for k, v := range resp.Kvs {
		logs.Debug("key[%s] values[%s]", k, v)
	}
	return
}

func initSec() (err error) {
	//log
	err = initLogger()
	if err != nil {
		logs.Error("init logger failed,err:%v", err)
		return
	}

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

	err = loadSecConf()

	logs.Info("init sec Succ")
	return
}
