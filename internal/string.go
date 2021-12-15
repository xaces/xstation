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

func StringIndex(s, sep string, n int) string {
	pos := 0
	for i := 0; i < n; i++ {
		lpos := strings.Index(s[pos:], sep)
		if lpos == -1 {
			break
		}
		pos += (lpos + 1)
	}
	return s[pos:]
}
