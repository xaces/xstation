package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)
// LoadJSONConf 初始化配置参数
func LoadJSONConf(jsonFile string, obj interface{}) {
	jsonFp, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("load error" + jsonFile)
		os.Exit(0)
	}
	defer jsonFp.Close()
	var jsString string
	iReader := bufio.NewReader(jsonFp)
	for {
		tString, err := iReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		jsString = jsString + tString
	}
	if err := json.Unmarshal([]byte(jsString), obj); err != nil {
		fmt.Println("json error " + jsonFile)
		os.Exit(0)
	}
}