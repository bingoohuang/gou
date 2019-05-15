package gou

import (
	"bufio"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// 参考链接： https://tech.mojotv.cn/2018/12/27/golang-logrus-tutorial
func InitLogger(logLevel, logDir, filename string) {
	baseLogPath := path.Join(logDir, filename)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}

	// 如果日志级别不是debug就不要打印日志到控制台了
	if logLevel == "debug" {
		logrus.SetOutput(os.Stderr)
	} else {
		setNull()
	}

	level := Decode(logLevel,
		"debug", logrus.DebugLevel,
		"info", logrus.InfoLevel,
		"warn", logrus.WarnLevel,
		"error", logrus.ErrorLevel,
		logrus.InfoLevel).(logrus.Level)
	logrus.SetLevel(level)

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{})
	logrus.AddHook(lfHook)
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	logrus.SetOutput(writer)
}
