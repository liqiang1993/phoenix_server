package setting

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

var AppSetting = &App{}

type Server struct {
	Port int
}

var ServerSetting = &Server{}

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

var DatabaseSetting = &Database{}

type Redis struct {
	Host             string
	Password         string
	MaxIdle          int
	MaxActive        int
	IdleTimeout      time.Duration
	ExpireTimeSecond int
}

var RedisSetting = &Redis{}

var cfg *ini.File

// InitConfig 初始化配置内容
func InitConfig() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		fmt.Printf("setting.Setup, fail to parse 'conf/app.ini': %s", err)
		panic("load config failed")
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	RedisSetting.IdleTimeout *= time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		fmt.Printf("cfg.MapTo %s err: %s", section, err)
		panic("get filed failed")
	}
}
