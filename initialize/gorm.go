package initialize

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"wolfcli/global"
)

// GormPostgreSql 初始化PostgreSql数据库
func GormPostgreSql() {
	dsn := "host=" + global.DBIP.String() + " user=" + global.PgUsername + " password=" + global.PgPassword + " dbname=" + global.DbName + " port=5432" + " " + "sslmode=disable TimeZone=Asia/Shanghai"
	postgresConfig := postgres.Config{
		DSN:                  dsn,  // DSN data source name
		PreferSimpleProtocol: true, // 禁用隐式 prepared statement
	}
	gormConfig := config()
	if global.Db, err = gorm.Open(postgres.New(postgresConfig), gormConfig); err != nil {
		zap.L().Error("PostgreSql连接异常", zap.Any("err", err))
		os.Exit(0)
	} else {
		sqlDB, _ := global.Db.DB()
		sqlDB.SetMaxIdleConns(100)
		sqlDB.SetMaxOpenConns(100)
	}
}

func config() (c *gorm.Config) {

	c = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	return
}
