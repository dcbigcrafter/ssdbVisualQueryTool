package main

import (
	"pojo"
	"strconv"
	"strings"
)

//com可能值的集合
var comMap = map[string]string{
	"QUERY": "QUERY", //查询
	//"POST":    "POST",    //增加用的查询接口
	"DELETE":  "DELETE",  //删除
	"CONNECT": "CONNECT", //连接
}

//对数据进行校验的方法
func CheckRqtData(param *pojo.RequestJson) string {
	//定义一个msg用于保存提示语
	var msg string
	//将com转为大写
	param.Com = strings.ToUpper(param.Com)
	//先判断com是否在com集合中
	if _, ok := comMap[param.Com]; !ok {
		msg = "com：" + param.Com + "，命令无法识别，无法继续操作！"
		return msg
	}
	//不管什么操作 dbIp和dbPort都是必须有的
	if param.Data.DbIP == "" {
		msg += "数据库ip地址不能为空."
	}
	if param.Data.DbPort == "" {
		msg += "数据库端口号不能为空."
	}
	//端口号不为空时 将其转换为int
	if param.Data.DbPort != "" {
		//将数据库端口号转换为数字
		dbPort, err := strconv.Atoi(param.Data.DbPort)
		if err != nil {
			msg += "不合法的端口号：" + param.Data.DbPort + "，端口号必须为数字."
		} else {
			//转换为数字正常 则存储
			param.Data.DbPortInt = dbPort
		}
	}
	//删除的话 key不能为空
	if param.Com == "DELETE" && param.Data.Key == "" {
		msg += "请勾选需要删除的数据！"
	}
	return msg
}
