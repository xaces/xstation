module xstation

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/json-iterator/go v1.1.11
	github.com/nats-io/nuid v1.0.1
	github.com/panjf2000/ants/v2 v2.4.6
	github.com/satori/go.uuid v1.2.0
	github.com/smallnest/rpcx v0.0.0-20210120041900-c2830baacdb1
	github.com/wlgd/xproto v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b // indirect
	gorm.io/driver/mysql v1.0.5
	gorm.io/driver/postgres v1.0.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.3
)

replace github.com/wlgd/xproto => ../xproto
