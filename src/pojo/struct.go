package pojo

//解析前端发来的json数据
type RequestJson struct {
	Com  string `json:"com"`  //请求标识 区分不同的请求
	Data params `json:"data"` //请求参数
}

//接收前端传来的参数
type params struct {
	Key       string `json:"key"`    //新增或者查询的key 删除是以逗号分开的key列表 连接数据库时是空
	Value     string `json:"value"`  //新增的value 查询、删除和连接数据库时是空
	DbIP      string `json:"dbIP"`   //连接数据库时传数据库ip
	DbPort    string `json:"dbPort"` //连接数据库时传数据库端口
	Number    string `json:"number"` //查询时传查询数量
	DbPortInt int    //存储转换成int的端口号
}

//包装key value 键值对
type Row struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
