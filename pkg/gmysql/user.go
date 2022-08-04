package gmysql

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm" //nolint:goimports
	"github.com/lucky-cheerful-man/phoenix_server/pkg/log"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/setting"
	"github.com/lucky-cheerful-man/phoenix_server/pkg/util"
	"time"
)

// Setup initializes the database instance
func Setup() *Mysql {
	var MysqlOperate Mysql

	db, err := gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("gmysql.Setup err: %s", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(setting.DatabaseSetting.MaxIdleConn)
	db.DB().SetMaxOpenConns(setting.DatabaseSetting.MaxOpenConn)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(setting.DatabaseSetting.ConnMaxLifeMinute))
	MysqlOperate.db = db

	return &MysqlOperate
}

type Mysql struct {
	db *gorm.DB
}

// InsertUser 注册接口
func (m *Mysql) InsertUser(name, password string) error {
	user := User{
		Name:       name,
		Password:   password,
		Nickname:   name,
		CreateTime: util.GetCurrentStr(),
		UpdateTime: util.GetCurrentStr(),
	}
	if err := m.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// CheckAuth 登陆认证接口
func (m *Mysql) CheckAuth(username, password string) (bool, string, string, error) {
	var auth User
	err := m.db.Select([]string{"id", "nickname", "image_path"}).Where(User{Name: username,
		Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, "", "", err
	}

	if len(auth.Nickname) > 0 {
		return true, auth.Nickname, auth.ImagePath, nil
	}

	return false, "", "", nil
}

// GetProfile 查询用户的属性信息
func (m *Mysql) GetProfile(name string) (string, string, error) {
	var auth User
	err := m.db.Select([]string{"nickname", "image_path"}).Where(User{Name: name}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", "", err
	}

	if len(auth.Nickname) > 0 {
		return auth.Nickname, auth.ImagePath, nil
	} else {
		log.Warnf("can not find valid profile info by name:%s", name)
		return "", "", errors.New("not found")
	}
}

// EditProfile 编辑用户属性
func (m *Mysql) EditProfile(name, path, nickname string) error {
	data := User{
		UpdateTime: util.GetCurrentStr(),
	}

	if len(path) != 0 {
		data.ImagePath = path
	}

	if len(nickname) != 0 {
		data.Nickname = nickname
	}

	if err := m.db.Model(&User{}).Where("name = ?", name).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
