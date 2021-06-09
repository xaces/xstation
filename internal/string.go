package internal

import (
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func ToJString(obj interface{}) string {
	data, err := jsoniter.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(data)
}

func StringToIntSlice(str, sep string) []int {
	strv := strings.Split(str, sep)
	var intv []int
	for _, v := range strv {
		if v == "" {
			continue
		}
		val, _ := strconv.Atoi(v)
		intv = append(intv, val)
	}
	return intv
}
