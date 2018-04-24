package util

/** file: responseJson.go
*  为每个错误代码建立一个对应的获得响应responseJson的方法
*  通过包名.方法名即可获得需要的responseJson，省去输入常量名的麻烦
 */

import (
	"common"
)

// ResponseJson: 返回给前端的响应json数据
/*   Code 错误代码
*    Msg  错误提示信息
*    Data 返回的json字符串(如果有的话)
 */
type ResponseJson struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

//增加成功
func InsertSuccess() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = "300000"
	responseJsn.Msg = common.ErrorOKInsertMsg
	return responseJsn
}

//删除成功
func DeleteSuccess() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorOKId
	responseJsn.Msg = common.ErrorOKDeleteMsg
	return responseJsn
}

//修改成功
func ModifySuccess() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorOKId
	responseJsn.Msg = common.ErrorOKModifyMsg
	return responseJsn
}

//查询成功
func SelectSuccess() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorOKId
	responseJsn.Msg = common.ErrorOKGetMsg
	return responseJsn
}

//json格式错误
func JsonFmtErr() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorJsonFmtErrId
	responseJsn.Msg = common.ErrorJsonFmtErrMsg
	return responseJsn
}

//数据校验失败
func CheckFailed() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorCheckFailedId
	responseJsn.Msg = common.ErrorCheckFailedMsg
	return responseJsn
}

//系统错误
func SystemErr() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorSystemErrId
	responseJsn.Msg = common.ErrorSystemErrMsg
	return responseJsn
}

//数据已存在（新增记录时）
func DataExistsErr() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorDataExistsErrId
	responseJsn.Msg = common.ErrorDataExistsErrMsg
	return responseJsn
}

//所要删除的数据不存在
func DelDataNotExistsErr() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorDelDataNotExistsErrId
	responseJsn.Msg = common.ErrorDelDataNotExistsMsg
	return responseJsn
}

//修改数据不存在
func EditDataNotExistsErr() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorEditDataNotExistsErrId
	responseJsn.Msg = common.ErrorEditDataNotExistsMsg
	return responseJsn
}

//查询的数据不存在
func SelectDataNotExistsErr() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorGetDataNotExistsErrId
	responseJsn.Msg = common.ErrorGetDataNotExistsMsg
	return responseJsn
}

//非法请求 如command是PUT POST DELETE之外的字符
func RequstIllegal() ResponseJson {
	var responseJsn ResponseJson
	responseJsn.Code = common.ErrorRequstIllegalId
	responseJsn.Msg = common.ErrorRequstIllegalMsg
	return responseJsn
}
