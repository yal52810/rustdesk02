package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
)

type Config struct {
	Path         string
	Level        string
	ReportCaller bool
}

func New(c *Config) *log.Logger {
	log.SetFormatter(&nested.Formatter{
		// HideKeys:        true,
		TimestampFormat: "[2006-01-02 15:04:05]",
		NoColors:        true,
		NoFieldsColors:  true,
		//FieldsOrder:     []string{"name", "age"},
	})

	// 强制输出到 stdout，忽略配置文件中的日志路径
	// 这是为了避免容器挂载点冲突问题
	log.SetOutput(os.Stdout)

	log.SetReportCaller(c.ReportCaller)

	level, err2 := log.ParseLevel(c.Level)
	if err2 != nil {
		level = log.DebugLevel
	}
	log.SetLevel(level)

	return log.StandardLogger()
}
