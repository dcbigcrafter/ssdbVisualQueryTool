#运行时生成以当前时间为文件名的日志文件
nohup bin/ssdbTool >logs/`date +%Y%m%d%H%M%S`.log 2>&1 &



