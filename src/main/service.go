package main

import (
	"dao"
	"fmt"
	"log"
	"pojo"
	"strconv"
	"strings"
	"util"
)

//查询操作的rest服务
func queryHandler(param pojo.RequestJson) util.ResponseJson {
	log.Println("in queryHandler()")
	//将查询数目转换为数字
	queryNum, err := strconv.Atoi(param.Data.Number)
	if err != nil {
		//出错的话设置为100
		queryNum = 100
	}
	//获取数据库连接
	dbClient, err := dao.GetConnection(param.Data.DbIP, param.Data.DbPortInt)
	if err != nil {
		//出错的话提示
		response := util.SystemErr()
		response.Msg = "，获取数据库连接失败，失败原因：" + err.Error()
		return response
	}
	//无错的情况下才需要关闭连接
	defer dbClient.Close()
	//对于新增数据的处理 当key不为空 且value不为空不为默认值时 是新增
	if param.Data.Key != "" && param.Data.Value != "" && param.Data.Value != "请输入Value" {
		//对待存储key做去空格处理
		key := strings.Replace(param.Data.Key, " ", "", -1)
		//对待存储value做去反斜线处理
		value := strings.Replace(param.Data.Value, "\\", "", -1)
		//执行插入操作
		err = dbClient.Set(key, value)
		if err != nil {
			//出错的话提示
			response := util.SystemErr()
			response.Msg = "执行插入操作失败，失败原因：" + err.Error()
			return response
		}
		response := util.InsertSuccess()
		return response
	} else {
		//查询操作的处理
		results, err := dbClient.Do("scan", param.Data.Key+"!", param.Data.Key+"~", int64(queryNum))
		if err != nil {
			responseJn := util.SystemErr()
			responseJn.Msg = "查询失败，原因：" + err.Error()
			return responseJn
		} else if len(results) == 1 && results[0] == "ok" {
			//未查到数据
			responseJn := util.SelectSuccess()
			responseJn.Msg = "查询结束，查询到0条数据"
			return responseJn
		} else if len(results) > 0 && results[0] != "ok" {
			//不为ok时也是一种异常
			responseJn := util.SystemErr()
			responseJn.Msg += "：" + results[0]
			return responseJn
		}
		//存在有效查询结果 则组织返回
		var rows []pojo.Row
		for i := 1; i < len(results); i += 2 {
			var rowData pojo.Row
			rowData.Key = results[i]
			rowData.Value = results[i+1]
			rows = append(rows, rowData)
		}
		//返回结果
		responseJn := util.SelectSuccess()
		responseJn.Data = rows
		return responseJn
	}

}

//删除操作的rest服务
func deleteHandler(param pojo.RequestJson) util.ResponseJson {
	log.Println("in deleteHandler()")
	//前端将需要删除的多个key以逗号分隔开 在此需要用逗号分割取得key列表
	deleteRows := strings.Split(param.Data.Key, ",")
	//获取数据库连接
	dbClient, err := dao.GetConnection(param.Data.DbIP, param.Data.DbPortInt)
	if err != nil {
		//出错的话提示
		response := util.SystemErr()
		response.Msg = "获取数据库连接失败，失败原因：" + err.Error()
		return response
	}
	//进行删除操作
	_, err = dbClient.Do("multi_del", deleteRows)
	dbClient.Close() //关闭连接
	if err != nil {
		responseJn := util.SystemErr()
		responseJn.Msg = "删除失败:" + err.Error()
		return responseJn
	}
	//删除成功 返回提示
	responseJn := util.DeleteSuccess()
	responseJn.Msg = "删除成功，总共删除了" + fmt.Sprintf("%d", len(deleteRows)-1) + "条记录"
	return responseJn
}

//改变ssdb连接的服务
func changeConHandler(param pojo.RequestJson) util.ResponseJson {
	log.Println("in changeConHandler()")
	//建立连接池
	err := dao.NewPool(param.Data.DbIP, param.Data.DbPortInt)
	//异常处理
	if err != nil {
		responseJn := util.SystemErr()
		responseJn.Msg = "数据库连接失败，原因：" + err.Error()
		return responseJn
	}
	//更改连接成功 返回记录
	responseJn := util.SelectSuccess()
	responseJn.Msg = "连接成功"
	return responseJn
}
