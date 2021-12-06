module xstation

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.2
	github.com/goftp/file-driver v0.0.0-20180502053751-5d604a0fc0c9
	github.com/goftp/server v0.0.0-20200708154336-f64f7c2d8a42
	github.com/jlaffaye/ftp v0.0.0-20211117213618-11820403398b // indirect
	github.com/json-iterator/go v1.1.11
	github.com/kardianos/service v1.2.0
	github.com/nats-io/nuid v1.0.1
	github.com/panjf2000/ants/v2 v2.4.6
	github.com/satori/go.uuid v1.2.0
	github.com/unrolled/secure v1.0.9
	github.com/wlgd/xproto v0.0.0-00010101000000-000000000000
	github.com/wlgd/xutils v0.0.0-20210701074559-e4b0685b2ff6
	gorm.io/driver/mysql v1.0.5
	gorm.io/driver/postgres v1.0.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)

replace github.com/wlgd/xproto => ../xproto

replace github.com/wlgd/xutils => ../xutils
