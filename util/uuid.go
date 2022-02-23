package util

import (
	"github.com/nats-io/nuid"
	uuid "github.com/satori/go.uuid"
)

//UUID 生成Guid字串
func UUID() string {
	u := uuid.NewV4()
	return u.String()
}

// ServeId 生成服务ID
func ServeId() string {
	return nuid.Next()
}
