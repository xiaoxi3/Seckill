package service

import (
	"github.com/astaxie/beego/logs"
)


var (
	secKillConf *SecKillConf
)

func InitService(serviceConf *SecKillConf){
	secKillConf = serviceConf
	logs.Debug("init service succ,config:%v",secKillConf)
}