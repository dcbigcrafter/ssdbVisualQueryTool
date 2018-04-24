package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
)

var (
	//自定义配置管理map key是"dbip:dbport" value是对应的配置 加了读写锁
	DbCustomConfig struct {
		CusConfMap map[string]DbConf
		Lock       sync.RWMutex
	}
	//数据库默认配置
	DefaultDbConf DbConf
	//rest服务配置
	ServiceConfig ServiceConf
)

func init() {
	//初始化map
	DbCustomConfig.CusConfMap = make(map[string]DbConf)
}

//所有配置的集合 用于接收文件中的json配置信息
type ConfigSet struct {
	Service ServiceConf `json:"service"` //rest服务的配置
	//Cleaner   CleanerConf `json:"cleaner"`   //session清理程序的配置
	DbDefault DbConf   `json:"dbDefault"` //数据库连接的默认配置
	DbCustom  []DbConf `json:"dbCustom"`  //针对特定数据库连接的自定义配置
}

//rest服务相关的详细配置
type ServiceConf struct {
	ServicePort string `json:"servicePort"` //rest服务的服务端口
}

//session清理程序的详细配置
type CleanerConf struct {
	ExpireTime int64 `json:"expireTime"` //session的过期失效时间 单位为秒
	FirstStart int64 `json:"firstStart"` //清理程序首次经过多少秒之后执行 单位为秒
}

//数据库连接的详细配置
type DbConf struct {
	Host             string `json:"host"`             //ssdb的ip或主机名
	Port             int    `json:"port"`             // ssdb的端口
	GetClientTimeout int    `json:"getClientTimeout"` //获取连接超时时间，单位为秒。默认值: 5
	MaxPoolSize      int    `json:"maxPoolSize"`      //最大连接池个数。默认值: 20
	MinPoolSize      int    `json:"minPoolSize"`      //最小连接池数。默认值: 5
	AcquireIncrement int    `json:"acquireIncrement"` //当连接池中的连接耗尽的时候一次同时获取的连接数。默认值: 5
	MaxWaitSize      int    `json:"maxWaitSize"`      //最大等待数目，当连接池满后，新建连接将等待池中连接释放后才可以继续，本值限制最大等待的数量，超过本值后将抛出异常。默认值: 1000
	HealthSecond     int    `json:"healthSecond"`     //连接池内缓存的连接状态检查时间隔，单位为秒。默认值: 5
	//Password         string `json:"password"`             //连接的密钥
	//Weight           int    `json:"weight"`             //权重，只在负载均衡模式下启用
}

//从json文件加载配置信息的方法
func LoadConfig(filename string) error {
	//读取json配置文件内容
	fileContentByte, err := ioutil.ReadFile(filename)
	//出现错误则返回错误
	if err != nil {
		//组成错误提示信息
		warnMsg := "读取配置文件时发生错误，错误原因为：" + err.Error()
		return errors.New(warnMsg)
	}
	//无错的情况下 将json信息unmarshal
	var conf ConfigSet
	err = json.Unmarshal(fileContentByte, &conf)
	if err != nil {
		//组成错误提示信息
		warnMsg := "解析配置文件时发生错误，错误原因为：" + err.Error()
		return errors.New(warnMsg)
	}
	//保存服务配置
	ServiceConfig = conf.Service
	//保存默认配置
	DefaultDbConf = conf.DbDefault
	//遍历自定义配置 并存入map
	DbCustomConfig.Lock.Lock() //上读写锁
	for _, dbCusConf := range conf.DbCustom {
		//组成该数据库的url
		dbUrl := dbCusConf.Host + ":" + fmt.Sprintf("%d", dbCusConf.Port)
		DbCustomConfig.CusConfMap[dbUrl] = dbCusConf //存入map
	}
	DbCustomConfig.Lock.Unlock() //解锁
	return nil
}
