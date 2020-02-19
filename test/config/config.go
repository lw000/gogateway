package config

import (
	"encoding/json"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"time"
)

type wsConfig struct {
	Host   string `json:"host"`
	Scheme string `json:"scheme"`
	Path   string `json:"path"`
}
type JsonConfig struct {
	Debug       int64    `json:"debug"`
	WsConf      wsConfig `json:"ws"`
	Count       int      `json:"count"`
	Millisecond int      `json:"millisecond"`
	Send        bool     `json:"send"`
}

func New() *JsonConfig {
	return &JsonConfig{}
}

// 配置日志文件
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d_%H%M",
		// rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		// rotatelogs.WithRotationCount(365),  // 最多存365个文件
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	log.AddHook(lfHook)

	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

func LoadJsonConfig(file string) (*JsonConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg JsonConfig
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, err
}
