package util

/** file: util.go
*  放置工具方法
*  与具体的工程解耦 与业务无关
 */

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

//获取body的data(json)转换为string
func GetRqtBodyStr(r *http.Request) (string, error) {
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

//获取纯数字格式的时间 格式:20170905164628
func GetDateDigit() (date string) {
	//获取时间戳
	tamp := time.Now().Unix()
	//格式化字符串按照年月日输出
	//go语言中格式化字符串特殊的含义月01,日02,时03,分04,秒05,年2006
	date = time.Unix(tamp, 0).Format("20060102150405")
	return date
}

//获取经过格式话的时间 格式：YYYY-MM-DD hh:mm:ss
func GetFormattedTime() (date string) {
	//获取时间戳
	tamp := time.Now().Unix()
	//格式化字符串按照年月日输出
	//go语言中格式化字符串特殊的含义月01,日02,时03,分04,秒05,年2006
	date = time.Unix(tamp, 0).Format("2006-01-02 15:04:05")
	return date
}

// 设置报文头信息 解决跨域问题
func SetHeader(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

//获取GET方式传递过来的参数
func GetGETParams(r *http.Request, v ...string) map[string]string {
	r.ParseForm()
	params := make(map[string]string)
	for _, param := range v {
		if len(r.Form[param]) == 0 {
			(params)[param] = "" //未传值则设为空字符串
		} else {
			(params)[param] = r.Form[param][0] //len(r.Form[param])不为0则获取其值
		}
	}
	return params
}

//结构体属性合法性校验的通用方法
func CheckParam(param string, paramCHName string, length int, fixedLength bool, nullable bool) (responseStr string, checkResult bool) {
	checkResult = true //将checkResult初始化为true
	//不可为空且参数值为空或空白字符串则校验失败
	if nullable == false && (len(param) == 0 || param == "") {
		responseStr = paramCHName + "不能为空"
		checkResult = false
	}
	//定长参数 校验长度是否可法
	if fixedLength && len(param) != 0 && len(param) != length {
		//如GroupId定长8个字符，不为空时，检查长度是否合法
		responseStr = paramCHName + "长度不合法"
		checkResult = false
	}
	//不定长参数 length不为0时 校验长度是否超过规定值
	if fixedLength == false && len(param) != 0 && len(param) > length && length != 0 {
		responseStr = paramCHName + "长度超出限制"
		checkResult = false
	}
	return responseStr, checkResult
}

//本方法用于将字符串的首字母变成大写
func StrFirstToUpper(str string) string {
	var upperStr string
	for i := 0; i < len(str); i++ {
		if i == 0 {
			upperStr += strings.ToUpper(string(str[i]))
		} else {
			upperStr += string(str[i])
		}
	}
	return upperStr
}

//本方法用于在golang中请求rest服务
func RequestRestService(requestUrl, requestData string) (bodyStr string, err error) {
	var jsonStr = []byte(requestData)
	req, err1 := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err1 != nil {
		fmt.Println("error occcurred in requestRestService(),content:", err1)
		return "", err1
	}

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println("error occcurred in requestRestService(),content:", err2)
		return "", err2
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr = string(body)
	return bodyStr, err
}

//从ssdb存储的key中提取出表名 提取从第一个小写字母到第一个非小写字母中间的字符串
//如果第一个字符就不是小写字母 那么返回该key
func GetTableName(str string) string {
	//将字符串转化为byte数组
	strByte := []byte(str)
	//定义一个byte数组用于接收表名
	var tableByte []byte
	//遍历字符串 判断是否为小写
	for _, v := range strByte {
		if v >= 'a' && v <= 'z' {
			tableByte = append(tableByte, v)
		} else {
			//不为小写则退出
			break
		}
	}
	//如果从一开始就不是小写字母 那么直接返回key
	if len(tableByte) == 0 {
		return str
	}
	return string(tableByte)
}

//获取页码以及每页显示数量
func GetPageNumSize(pageNumString, pageSizeString string) (pageNum, pageSize int, err error) {
	var warnMsg string //保存转换为数字失败时的错误提示
	//分页查询的页码转为数字
	pageNum, err2 := strconv.Atoi(pageNumString)
	if err2 != nil {
		warnMsg += "页码：" + pageNumString + ",转换成整数失败。"
	}
	//分页查询的每页显示条数转为数字
	pageSize, err1 := strconv.Atoi(pageSizeString)
	if err1 != nil {
		warnMsg += "每页显示条数：" + pageSizeString + ",转换成整数失败。"
	}
	//判断pageSize是否小于等于0 因为后面需要做被除数
	if pageSize == 0 {
		warnMsg += "每页显示条数不可以为0。"
	}
	//判断pageSize是否小于0 小于0则改为20
	if pageSize < 0 {
		pageSize = 20
	}
	//判断pageSize是否小于等于0 是的话改为1
	if pageNum <= 0 {
		pageNum = 1
	}
	//如果错误提示不为空 需要返回错误信息
	if warnMsg != "" {
		return 0, 0, errors.New(warnMsg)
	}
	return pageNum, pageSize, nil
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成序列号
func GetSerialNumber() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		//发生错误时用时间戳作为序列号
		return fmt.Sprintf("%d", time.Now().Unix())
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

//去除字符串中的空格 tab 换行 等一系列字符串处理
func TrimString(str string) string {
	//去字符串中的空格
	str = strings.Replace(str, " ", "", -1)
	//去字符串中的tab空格
	str = strings.Replace(str, "\n", "", -1)
	//去字符串中的换行符
	str = strings.Replace(str, "\t", "", -1)
	return str
}

//将string数组里的元素正序排列 执行完成后原数组即为有序 本排序区分大小写
func SortStringArray(stringArray *[]string) {
	//强转类型将[]string转为StringSlice
	stringSlice := sort.StringSlice(*stringArray)
	//排序
	sort.Sort(stringSlice)
}
