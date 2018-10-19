package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	etcd_client "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"SecKill/SecProxy/service"
)

func initEtcd()(err error){

	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	etcdClient = cli
	return
	
}
var (
	redisPool *redis.Pool
	etcdClient *etcd_client.Client
)

func initRedis()(err error){
	redisPool = &redis.Pool{
		MaxIdle: secKillConf.RedisConf.RedisMaxIdle,
		MaxActive: secKillConf.RedisConf.RedisMaxActive,
		IdleTimeout:time.Duration(secKillConf.RedisConf.RedisIdleTimeout)*time.Second,
		Dial: func()(redis.Conn,error){
          return redis.Dial("tcp",secKillConf.RedisConf.RedisAddr)
		},
	}

	conn := redisPool.Get()
	defer conn.Close()

   _,err =	conn.Do("ping")
   if err != nil{
	   logs.Error("ping redis failed,err:%v",err)
   }
	return
}

//log级别的转换

func convertLogLevel(level string) int{
   switch(level){
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

func initLogger()(err error){
	config := make(map[string]interface{})
	config["filename"] = secKillConf.LogPath
	config["level"] = convertLogLevel(secKillConf.LogLevel)

	configStr,err := json.Marshal(config)
	if err != nil{
		fmt.Println("Marshal failed,err",err)
	}

	logs.SetLogger(logs.AdapterFile,string(configStr))
	return 
}

func loadSecConf()(err error){
 
 resp,err :=etcdClient.Get(context.Background(),secKillConf.EtcdConf.EtcdSecProductKey)	
 if err != nil{
	 logs.Error("get [%s] from etcd failed,err:%v",secKillConf.EtcdConf.EtcdSecProductKey,err)
	 return
 }
 //把返回的json转换成对象
 var secProductInfo []service.SecProductInfoConf
 for k,v := range resp.Kvs{
	 logs.Debug("key[%s] values[%s]",k,v)
	 err = json.Unmarshal(v.Value,&secProductInfo)
	 if err != nil{
		 logs.Error("Unmarshal sec product info failed,err:%v",err)
	     return 
	}
	logs.Debug("sec info config is [%v]",secProductInfo)
 }
 secKillConf.RwSecProductLock.Lock()
 for _,v := range secProductInfo{
	 secKillConf.SecProductInfoMap[v.ProductId] = &v
 }
 secKillConf.RwSecProductLock.Unlock()

  return
}

func initSec()(err error){

	 err = initLogger()

	 if err != nil{
		 logs.Error("init logger failed ,err:%v",err)
		 return
	 }

	//etcd
	err = initEtcd()
	if err != nil{
		logs.Error("init etcd failed,err:%v",err)
		return
	}

	//redis
	err = initRedis()
	if err != nil{
		logs.Error("init redis failed,err:%v",err)
		return
	}

	err = loadSecConf()

	if err != nil{
		logs.Error("init conf failed,err:%v",err)
		return
	}

	service.InitService(secKillConf)
		
	initSecProductWatcher()

	logs.Info("init sec Succ")
	 


	return


}
//监听etcd的变化
func initSecProductWatcher(){
	go watchSecProductKey(secKillConf.EtcdConf.EtcdSecProductKey)
}

func watchSecProductKey(key string){
   cli,err := etcd_client.New(etcd_client.Config{
	   Endpoints: []string{"127.0.0.1:2379"},
	   DialTimeout: 5*time.Second,

   })

   if err != nil{
	   logs.Error("connect etcd failed,err",err)
       return
	}

  logs.Debug("begin watch key:%s",key)
   
  for{
	  rch := cli.Watch(context.Background(),key)
	  var secProductInfo []service.SecProductInfoConf
	  var getConfSucc = true

	  for wresp := range rch{
		   for _,ev := range wresp.Events{
			   if ev.Type == mvccpb.DELETE{
				   logs.Warn("key[%s] config deleted",key)
				   continue
			   }

			   if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key{
				   err = json.Unmarshal(ev.Kv.Value,&secProductInfo)
				   if err != nil{
					   logs.Error("key ,unmarshal,err:%v",err)
					   getConfSucc = false
					   continue
				   }
			   }
			   logs.Debug("get config from etcd,%s %q:%q\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
		   }

		   if getConfSucc{
			   logs.Debug("get config from etcd succ,%v",secProductInfo)
			   updateSecProductInfo(secProductInfo)
		   }
	  }
  }
  
}

func updateSecProductInfo(secProductInfo []service.SecProductInfoConf){
	
	var tmp map[int]*service.SecProductInfoConf = make(map[int]*service.SecProductInfoConf,1024)
	for _,v := range secProductInfo{
		tmp[v.ProductId] = &v
	}

	secKillConf.RwSecProductLock.Lock()

		secKillConf.SecProductInfoMap = tmp
	
	secKillConf.RwSecProductLock.Unlock()
}
