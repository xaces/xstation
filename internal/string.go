package internal

import (
	jsoniter "github.com/json-iterator/go"
)

func ToJString(obj interface{}) string {
	data, err := jsoniter.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(data)
}
