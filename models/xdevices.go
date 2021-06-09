package models

type XDeviceOpt struct {
	VehiNo   string `json:"vehiNo"`
	VehiName string `json:"vehiName"`
	Icons    string `json:"icon"`
}

type XDevice struct {
	Model
	XDeviceOpt
}

// TableName 表名
func (s *XDevice) TableName() string {
	return "t_xdevice"
}
