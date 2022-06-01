package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

var logToFile *logrus.Logger

//日志文件名
var loggerFile string

func SetLogFile(file string) {
	loggerFile = file
}

//初始化
func init() {
	// SetLogFile(filepath.Join(conf.Conf.MyLog.Path, conf.Conf.MyLog.Name))
	SetLogFile(filepath.Join(viper.GetString("myLog.path"), viper.GetString("myLog.name")))
}

func Log() *logrus.Logger {
	//文件输出
	if viper.GetString("myLog.model") == "file" {
		return logFile()
	} else {
		//控制台输出
		if log == nil {
			log = logrus.New()
			log.Out = os.Stdout
			log.Formatter = &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}
			log.SetLevel(logrus.DebugLevel)
		}
	}
	return log
}

func logFile() *logrus.Logger {

	if logToFile == nil {
		logToFile = logrus.New()

		logToFile.SetLevel(logrus.DebugLevel)

		// 设置 rotatelogs  返回写日志对象logWriter
		logWriter, _ := rotatelogs.New(
			// 分割后的文件名称
			loggerFile+"_%Y%m%d.log",

			// 设置最大保存时间
			rotatelogs.WithMaxAge(30*24*time.Hour),

			// 设置日志切割时间间隔(1天)
			rotatelogs.WithRotationTime(24*time.Hour),
		)

		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}

		//设置时间格式
		lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})

		// 新增 Hook
		logToFile.AddHook(lfHook)
	}
	return logToFile
}
