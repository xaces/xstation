package rpc

import (
	"xstation/entity/cache"

	jsoniter "github.com/json-iterator/go"
)

type IRpc interface {
	Subscribe(string, func([]byte)) error
	Publish(string, interface{}) error
	Release()
}

var gRpc IRpc

func Handler(b []byte) {
	j := jsoniter.Get(b)
	if j.LastError() != nil {
		return
	}
	// fmt.Printf("%s\n", b)
	code := j.Get("code").ToInt()
	switch code {
	case 5000: // 同步服务器信息

	case 5010, 5011: // 设备管理
		var vehis []cache.DeviceInfo
		j.Get("data").ToVal(&vehis)
		for _, v := range vehis {
			if code == 5010 {
				cache.NewDevice(v)
			} else {
				cache.DelDevice(v.No)
			}
		}
	case 5012: // 比如指定Ftp上传
	}
}

// 这里Topic用服务guid
func Run(s IRpc, topic string) {
	if s == nil {
		return
	}
	gRpc = s
	gRpc.Subscribe(topic, Handler)
}

func Publish(topic string, v interface{}) error {
	if gRpc == nil {
		return nil
	}
	return gRpc.Publish(topic, v)
}

func Shutdown() {
	if gRpc == nil {
		return
	}
	gRpc.Release()
}
