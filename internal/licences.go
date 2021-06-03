package internal

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

const (
	licencesFile   = "licences"
	identification = "nlasfl2wrwfsnsfs131#$%fs"
)

type LicensingInf struct {
	Licences   string `json:"licences"`   // 授权信息
	MQAddress  string `json:"mqAddress"`  // MQ订阅信息
	RpcAddress string `json:"rpcAddress"` // rpc地址
}

// WriteLicences 写入文件
func WriteLicences(lice *LicensingInf) error {
	data, err := json.Marshal(lice)
	if err != nil {
		return err
	}
	os.Remove(licencesFile)
	fp, err := os.Create(licencesFile)
	if err == nil {
		fp.Write([]byte(identification))
		fp.Write(data)
	}
	defer fp.Close()
	return nil
}

func ReadLicences() (lice *LicensingInf) {
	fp, err := os.Open(licencesFile)
	if err != nil {
		return
	}
	defer fp.Close()
	var data []byte
	buf := make([]byte, 1024)
	bfRd := bufio.NewReader(fp)
	for {
		n, err := bfRd.Read(buf)
		data = append(data, buf[:n]...)
		if err != nil {
			break
		}
	}
	if !strings.HasPrefix(string(data), identification) {
		data = data[:0]
	}
	json.Unmarshal(data[len(identification):], lice)
	return
}
