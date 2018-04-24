package main

import (
	"dao"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"pojo"
	"util"
)

//请求分发 根据请求路径转至不同路由
func restHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("-------------in dispatchHandler--------------")
	//针对POST提交方式的处理
	if r.Method == "POST" {
		var requestJs pojo.RequestJson
		//通过对request的解析 获得前端传来的参数
		bodyStr, err := util.GetRqtBodyStr(r)
		if err != nil {
			response := util.SystemErr()
			response.Msg += "，读取请求内容异常，原因为：" + err.Error()
			responseByte, _ := json.Marshal(response)
			w.Write(responseByte)
			return
		}
		//获取访问者ip
		clientIp := r.RemoteAddr
		//记录本次请求的相关信息
		log.Printf("IP为%v的访问者进行了如下请求：%v\n", clientIp, bodyStr)
		//反序列化请求参数
		err = json.Unmarshal([]byte(bodyStr), &requestJs)
		//反序列化失败则返回错误提示
		if err != nil {
			response := util.JsonFmtErr()
			response.Msg += "发生在反序列化请求内容，原因为：" + err.Error()
			responseByte, _ := json.Marshal(response)
			w.Write(responseByte)
			return
		}
		//数据合法性校验
		warnMsg := CheckRqtData(&requestJs)
		//校验不通过 返回提示信息
		if warnMsg != "" {
			responseJsn := util.CheckFailed()
			responseJsn.Msg += "：" + warnMsg
			responseJsnByte, _ := json.Marshal(responseJsn)
			w.Write(responseJsnByte)
			return
		}
		//校验通过 针对不同com做不同的处理 将数据传至service层
		switch requestJs.Com {
		case "QUERY":
			serviceHandler(w, r, requestJs, queryHandler)
		case "DELETE":
			serviceHandler(w, r, requestJs, deleteHandler)
		case "CONNECT":
			serviceHandler(w, r, requestJs, changeConHandler)
		default:
			responseJsn := util.RequstIllegal()
			responseJsn.Msg += "：com=" + requestJs.Com
			responseJsnByte, _ := json.Marshal(responseJsn)
			w.Write(responseJsnByte)
			return
		}
	} else {
		responseJsn := util.RequstIllegal()
		responseJsn.Msg = "不支持的提交方式：" + r.Method + "！"
		responseJsnByte, _ := json.Marshal(responseJsn)
		w.Write(responseJsnByte)
		return
	}
}

//打开查询页
func defaultPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in defaultPageHandler()")
	response, _ := template.ParseFiles("./webRoot/index.html")
	response.Execute(w, nil)
}

//查看连接池状态的服务
func poolInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in poolInfoHandler()")
	msg := dao.StatusOfPools()
	w.Write([]byte(msg))
}

//调用具体的rest服务处理函数 获取处理结果返回给前端
func serviceHandler(w http.ResponseWriter, r *http.Request, param pojo.RequestJson, f func(pojo.RequestJson) util.ResponseJson) {
	respJsn := f(param)
	respByte, _ := json.Marshal(respJsn)
	w.Write(respByte)
}
