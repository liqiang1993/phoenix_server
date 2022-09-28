package gmysql

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/lucky-cheerful-man/phoenix_server/src/config"
	"github.com/lucky-cheerful-man/phoenix_server/src/log"
	"time"
)

// Setup initializes the database instance
func Setup() *Mysql {
	var MysqlOperate Mysql

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
	MysqlOperate.db = db

	return &MysqlOperate
}

type Mysql struct {
	db *gorm.DB
}

// InsertUser 注册接口
func (m *Mysql) InsertUser(name, password string) error {
	user := UserInfoTab{
		UserName:    name,
		Password:    password,
		NickName:    name,
		GmtCreated:  time.Now(),
		GmtModified: time.Now(),
	}
	if err := m.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// CheckAuth 登陆认证接口
func (m *Mysql) CheckAuth(username, password string) (bool, string, string, error) {
	var auth UserInfoTab
	err := m.db.Select([]string{"id", "nick_name", "image"}).Where(UserInfoTab{UserName: username,
		Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, "", "", err
	}

	if len(auth.NickName) > 0 {
		return true, auth.NickName, auth.Image, nil
	}

	return false, "", "", nil
}

// GetProfile 查询用户的属性信息
func (m *Mysql) GetProfile(name string) (string, string, error) {
	var auth UserInfoTab
	err := m.db.Select([]string{"nick_name", "image"}).Where(UserInfoTab{UserName: name}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", "", err
	}

	if len(auth.NickName) > 0 {
		return auth.NickName, auth.Image, nil
	} else {
		log.Warn("can not find valid profile info by name:%s", name)
		return "", "", errors.New("not found")
	}
}

// EditProfile 编辑用户属性
func (m *Mysql) EditProfile(name, path, nickname string) error {
	data := UserInfoTab{
		GmtModified: time.Now(),
	}

	if len(path) != 0 {
		data.Image = path
	}

	if len(nickname) != 0 {
		data.NickName = nickname
	}

	if err := m.db.Model(&UserInfoTab{}).Where("user_name = ?", name).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
