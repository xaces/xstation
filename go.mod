module xstation

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.2
	github.com/json-iterator/go v1.1.11
	github.com/nats-io/nats-server/v2 v2.3.2 // indirect
	github.com/nats-io/nuid v1.0.1
	github.com/panjf2000/ants/v2 v2.4.6
	github.com/satori/go.uuid v1.2.0
	github.com/unrolled/secure v1.0.9
	github.com/urfave/cli/v2 v2.3.0
	github.com/wlgd/xproto v0.0.0-00010101000000-000000000000
	github.com/wlgd/xutils v0.0.0-20210701074559-e4b0685b2ff6
	gorm.io/driver/mysql v1.0.5
	gorm.io/driver/postgres v1.0.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)

replace github.com/wlgd/xproto => ../xproto

replace github.com/wlgd/xutils => ../xutils
