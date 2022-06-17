package log

import (
	"entry-task-rpc/pkg/setting"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = logrus.New()

func InitLog() {
	logger := &lumberjack.Logger{
		Filename:   setting.AppSetting.LogFileName,
		MaxSize:    setting.AppSetting.LogMaxSize,    // 日志文件大小，单位是 MB
		MaxBackups: setting.AppSetting.LogMaxBackups, // 最大过期日志保留个数
		MaxAge:     setting.AppSetting.LogMaxAgeDay,  // 保留过期文件最大时间，单位 天
		Compress:   setting.AppSetting.LogCompress,   // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
	}

	log.SetOutput(logger)
	log.SetLevel(getLevel(setting.AppSetting.LogLevel))
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
}

// 获取本次启用的日志级别
func getLevel(level string) logrus.Level {
	switch level {
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	}

	return logrus.InfoLevel
}

// Infof 封装一层info日志
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf 封装一层warn日志
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf 封装一层error日志
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf 封装一层Fatal日志
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
