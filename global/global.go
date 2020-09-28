package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net"
	"time"
)

var (
	WOLF_LOG     *zap.Logger
	IP           net.IP
	DBIP         net.IP
	PicIP        net.IP
	PicPort      string
	Day          string
	DefaultOrgId string
	PgUsername   string
	PgPassword   string
	DbName       string
	Db           *gorm.DB
	StartTime    int64
	EndTime      int64
	ThresHold    float64
	HttpTimeOut  time.Duration
	ParamNum     int
)
