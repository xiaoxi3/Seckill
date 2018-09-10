package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
}

const (
	etcdSecKey = "/xiaoxi/seckill/product"
)

func SetLogConfToEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed,err", err)
	}
	fmt.Println("connect succ")
	defer cli.Close()

	var SecInfoConfArr []SecInfoConf

	SecInfoConfArr = append(
		SecInfoConfArr,
		SecInfoConf{
			ProductId: 1022,
			StartTime: 1536570579,
			EndTime:   1536725271,
			Status:    0,
			Total:     10000,
			Left:      10000,
		},
	)

	data, err := json.Marshal(SecInfoConfArr)
	if err != nil {
		fmt.Println("json failed", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, etcdSecKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, etcdSecKey)
	cancel()

	if err != nil {
		fmt.Println("get failed err", err)
		return
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s=>%s\n", ev.Key, ev.Value)
	}
}

func main() {
	SetLogConfToEtcd()
}
