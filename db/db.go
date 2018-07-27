package db

import (
	"log"
	"time"

	"../model"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	MyDB     = "plog"
	username = "*****"
	password = "****"
	addr     = "*****"
)

// 全局执行插入操作的对象
var cli client.Client

func Init() error {

	conn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})

	if err != nil {
		return err
	}

	cli = conn
	return nil
}

//Insert
func WritesPoints(weblog model.WebLog) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{
		"level":   weblog.Type,
		"user":    weblog.User,
		"project": weblog.Project,
	}
	fields := map[string]interface{}{
		"tag":    weblog.Tag,
		"detail": weblog.Detail,
	}

	pt, err := client.NewPoint(
		"weblog",
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	if err := cli.Write(bp); err != nil {
		log.Fatal(err)
	}
}
