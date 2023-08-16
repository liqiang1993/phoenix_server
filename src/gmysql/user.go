package gmysql

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lucky-cheerful-man/phoenix_server/src/log"
)

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
