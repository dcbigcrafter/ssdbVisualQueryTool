package common

const (
	//服务调用没有发生任何错误
	ErrorOKId           = "200000" //所有的服务都以这个编码为正确无误调用之后的返回值
	ErrorOKMsg          = "OK"     //所有的服务都以这个消息为200000错误码对应的消息
	ErrorOKInsertMsg    = "增加成功"
	ErrorOKDeleteMsg    = "删除成功"
	ErrorOKModifyMsg    = "修改成功"
	ErrorOKGetMsg       = "数据查询成功"
	ErrorOKSetPwdMsg    = "密码重置成功"
	ErrorOKSetStateMsg  = "状态修改成功"
	ErrorOKOperStateMsg = "注销成功"
	//机构ID为空
	ErrorGroupIdIsNullId  = "200001"
	ErrorGroupIdIsNullMsg = "机构标识不能为空"

	//输入数据为空
	ErrorDataIsNullId  = "200002"
	ErrorDataIsNullMsg = "数据为空"

	//操作者Id为空
	ErrorOperatorIdIsNullId  = "200003"
	ErrorOperatorIdIsNullMsg = "操作者Id不能为空"

	//json格式错误
	ErrorJsonFmtErrId  = "200004"
	ErrorJsonFmtErrMsg = "json格式错误"

	//数据校验失败
	ErrorCheckFailedId  = "200005"
	ErrorCheckFailedMsg = "提交的数据未能通过校验"

	//系统错误
	ErrorSystemErrId  = "200006"
	ErrorSystemErrMsg = "系统错误"

	//数据已存在（新增记录时）
	ErrorDataExistsErrId  = "200007"
	ErrorDataExistsErrMsg = "数据已存在"

	//查询的数据不存在
	ErrorGetDataNotExistsErrId = "200008"
	ErrorGetDataNotExistsMsg   = "查询数据不存在"

	//删除数据不存在
	ErrorDelDataNotExistsErrId = "200009"
	ErrorDelDataNotExistsMsg   = "删除数据不存在"

	//修改数据不存在
	ErrorEditDataNotExistsErrId = "200010"
	ErrorEditDataNotExistsMsg   = "修改数据不存在"

	//系统查询错误
	ErrorSystemSelectErrId  = "200011"
	ErrorSystemSelectErrMsg = "系统查询错误"

	//非法请求 如command是PUT POST DELETE之外的字符
	ErrorRequstIllegalId  = "200012"
	ErrorRequstIllegalMsg = "无效的请求参数"

	//数据库打开失败
	ErrorOpenSqliteId  = "200013"
	ErrorOpenSqliteMsg = "数据库打开失败"

	//新增数据失败
	ErrorInsertId  = "200014"
	ErrorInsertMsg = "新增数据失败"

	//修改数据失败
	ErrorEditId  = "200015"
	ErrorEditMsg = "修改数据失败"

	//删除数据失败
	ErrorDeleteId  = "200016"
	ErrorDeleteMsg = "删除数据失败"
)
