package log

import (
	"github.com/lucky-cheerful-man/phoenix_server/src/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = logrus.New()

func init() {
	logger := &lumberjack.Logger{
		Filename:   config.ReferGlobalConfig().AppSetting.LogFileName,
		MaxSize:    config.ReferGlobalConfig().AppSetting.LogMaxSize,    // 日志文件大小，单位是 MB
		MaxBackups: config.ReferGlobalConfig().AppSetting.LogMaxBackups, // 最大过期日志保留个数
		MaxAge:     config.ReferGlobalConfig().AppSetting.LogMaxAgeDay,  // 保留过期文件最大时间，单位 天
		Compress:   config.ReferGlobalConfig().AppSetting.LogCompress,   // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
	}

	log.SetOutput(logger)
	log.SetLevel(getLevel(config.ReferGlobalConfig().AppSetting.LogLevel))
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

// Info 封装一层info日志
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn 封装一层warn日志
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error 封装一层error日志
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
