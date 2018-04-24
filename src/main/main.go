package main

import (
	"config"
	"dao"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	//使用config包的方法读取etc下的配置文件(json格式)
	err := config.LoadConfig("./etc/conf.json")
	if err != nil {
		log.Printf("加载配置文件异常，原因如下：%v\n", err)
		return
	}
	log.Println("配置文件读取成功！")
	log.Println("使用自定义配置的数据库，连接信息如下：")
	//通过自定义配置文件初始化数据库连接
	msg := dao.InitCustomConnection()
	//msg不为空 则打印输出
	if msg != "" {
		fmt.Printf(msg)
	}

	//静态文件的路由
	http.Handle("/css/", http.FileServer(http.Dir("./webRoot")))
	http.Handle("/images/", http.FileServer(http.Dir("./webRoot")))
	http.Handle("/js/", http.FileServer(http.Dir("./webRoot")))

	//dynamic route
	http.HandleFunc("/", defaultPageHandler)
	http.HandleFunc("/rest", restHandler)
	http.HandleFunc("/info", poolInfoHandler)

	log.Printf("服务通过端口：%v启动！\n", config.ServiceConfig.ServicePort)
	err = http.ListenAndServe(":"+config.ServiceConfig.ServicePort, nil)
	if err != nil {
		log.Printf("http.ListenAndServe error = [%v]\n", err)
		return
	}
	log.Printf("in main(), after http.ListenAndServe().\n")
	return
}
