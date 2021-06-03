package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

const (
	contentType = "application/json;charset=utf8"
)

type HttpResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// httpPost http post 请求
func HttpPost(url string, requset interface{}, result interface{}) error {
	bs, err := json.Marshal(requset)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bs)
	recv, err := http.Post(url, contentType, body)
	if err != nil {
		return err
	}
	defer recv.Body.Close()
	content, err := ioutil.ReadAll(recv.Body)
	if err != nil {
		return err
	}
	var res HttpResult
	if err := jsoniter.Unmarshal(content, &res); err != nil {
		return err
	}
	if res.Code != 200 {
		return errors.New(res.Msg)
	}
	jsoniter.Get(content, "data").ToVal(result)
	return nil
}

// httpPost http post 请求
func HttpGet(url string, result interface{}) error {
	recv, err := http.Get(url)
	if err != nil {
		return err
	}
	defer recv.Body.Close()
	content, err := ioutil.ReadAll(recv.Body)
	if err != nil {
		return err
	}
	var res HttpResult
	if err := jsoniter.Unmarshal(content, &res); err != nil {
		return err
	}
	if res.Code != 200 {
		return errors.New(res.Msg)
	}
	jsoniter.Get(content, "data").ToVal(result)
	return nil
}
