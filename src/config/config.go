package config

import (
	"fmt"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	RuntimeRootPath string
	Salt            string
	RootPictureDir  string

	LogLevel      string // 日志文件展示级别
	LogFileName   string // 日志文件存放路径与名称
	LogMaxSize    int    // 日志文件大小，单位是 MB
	LogMaxBackups int    // 最大过期日志保留个数
	LogMaxAgeDay  int    // 保留过期文件最大时间，单位 天
	LogCompress   bool   // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
}

type Server struct {
	RegisterAddress       string
	RegisterServerName    string
	RegisterServerVersion string
}

type Database struct {
	Type              string
	User              string
	Password          string
	Host              string
	Name              string
	TablePrefix       string
	MaxIdleConn       int
	MaxOpenConn       int
	ConnMaxLifeMinute int
}

type Redis struct {
	Host             string
	Password         string
	MaxIdle          int
	MaxActive        int
	IdleTimeout      time.Duration
	ExpireTimeSecond int
}

type GlobalConfig struct {
	cfg             *ini.File
	AppSetting      App
	ServerSetting   Server
	DatabaseSetting Database
	RedisSetting    Redis
}

var globalConfig *GlobalConfig

// InitConfig 初始化配置内容
func init() {
	var err error
	globalConfig = new(GlobalConfig)
	//globalConfig.cfg, err = ini.Load("/Users/qiang.lilq/interesting/go/phoenix_server/localconf/app.ini")
	globalConfig.cfg, err = ini.Load("/app/localconf/app.ini")
	if err != nil {
		fmt.Printf("config.Setup, fail to parse 'localconf/app.ini': %s", err)
		panic("load config failed")
	}

	mapTo("app", &globalConfig.AppSetting)
	mapTo("server", &globalConfig.ServerSetting)
	mapTo("database", &globalConfig.DatabaseSetting)
	mapTo("redis", &globalConfig.RedisSetting)

	globalConfig.RedisSetting.IdleTimeout *= time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := globalConfig.cfg.Section(section).MapTo(v)
	if err != nil {
		fmt.Printf("cfg.MapTo %s err: %s", section, err)
		panic("get filed failed")
	}
}

func ReferGlobalConfig() *GlobalConfig {
	return globalConfig
}
