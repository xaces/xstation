module xstation

go 1.16

require (
	github.com/FishGoddess/cachego v0.2.5
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.2
	github.com/go-co-op/gocron v1.13.0 // indirect
	github.com/goftp/file-driver v0.0.0-20180502053751-5d604a0fc0c9
	github.com/goftp/server v0.0.0-20200708154336-f64f7c2d8a42
	github.com/jlaffaye/ftp v0.0.0-20211117213618-11820403398b // indirect
	github.com/json-iterator/go v1.1.11
	github.com/kardianos/service v1.2.0
	github.com/lesismal/arpc v1.2.7
	github.com/nats-io/nuid v1.0.1
	github.com/panjf2000/ants/v2 v2.4.6
	github.com/satori/go.uuid v1.2.0
	github.com/unrolled/secure v1.0.9
	github.com/wlgd/xproto v0.0.0-00010101000000-000000000000
	github.com/wlgd/xutils v0.0.0-20220511104346-0bbafde87f52
)

replace github.com/wlgd/xproto => ../xproto
