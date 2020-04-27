package lo

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/bingoohuang/gou/file"

	"github.com/bingoohuang/gou/local"
	"github.com/bingoohuang/gou/str"
	"github.com/thoas/go-funk"

	"github.com/bingoohuang/gou/lang"

	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"
)

// DeclareLogPFlags declares the log pflags.
func DeclareLogPFlags() {
	pflag.StringP(LoglevelKey, "", "info", "debug/info/warn/error")
	pflag.StringP(LogdirKey, "", "", "log dir")
	pflag.BoolP(LogrusKey, "", true, "enable log to file")
}

// DeclareLogFlags declares the log flags.
func DeclareLogFlags() {
	flag.String(LoglevelKey, "info", "debug/info/warn/error")
	flag.String(LogdirKey, "", "log dir")
	flag.Bool(LogrusKey, true, "enable log to file")
}

// TextFormatter extends the prefixed.TextFormatter with line joining.
type TextFormatter struct {
}

var reNewLines = regexp.MustCompile(`\r?\n`) // nolint

// Format formats the log output.
func (f *TextFormatter) Format(e *logrus.Entry) ([]byte, error) {
	b := bytes.Buffer{}

	b.WriteString(e.Time.Format("2006-01-02 15:04:05.000") + " ")
	// align the longest WARNING, which has the length of 7
	b.WriteString(fmt.Sprintf("%7s ", strings.ToUpper(e.Level.String())))
	b.WriteString(fmt.Sprintf("%d --- ", os.Getpid()))
	b.WriteString(fmt.Sprintf("[%d] ", lang.CurGoroutineID().Uint64()))
	b.WriteString(fmt.Sprintf("[%s] ", local.String(local.KeyTraceID)))

	caller := getCaller()

	b.WriteString(fmt.Sprintf("%-20s", fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line)))
	b.WriteString(" : ")

	if len(e.Data) > 0 {
		keys := funk.Keys(e.Data).([]string)
		sort.Strings(keys)

		for _, k := range keys {
			b.WriteString(fmt.Sprintf("%s:%s ",
				str.TryQuote(k, ":='"),
				str.TryQuote(fmt.Sprintf("%v", e.Data[k]), ":='")))
		}
	}

	b.WriteString(reNewLines.ReplaceAllString(e.Message, ``) + "\n")

	return b.Bytes(), nil
}

// SetupLog setup log parameters.
func SetupLog() io.Writer {
	l, err := logrus.ParseLevel(viper.GetString(LoglevelKey))
	if err != nil {
		l = logrus.InfoLevel
	}

	logrus.SetLevel(l)

	// https://stackoverflow.com/a/48972299
	formatter := &TextFormatter{}

	if !viper.GetBool(LogrusKey) {
		logrus.SetFormatter(formatter)

		return os.Stdout
	}

	appName := filepath.Base(os.Args[0])
	logdir := viper.GetString(LogdirKey)

	if logdir == "" {
		logdir = file.HomeDirExpand("~/logs/" + appName)
	}

	if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
		logrus.Panicf("failed to create %s error %v\n", logdir, err)
	}

	return initLogger(l, logdir, appName+".log", formatter)
}

const (
	LoglevelKey = "loglevel"
	LogdirKey   = "logdir"
	LogrusKey   = "logrus"

	LogTimeFormatKey     = "logTimeFormat"
	LogMaxBackupsDaysKey = "logMaxBackupsDays"
	LogDebugKey          = "logDebug"
)

// 参考链接： https://tech.mojotv.cn/2018/12/27/golang-logrus-tutorial
// nolint gomnd
func initLogger(level logrus.Level, logDir, filename string, formatter logrus.Formatter) io.Writer {
	viper.SetDefault(LogMaxBackupsDaysKey, 7)
	viper.SetDefault(LogDebugKey, false)
	viper.SetDefault(LogTimeFormatKey, "2006-01-02")

	writer, err := NewRotateFile(filepath.Join(logDir, filename),
		MaxBackupsDays(viper.GetInt(LogMaxBackupsDaysKey)),
		TimeFormat(viper.GetString(LogTimeFormatKey)),
		Debug(viper.GetBool(LogDebugKey)))
	if err != nil {
		logrus.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}

	logrus.SetLevel(level)
	logrus.AddHook(lfshook.NewHook(writer, formatter))
	logrus.SetOutput(ioutil.Discard)
	logrus.SetReportCaller(true)

	return writer
}
