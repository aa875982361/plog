# plog

## 简介 

>之前一直想搭建一个属于自己 的前端日志系统，参考了挺多资料，现在来正式 搭建一个吧

>技术选型原因 后端语言 相对来说 golang比较熟悉，关于日志系统，在面试头条的时候有了解到时间序列数据库，上网找了挺多资料，发现influx 是用 go 实现的，所以想尝试一下，当然 这个也不是说最好的，对于我来说 还是挺合适的

技术栈： golang + influx  + grafana 

## influx 安装

influx 的安装 可以看看这位大神的[博客](https://www.jianshu.com/p/d2935e99006e)

###  账密设置

在这个文件里面设置：
`./db/db.go`

## goalng 代码编写

先装几个依赖 

``` cmd
go get github.com/gorilla/mux
go get go.uber.org/zap/zapcore
go get go.uber.org/zap/
go get github.com/influxdata/influxdb/client/v2
```

### 主要代码

``` javascript

func main() {
	// 初始化数据库操作结构体
	err := db.Init()

	if err != nil {
		log.Panic("Init database err:" + err.Error())
	}
	// 初始化后端日志系统
	err = logger.Init(".")
	if err != nil {
		log.Panic("Init logger err:" + err.Error())
	}
	defer logger.Sync()

	fmt.Println("Service listen at port 6609")
	// 初始化监听事件
	err = http.ListenAndServe(":6609", setRouter())

	if err != nil {
		logger.Error(err)
	}
}

// 设置路由
func setRouter() *mux.Router {
	r := mux.NewRouter()
	// 监听 /api/weblog
	r.HandleFunc("/api/weblog", controller.HandleWeblog)

	return r
}
```

### 启用

``` cmd
go build plog.go
nohup ./plog &
```

api示例：

![api示例](https://www.github.com/aa875982361/plingWeb/raw/master/markdown/1532676689511.png)

POST 127.0.0.1:6609/api/weblog

``` json
{
  "type": "info", // 等级 
  "project": "test", // 项目名称
  "user": "plinghuang", // 使用人
  "tag": "dddd", // 简要概述
  "detail": " ddsadasd", // 具体情况
  "createtime": 1532490304939 // 时间戳
}
```

### 安全性

之前就有想过,担心别人滥用我的接口去操作数据,暂时想到有几层防御机制
1. 前端数据加密传输
2. 根据时间戳 去判定 是否小于当前时间太多
3. 设置使用量 给固定的人 在某段时间可以使用的内存大小 如果快超出则预警 检查

## 前端使用 

1. 采用主动埋点 的方式 主动触发日志系统
2. 采用 ajax 发送请求 
3. 采用 队列结构 在发送的时候做控制 防止发送太频繁

参考链接：[美团点评前端无痕埋点实践](https://tech.meituan.com/mt-mobile-analytics-practice.html)
[log + ajax，前端日志解决方案](https://github.com/eshengsky/lajax)
[可视化工具配置采集节点，在前端自动解析配置并上报埋点数据](https://github.com/mixpanel/mixpanel-js)
[构建web前端异常监控系统](http://www.aliued.cn/2012/10/27/%E6%9E%84%E5%BB%BAweb%E5%89%8D%E7%AB%AF%E5%BC%82%E5%B8%B8%E7%9B%91%E6%8E%A7%E7%B3%BB%E7%BB%9F-fdsafe.html)