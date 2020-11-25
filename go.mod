module wolfcli

go 1.14

require (
	github.com/astaxie/beego v1.12.3
	github.com/chrislusf/seaweedfs v0.0.0-20201116054817-c0d279c54e56
	github.com/dustin/go-humanize v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.1.1
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/postgres v1.0.5
	gorm.io/gorm v1.20.6
)

replace go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200425165423-262c93980547
