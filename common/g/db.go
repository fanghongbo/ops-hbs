package g

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// 初始化数据库连接
func InitDB() {
	var (
		dbConn string
		err    error
	)

	dbConn = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Db,
	)

	db, err = gorm.Open("mysql", dbConn)
	if err != nil {
		dlog.Fatalf("failed to connect mysql server: %s", err.Error())
		return
	}

	if err = db.DB().Ping(); err != nil {
		dlog.Fatalf("ping db err: %s", err.Error())
		return
	}

	if db != nil {
		db.LogMode(config.Debug)
		db.DB().SetMaxIdleConns(config.Database.MaxIdle)
		db.DB().SetMaxOpenConns(config.Database.MaxConn)
	}
}

// 获取db实例
func DB() *gorm.DB {
	return db
}
