package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type baiduLocation struct {
	Status int `json:"status"`
	Result struct {
		Address string `json:"formatted_address"`
	} `json:"result"`
}

type MapOption struct {
	Name string
	Key  string
}

var opt MapOption

// HttpGet http get 请求
func httpGet(url string, result interface{}) error {
	recv, err := http.Get(url)
	if err != nil {
		return err
	}
	defer recv.Body.Close()
	content, err := ioutil.ReadAll(recv.Body)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(content, &result)
}

func GetLocation(longtitude, latitude float32) string {
	switch opt.Name {
	case "Baidu":
		urlstr := fmt.Sprintf("http://api.map.baidu.com/reverse_geocoding/v3/?ak=%s&output=json&coordtype=wgs84ll&location=%f,%f", opt.Key, latitude, longtitude)
		var lo baiduLocation
		if err := httpGet(urlstr, &lo); err == nil && lo.Status == 0 {
			return lo.Result.Address
		}
	}
	return ""
}
