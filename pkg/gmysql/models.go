package gmysql

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DBInterface interface {
	InsertUser(name, password string) error
	CheckAuth(username, password string) (bool, string, string, error)
	GetProfile(name string) (string, string, error)
	EditProfile(name, path, nickname string) error
}

type User struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	Nickname   string `json:"nickname"`
	ImagePath  string `json:"image_path"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
