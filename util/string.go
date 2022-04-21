package util

import (
	"fmt"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func JString(v interface{}) string {
	str, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return str
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

func FilePath(s, deviceNo string) string {
	sarr := strings.Split(s, "/")
	i := len(sarr) - 1
	return fmt.Sprintf("%s/%s/%s/%s/%s", sarr[i-2],deviceNo, sarr[i-3], sarr[i-1], sarr[i])
}

func FilePicPath(s, deviceNo string) string {
	sarr := strings.Split(s, "/")
	i := len(sarr) - 1
	return fmt.Sprintf("%s/%s/%s/%s", sarr[i-1], deviceNo, sarr[i-2], sarr[i])
}

// ftp://admin:123456@127.0.0.1:2211
func FtpUriParse(s string) (port int, user, pswd string) {
	arrs := strings.Split(s, "@")
	s1 := strings.Split(arrs[0], ":")
	user = s1[1][2:]
	pswd = s1[2]
	s2 := strings.Split(arrs[1], ":")
	port, _ = strconv.Atoi(s2[1])
	return
}
