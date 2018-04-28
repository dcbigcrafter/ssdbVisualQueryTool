package dao

import (
	"config"
	"errors"
	"fmt"
	"github.com/seefan/gossdb"
	"github.com/seefan/gossdb/conf"
	"log"
	"strings"
	"sync"
	"util"
)

var (
	//连接池管理map key是"dbip:dbport" value是对应的连接池 加了读写锁
	conPools struct {
		PoolMap map[string]*gossdb.Connectors
		Lock    sync.RWMutex
	}
)

func init() {
	//初始化map
	conPools.PoolMap = make(map[string]*gossdb.Connectors)
}

//用于新建连接池的方法
func getConPool(dbConfig config.DbConf) (*gossdb.Connectors, error) {
	connector, err := gossdb.NewPool(&conf.Config{
		Host:             dbConfig.Host,             //ssdb的ip或主机名
		Port:             dbConfig.Port,             //ssdb的端口
		MinPoolSize:      dbConfig.MinPoolSize,      //最小连接池数
		MaxPoolSize:      dbConfig.MaxPoolSize,      //最大连接池个数
		MaxWaitSize:      dbConfig.MaxWaitSize,      //最大等待数目,当连接池满后，新建连接将等待池中连接释放后才可以继续
		AcquireIncrement: dbConfig.AcquireIncrement, //当连接池中的连接耗尽的时候一次同时获取的连接数
		GetClientTimeout: dbConfig.GetClientTimeout, //获取连接超时时间，单位为秒
		HealthSecond:     dbConfig.HealthSecond,     //连接池内缓存的连接状态检查时间隔，单位为秒
	})
	//如果发生错误
	if err != nil {
		return nil, err
	}
	return connector, nil
}

//初始化自定义配置的数据库连接
func InitCustomConnection() (msg string) {
	//为了后台输出好看 在每一行前加的空格 统一定义便于更改
	var tabSpaces string = "\t\t\t"
	//没有自定义配置的情况
	config.DbCustomConfig.Lock.RLock()                         //读取锁
	dbCustomConfigLen := len(config.DbCustomConfig.CusConfMap) //取得自定义配置的个数
	config.DbCustomConfig.Lock.RUnlock()                       //解锁
	if dbCustomConfigLen == 0 {
		msg = tabSpaces + "当前数据库配置中没有自定义配置，无需初始化连接。"
		return msg
	}
	//当存在自定义配置时 加载自定义配置中的数据库连接
	config.DbCustomConfig.Lock.RLock() //读取锁
	for _, dbCustomConf := range config.DbCustomConfig.CusConfMap {
		//组成该数据库的url
		dbUrl := dbCustomConf.Host + ":" + fmt.Sprintf("%d", dbCustomConf.Port)
		//初始化数据库
		connector, err := getConPool(dbCustomConf)
		if err != nil {
			//出错的情况下忽略该数据库的初始化 打印日志 并继续执行
			errMsg := strings.Replace(err.Error(), "\n", "", -1) //错误信息去换行符处理
			msg += tabSpaces + "#地址为" + dbUrl + "的数据库启动信息为：\n"
			msg += tabSpaces + "    -启动失败，错误日志：" + errMsg + "\n"
			continue
		}
		//无错误就将该连接池保存至map
		conPools.Lock.Lock() //锁住map
		conPools.PoolMap[dbUrl] = connector
		conPools.Lock.Unlock() //解锁map
		//将自定义片配置放入map
		config.DbCustomConfig.CusConfMap[dbUrl] = dbCustomConf
		//输出日志信息
		msg += tabSpaces + "#地址为" + dbUrl + "的数据库启动信息为：\n"
		msg += tabSpaces + "    -启动成功！\n"
	}
	config.DbCustomConfig.Lock.RUnlock() //解锁
	return msg
}

