package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"../constant"
	"../db"
	"../logger"
	"../model"
	"../utils/feedback"
)

func InitHeaders(w http.ResponseWriter) {
	if w == nil {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	return
}

func HandleWeblog(w http.ResponseWriter, r *http.Request) {

	if w == nil || r == nil {
		fmt.Println("resposeWriter or request is nil")
		return
	}

	InitHeaders(w)

	fb := feedback.NewFeedBack(w)

	var wl model.WebLog

	if r.Method == http.MethodPost {
		result, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Warning(fmt.Errorf("read request body failed: %v", err))
			fb.Code(constant.GLOBAL_SYS_ERR).Msg(constant.GLOBAL_SYS_ERR_MSG).Response()
			return
		}
		r.Body.Close()

		err = json.Unmarshal([]byte(result), &wl)
		if err != nil {
			logger.Warning(fmt.Errorf("unmarshal web log failed: %v", err))
			fb.Code(constant.GLOBAL_PARM_ERR).Msg(constant.GLOBAL_PARM_ERR_MSG).Response()
			return
		}
	}

	if r.Method == http.MethodGet {
		rqs, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			logger.Warning(fmt.Errorf("parse url failed: %v", err))
			fb.Code(constant.GLOBAL_PARM_ERR).Msg(constant.GLOBAL_PARM_ERR_MSG).Response()
			return
		}
		if rqs["user"] != nil {
			wl.User = rqs["user"][0]
		}
		if rqs["type"] != nil {
			wl.Type = rqs["type"][0]
		}
		if rqs["project"] != nil {
			wl.Project = rqs["project"][0]
		}
		if rqs["tag"] != nil {
			wl.Tag = rqs["tag"][0]
		}
		if rqs["detail"] != nil {
			wl.Detail = rqs["detail"][0]
		}
		if rqs["createTime"] != nil {
			wl.CreateTime, err = strconv.ParseInt(rqs["createTime"][0], 10, 64)
		}

		if err != nil {
			logger.Warning(fmt.Errorf("parse createTime failed: %v", err))
			fb.Code(constant.GLOBAL_PARM_ERR).Msg(constant.GLOBAL_PARM_ERR_MSG).Response()
			return
		}
	}

	bufs, err := json.Marshal(wl)
	if err != nil {
		logger.Error(err)
	}

	fmt.Println(string(bufs))

	fb.Code(constant.GLOBAL_SUCCESS).Response()

	db.WritesPoints(wl)
}
