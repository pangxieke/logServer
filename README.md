# 日志服务
1. http服务，接受日志上报，存储（kafka，redis）
2. 常驻后台服务，定时读取数据（kafka，redis），将数据上报Elasticsearch

## config
配置项

消息队列，可以使用kafka或者redis作为消息队列。
```
# kafka/redis
storage = kafka
# kafka topic /redis key
topic = log_topic
# 是否同步存储
sync = true

```
开启定时服务，将[redis/kafka]中队列数据读取处理，发送到elastic
```$xslt
[time]
# 开启定时同步
sync = true
# 同步时间 s
sync_time=5
```
kafka, redis, elastic配置项
```$xslt
[kafka]
addrs = 127.0.0.1:9092

[redis]
addrs = 127.0.0.1:7001,127.0.0.1:7002

[elasticsearch]
hosts=http://127.0.0.1:9200
user=
secret=
```

## 启动
```$xslt
go run main/main.go
```

## 测试
```$xslt
python3 test/upload.py
```



