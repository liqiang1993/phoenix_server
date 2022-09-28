package gmysql

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type DBInterface interface {
	InsertUser(name, password string) error
	CheckAuth(username, password string) (bool, string, string, error)
	GetProfile(name string) (string, string, error)
	EditProfile(name, path, nickname string) error
}

type UserInfoTab struct {
	UserName    string    `json:"user_name"`
	Password    string    `json:"password"`
	NickName    string    `json:"nick_name"`
	Image       string    `json:"image"`
	GmtCreated  time.Time `json:"gmt_created"`
	GmtModified time.Time `json:"gmt_modified"`
}
