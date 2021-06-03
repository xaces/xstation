package service

import "xstation/internal"

var (
	_snowflake *internal.Snowflake = nil
)

func PrimaryKey() int64 {
	return _snowflake.NextId()
}

// Init 初始化服务
func Init() error {
	if err := sqlInit(); err != nil {
		return err
	}
	_snowflake, _ = internal.NewSnowflake(0xFF)
	return nil
}
