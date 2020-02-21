package lo

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"

	"github.com/bingoohuang/gou/str"
	"github.com/spf13/pflag"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DeclareLogPFlags declares the log pflags.
func DeclareLogPFlags() {
	pflag.StringP("loglevel", "", "info", "debug/info/warn/error")
	pflag.StringP("logdir", "", "var/logs", "log dir")
	pflag.BoolP("logrus", "", true, "enable logrus")
}

// DeclareLogFlags declares the log flags.
func DeclareLogFlags() {
	flag.String("loglevel", "info", "debug/info/warn/error")
	flag.String("logdir", "var/logs", "log dir")
	flag.Bool("logrus", true, "enable logrus")
}

// TextFormatter extends the prefixed.TextFormatter with line joining.
type TextFormatter struct {
	Skip            int
	PrintCallerInfo bool
}

var reNewLines = regexp.MustCompile(`\r?\n`) // nolint

// Format formats the log output.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := bytes.Buffer{}

	b.WriteString(entry.Time.Format("2006-01-02 15:04:05.000") + " ")

	if f.PrintCallerInfo {
		// getting caller info - it's expensive.
		if _, file, line, ok := runtime.Caller(f.Skip); ok {
			//funcName := runtime.FuncForPC(pc).Name()
			b.WriteString(fmt.Sprintf("%s:%d ", path.Base(file), line))
		}
	}

	b.WriteString("[" + entry.Level.String() + "] ")
	b.WriteString(reNewLines.ReplaceAllString(entry.Message, ``) + "\n")

	return b.Bytes(), nil
}

// SetupLog setup log parameters.
func SetupLog() io.Writer {
	loglevel := viper.GetString("loglevel")
	logrus.SetLevel(str.Decode(loglevel,
		"debug", logrus.DebugLevel, "info", logrus.InfoLevel, "warn", logrus.WarnLevel,
		"error", logrus.ErrorLevel, logrus.InfoLevel).(logrus.Level))

	viper.SetDefault("ContextHookSkip", 7)

	// https://stackoverflow.com/a/48972299
	formatter := &TextFormatter{
		Skip:            viper.GetInt("ContextHookSkip"),
		PrintCallerInfo: viper.GetBool("PrintCallerInfo"),
	}

	if !viper.GetBool("logrus") {
		logrus.SetFormatter(formatter)

		return os.Stdout
	}

	logdir := viper.GetString("logdir")
	if logdir != "" {
		if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
			logrus.Panicf("failed to create %s error %v\n", logdir, err)
		}

		return initLogger(str.Decode(loglevel,
			"debug", logrus.DebugLevel, "info", logrus.InfoLevel, "warn", logrus.WarnLevel,
			"error", logrus.ErrorLevel, logrus.InfoLevel).(logrus.Level), logdir, filepath.Base(os.Args[0])+".log", formatter)
	}

	logrus.SetFormatter(formatter)

	return os.Stdout
}

// 参考链接： https://tech.mojotv.cn/2018/12/27/golang-logrus-tutorial
// nolint gomnd
func initLogger(level logrus.Level, logDir, filename string, formatter logrus.Formatter) io.Writer {
	baseLogPath := path.Join(logDir, filename)
	writer, err := NewDailyFile(baseLogPath, 7) // 文件最大保存7天

	if err != nil {
		logrus.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}

	logrus.SetLevel(level)

	//writerMap := lfshook.WriterMap{
	//	logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
	//	logrus.InfoLevel:  writer,
	//	logrus.WarnLevel:  writer,
	//	logrus.ErrorLevel: writer,
	//	logrus.FatalLevel: writer,
	//	logrus.PanicLevel: writer,
	//}

	writerMap := writer

	logrus.AddHook(lfshook.NewHook(writerMap, formatter))
	logrus.SetOutput(&Discarder{})

	return writer
}

// A Discarder sends all writes to ioutil.Discard.
type Discarder struct{}

// Write implements io.Writer.
func (d *Discarder) Write(b []byte) (int, error) { return ioutil.Discard.Write(b) }
