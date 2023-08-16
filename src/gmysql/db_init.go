package gmysql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lucky-cheerful-man/phoenix_server/src/config"
	"github.com/lucky-cheerful-man/phoenix_server/src/log"
)

var DBOperate *Mysql

// Setup initializes the database instance
func init() {
	DBOperate = new(Mysql)

	db, err := gorm.Open(config.ReferGlobalConfig().DatabaseSetting.Type, fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.ReferGlobalConfig().DatabaseSetting.User,
		config.ReferGlobalConfig().DatabaseSetting.Password,
		config.ReferGlobalConfig().DatabaseSetting.Host,
		config.ReferGlobalConfig().DatabaseSetting.Name))

	if err != nil {
		log.Error("mysql.Setup err: %s", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.ReferGlobalConfig().DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(config.ReferGlobalConfig().DatabaseSetting.MaxIdleConn)
	db.DB().SetMaxOpenConns(config.ReferGlobalConfig().DatabaseSetting.MaxOpenConn)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(config.ReferGlobalConfig().DatabaseSetting.ConnMaxLifeMinute))
	DBOperate.db = db
}

type DBInterface interface {
	InsertUser(name, password string) error
	CheckAuth(username, password string) (bool, string, string, error)
	GetProfile(name string) (string, string, error)
	EditProfile(name, path, nickname string) error
}
