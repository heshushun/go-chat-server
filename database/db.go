package database

import (
	"fmt"
	"go-chat/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var Db *gorm.DB

// InitDatabase 初始化数据库的连接
func InitDatabase() {
	var err error
	Db, err = gorm.Open(config.DatabaseSetting.DbType,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.DatabaseSetting.DbUser,
			config.DatabaseSetting.DbPassword,
			config.DatabaseSetting.DbHost,
			config.DatabaseSetting.DbPort,
			config.DatabaseSetting.DbName,
		))

	if err != nil {
		log.Fatal("mysql open error", err)
	}

	Db.SingularTable(true)

	Db.DB().SetMaxIdleConns(10)

	Db.DB().SetMaxOpenConns(100)

	Db.DB().SetConnMaxLifetime(5 * time.Second)

}
