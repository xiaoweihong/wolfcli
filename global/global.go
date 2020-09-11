package global

import (
	"go.uber.org/zap"
	"net"
)

var (
	WOLF_LOG     *zap.Logger
	IP           net.IP
	DefaultOrgId string
)