//新建数据库连接池 供其他包调用的
func NewPool(dbIp string, dbPort int) error {
	//组成该数据库的url
	dbUrl := dbIp + ":" + fmt.Sprintf("%d", dbPort)
	//查看map中是否有此数据库的连接池
	conPools.Lock.RLock() //读取锁
	pool, ok := conPools.PoolMap[dbUrl]
	conPools.Lock.RUnlock() //解锁
	//如果存在
	if ok {
		/*
		* 判断数据库是否可以连接 不知道这样判断是否可行
		 */
		//尝试从连接池获取一个连接
		dbClient, err := pool.NewClient()
		if err != nil {
			//有err说明数据库状态有问题
			return err
		} else {
			//没err说明数据库正常
			dbClient.Close()
			return nil
		}
	}
	//不存在的话 需要使用配置文件新建连接池
	var dbConf config.DbConf
	//先看自定义配置里有没有该数据库的配置
	config.DbCustomConfig.Lock.RLock() //上读取锁
	dbCusConf, ok := config.DbCustomConfig.CusConfMap[dbUrl]
	config.DbCustomConfig.Lock.RUnlock() //解锁
	if ok {
		//存在则使用自定义配置
		dbConf = dbCusConf
	} else {
		//不存在则使用默认配置
		dbConf = config.DefaultDbConf
	}
	//给dbConf赋值 数据库ip 端口号
	dbConf.Host = dbIp
	dbConf.Port = dbPort
	//通过默认配置新建连接池
	connectionPool, err := getConPool(dbConf)
	if err != nil {
		//有err说明数据库连接有问题
		return err
	}
	//无错的情况下 将本连接池加入map
	conPools.Lock.Lock() //读写锁
	conPools.PoolMap[dbUrl] = connectionPool
	conPools.Lock.Unlock() //解锁
	return nil
}

//从某个连接池中获取一个数据库连接
func GetConnection(dbIp string, dbPort int) (*gossdb.Client, error) {
	//组成该数据库的url
	dbUrl := dbIp + ":" + fmt.Sprintf("%d", dbPort)
	//查看该数据库的连接池是否存在
	conPools.Lock.RLock() //读锁
	pool, ok := conPools.PoolMap[dbUrl]
	conPools.Lock.RUnlock() //解锁
	if !ok {
		log.Printf("从连接池获取连接时，发现%v的连接池没建立\n", dbUrl)
		return nil, errors.New("请先连接数据库")
	}
	//尝试获取连接
	c, err := pool.NewClient()
	if err != nil {
		return nil, err
	}
	//测试连接
	if !c.Ping() {
		c.Close() //关闭连接
		return nil, errors.New("数据库连接失败")
	}
	return c, nil
}

//查看已经连接的所有连接池状态
func StatusOfPools() string {
	//返回的信息
	var msg string
	//判断连接池管理map里面有没有连接
	conPools.Lock.RLock()                //map上锁
	conPoolsLen := len(conPools.PoolMap) //获取连接池的个数
	if conPoolsLen == 0 {
		msg = "当前暂无数据库连接"
		conPools.Lock.RUnlock() //解锁
		return msg
	}
	var dbUrls []string                 //数据库连接名称列表
	dburlMap := make(map[string]string) //数据库连接与对应的连接池信息
	//存在数据库连接则逐一显示信息
	for k, v := range conPools.PoolMap {
		//保存数据库连接名称
		dbUrls = append(dbUrls, k)
		msg = "地址为" + k + "的数据库连接池信息为：\n"
		//先尝试获取一个连接
		c, err := v.NewClient()
		if err != nil {
			msg += "\t--无法正常获取数据库连接，数据库或许已经关闭，异常信息：" + err.Error()
			continue
		}
		//ping数据库 看能否查询数据
		ok := c.Ping()
		if ok {
			msg += "\t--能够正常查询数据\n"
		} else {
			msg += "\t--无法查询数据\n"
		}
		c.Close() //关闭连接
		msg += "\t--" + v.Info() + "\n\n"
		//保存数据库连接名称与该数据库连接池信息
		dburlMap[k] = msg
	}
	conPools.Lock.RUnlock() //解锁
	//排序数据库连接名称列表
	util.SortStringArray(&dbUrls)
	//将排序后的信息整合
	msg = "" //清空
	//遍历dburlMap 取出对应的连接池信息
	for _, dbUrl := range dbUrls {
		msg += dburlMap[dbUrl]
	}
	return msg
}
