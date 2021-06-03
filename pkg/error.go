package pkg

import (
	"errors"
)

var (
	ErrServeIsRunning = errors.New("serve is running")    // ErrServeIsRunning 服务已运行
	ErrHttpURL        = errors.New("http url error")      // ErrHttpURL 请求URL错误
	ErrPostRequest    = errors.New("request error")       // ErrPostRequest post请求错误
	ErrParameter      = errors.New("parameter error")     // ErrParameter 参数错误
	ErrStationLogin   = errors.New("station login error") // ErrStationLogin 登录错误
	ErrInner          = errors.New("server error")        // ErrInner 内部错误(更新)
)
