package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	secKillConf = &SecKillConf{}
)

type RedisConf struct {
	redisAddr        string
	redisMaxIdle     int
	redisMaxActive   int
	redisIdleTimeout int
}

type EtcdConf struct {
	etcdAddr    string
	etcdTimeout int
	etcdSecKey  string
}

type SecKillConf struct {
	redisConf RedisConf
	etcdConf  EtcdConf
	logPath   string
	logLevel  string
}

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
}

func initConfig() (err error) {
	redisAddr := beego.AppConfig.String("redis_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")

	logs.Debug("read config succ, redis addr:%v", redisAddr)
	logs.Debug("read config succ, etcd addr:%v", etcdAddr)

	secKillConf.redisConf.redisAddr = redisAddr
	secKillConf.etcdConf.etcdAddr = etcdAddr

	if len(redisAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("init config failed,redis[%s] or etcd[%s] config is null", redisAddr, etcdAddr)
		return
	}

	redisMaxIdle, err := beego.AppConfig.Int("redis_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed,read redisMaxIdle %v", err)
		return
	}
	secKillConf.redisConf.redisMaxIdle = redisMaxIdle

	redisMaxActive, err := beego.AppConfig.Int("redis_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed,read redisMaxActive %v", err)
		return
	}
	secKillConf.redisConf.redisMaxActive = redisMaxActive

	redisIdleTimeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed,read redisIdleTimeout %v", err)
		return
	}
	secKillConf.redisConf.redisIdleTimeout = redisIdleTimeout

	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed,read etcdTimeout %v", err)
		return
	}
	secKillConf.etcdConf.etcdTimeout = etcdTimeout
	secKillConf.etcdConf.etcdSecKey = beego.AppConfig.String("etcd_sec_key")
	if len(secKillConf.etcdConf.etcdSecKey) == 0 {
		err = fmt.Errorf("init config failed,read etcd_sec_key error:%v", err)
		return
	}

	secKillConf.logPath = beego.AppConfig.String("log_path")
	secKillConf.logLevel = beego.AppConfig.String("log_level")

	return
}
