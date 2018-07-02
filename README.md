# 使用方法
## 1.下载第三方包
>由于使用了seefan的ssdb客户端，因此需要使用以下命令将相关代码下载到GOPATH。
```
go get github.com/seefan/gossdb
```
## 2.编译程序
>进入ssdbVisualQueryTool目录，设置该路径为GOPATH，执行编译，参考编译命令如下：
```
export GOPATH=`pwd`
go build -o bin/ssdbTool main
```
## 3.修改配置文件
>配置文件为ssdbVisualQueryTool/etc/conf.json文件，其中service部分为访问端口号的配置；dbDefault为数据库连接池的默认配置；dbCustom为自定义数据库连接池配置，可用于对常用数据库连接的优化。

## 4.执行程序
>bin目录下的start.sh中提供了以下启动命令作为参考。由于读取静态网页文件用的是相对路径，以根目录ssdbVisualQueryTool为基准，因此请在根目录通过bin/ssdbTool的方式启动程序，否则可能会在加载静态页面的时候抛出异常。
```
nohup bin/ssdbTool >logs/`date +%Y%m%d%H%M%S`.log 2>&1 &
```
## 5.访问ssdb可视化工具
>打开浏览器输入localhost:配置的端口号即可访问ssdb可视化工具，本例为localhost:8080，打开之后界面如下。前端技术有限，界面略丑~

![ssdb可视化工具界面](https://github.com/dcbigcrafter/screenShorts/raw/master/ssdbVisualQueryTool/1.png)
# 支持的操作
## 1.连接数据库
>输入ssdb的IP地址以及端口号进行连接，连接成功则在右侧以绿色文字提示数据库已连接，失败则以红色字体提示连接失败以及失败原因。

![连接数据库的操作](https://github.com/dcbigcrafter/screenShorts/raw/master/ssdbVisualQueryTool/2.png)
## 2.查询数据
>查询以输入字符为key或者以输入字符开头的key对应的value。

![查询数据库](https://github.com/dcbigcrafter/screenShorts/raw/master/ssdbVisualQueryTool/4.png)
## 3.增加数据
>输入key和value以执行新增或修改的操作。

![增加数据](https://github.com/dcbigcrafter/screenShorts/raw/master/ssdbVisualQueryTool/3.png)
## 4.删除数据
>勾选查询出来的数据，进行删除。

![增加数据](https://github.com/dcbigcrafter/screenShorts/raw/master/ssdbVisualQueryTool/5.png)
# 开发缘由
ssdb自带的命令行连接工具存在以下缺陷：查询结果堆砌在一起不够清晰明了；删除多个key比较麻烦。平常的开发中又需要频繁的查询或删除ssdb中的测试数据，使用命令行进行查询删除操作非常不方便。索性就利用刚学到的goweb知识开发了一套ssdb可视化管理工具，对于ssdb的连接以及增删改查操作都能够通过页面进行操作，比命令行工具方便很多。本来只是将服务挂在了公司内部服务器分配给个人的端口上自己用，后来发现公司很多同事甚至产品经理也在用，已然成了大家的标配开发工具，感觉还是蛮有成就感的。
