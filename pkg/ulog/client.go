package ulog

import (
	"context"
	"fmt"
	"market_aggregate/pkg/util"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type Logger struct {
	// logger *logx.RotateLogger

	logger logx.Logger
}

type LogConfig struct {
	delimiter string
	days      int
	gzip      bool
}

func (l *Logger) Init(ctx context.Context, conf logx.LogConf) error {

	// logger_rule := logx.DefaultRotateRule(util.CurrTimeString(), config.delimiter, config.days, config.gzip)

	// var err error
	// l.logger, err = logx.NewLogger(util.CurrTimeString(), logger_rule, false)
	logx.MustSetup(conf)
	l.logger = logx.WithContext(ctx)
	// l.logger.MustSetUp(conf)

	return nil
}

var SingleLogger *Logger
var lock = &sync.Mutex{}

func Info(v ...interface{}) {
	LOGGER().logger.Info(v)
}

func Infof(format string, v ...interface{}) {
	LOGGER().logger.Infof(format, v)
}

func Error(v ...interface{}) {
	LOGGER().logger.Error(v)
}

func Errorf(format string, v ...interface{}) {
	LOGGER().logger.Errorf(format, v)
}

// func Severe(v ...interface{}) {
// 	LOGGER().logger.Severe(v)
// 	logx.Severe()
// }

// func Severef(format string, v ...interface{}) {
// 	LOGGER().logger.Severef(format, v)
// }

func Slow(v ...interface{}) {
	LOGGER().logger.Slow(v)
}

func Slowf(format string, v ...interface{}) {
	LOGGER().logger.Slowf(format, v)
}

func LOGGER() *Logger {
	if SingleLogger == nil {
		lock.Lock()
		defer lock.Unlock()

		if SingleLogger == nil {
			SingleLogger = new(Logger)
			fmt.Println("Init Single Logger")
		} else {
			fmt.Println("Second Judge Logger")
		}
	} else {
		// fmt.Println("Single already created!")
	}

	return SingleLogger
}

// func LOG_INIT(log_config *LogConfig) error {
// 	return LOGGER().Init(log_config)
// }

func LOG_INIT(ctx context.Context, conf logx.LogConf) error {
	return LOGGER().Init(ctx, conf)
}

func write_log1() {
	for {
		Info("[1] " + util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log2() {
	for {
		// Info("[2] " + util.CurrTimeString())

		Infof("f[2] %s ", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log3() {
	for {
		Errorf("f[3] %s", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

func write_log4() {
	for {
		Slowf("[4] %s", util.CurrTimeString())
		time.Sleep(time.Second * 1)
	}
}

/*
	ServiceName         string `json:",optional"`
	Mode                string `json:",default=console,options=[console,file,volume]"`
	Encoding            string `json:",default=json,options=[json,plain]"`
	TimeFormat          string `json:",optional"`
	Path                string `json:",default=logs"`
	Level               string `json:",default=info,options=[info,error,severe]"`
	Compress            bool   `json:",optional"`
	KeepDays            int    `json:",optional"`
	StackCooldownMillis int    `json:",default=100"`

  Compress: true
  KeepDays: 0
  Level: "info"
  Mode: "file"
  #Mode: "console"
  Path: "./log"
  ServiceName: "client"
  StackCooldownMillis: 100
*/

func get_test_conf() logx.LogConf {
	return logx.LogConf{
		ServiceName: "client",
		Mode:        "file",
		Encoding:    "json",
		TimeFormat:  "2006-01-02 15:04:05",
		Path:        "logs",
		Level:       "error",
		Compress:    true,
		KeepDays:    3,
	}
}

func TestLog1() {
	ctx := context.Background()
	conf := get_test_conf()

	LOG_INIT(ctx, conf)

	go write_log1()

	go write_log2()

	go write_log3()

	go write_log4()

	select {}
}
