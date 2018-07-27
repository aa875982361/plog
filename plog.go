package main

import (
	"fmt"
	"log"
	"net/http"

	"./controller"
	"./db"
	"./logger"

	"github.com/gorilla/mux"
)

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
