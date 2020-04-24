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
	"runtime"
	"sort"
	"strings"

	"github.com/bingoohuang/gou/file"

	"github.com/bingoohuang/gou/local"
	"github.com/bingoohuang/gou/str"
	"github.com/thoas/go-funk"

	"github.com/bingoohuang/gou/lang"

	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"

	"github.com/spf13/pflag"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DeclareLogPFlags declares the log pflags.
func DeclareLogPFlags() {
	pflag.StringP("loglevel", "", "info", "debug/info/warn/error")
	pflag.StringP("logdir", "", "", "log dir")
	pflag.BoolP("logrus", "", true, "enable logrus")
}

// DeclareLogFlags declares the log flags.
func DeclareLogFlags() {
	flag.String("loglevel", "info", "debug/info/warn/error")
	flag.String("logdir", "", "log dir")
	flag.Bool("logrus", true, "enable logrus")
}

// TextFormatter extends the prefixed.TextFormatter with line joining.
type TextFormatter struct {
	Skip              int
	NoPrintCallerInfo bool
}

var reNewLines = regexp.MustCompile(`\r?\n`) // nolint

// Format formats the log output.
func (f *TextFormatter) Format(e *logrus.Entry) ([]byte, error) {
	b := bytes.Buffer{}

	b.WriteString(e.Time.Format("2006-01-02 15:04:05.000") + " ")
	b.WriteString(fmt.Sprintf("%5s ", strings.ToUpper(e.Level.String())))
	b.WriteString(fmt.Sprintf("%d --- ", os.Getpid()))
	b.WriteString(fmt.Sprintf("[%d] ", lang.CurGoroutineID().Uint64()))
	b.WriteString(fmt.Sprintf("[%s] ", local.String(local.KeyTraceID)))

	if !f.NoPrintCallerInfo {
		// getting caller info - it's expensive.
		if _, file, line, ok := runtime.Caller(f.Skip); ok {
			fl := fmt.Sprintf("%s:%d", filepath.Base(file), line)
			b.WriteString(fmt.Sprintf("%-20s", fl))
		}
	}

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
	loglevel := viper.GetString("loglevel")

	l, err := logrus.ParseLevel(loglevel)
	if err != nil {
		l = logrus.InfoLevel
	}

	logrus.SetLevel(l)

	viper.SetDefault("contextHookSkip", 8)

	// https://stackoverflow.com/a/48972299
	formatter := &TextFormatter{
		Skip:              viper.GetInt("contextHookSkip"),
		NoPrintCallerInfo: viper.GetBool("noPrintCallerInfo"),
	}

	if !viper.GetBool("logrus") {
		logrus.SetFormatter(formatter)

		return os.Stdout
	}

	appName := filepath.Base(os.Args[0])
	logdir := viper.GetString("logdir")

	if logdir == "" {
		logdir = file.HomeDirExpand("~/logs/" + appName)
	}

	if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
		logrus.Panicf("failed to create %s error %v\n", logdir, err)
	}

	return initLogger(l, logdir, appName+".log", formatter)
}

// 参考链接： https://tech.mojotv.cn/2018/12/27/golang-logrus-tutorial
// nolint gomnd
func initLogger(level logrus.Level, logDir, filename string, formatter logrus.Formatter) io.Writer {
	viper.SetDefault("logMaxBackups", 7)
	viper.SetDefault("logDebug", false)
	viper.SetDefault("logTimeFormat", "20060102")

	maxBackupsDays := viper.GetInt("logMaxBackupsDays")
	timeFormat := viper.GetString("logTimeFormat")
	logDebug := viper.GetBool("logDebug")

	writer, err := NewRotateFile(filepath.Join(logDir, filename),
		MaxBackupsDays(maxBackupsDays), TimeFormat(timeFormat), Debug(logDebug))
	if err != nil {
		logrus.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}

	logrus.SetLevel(level)
	logrus.AddHook(lfshook.NewHook(writer, formatter))
	logrus.SetOutput(ioutil.Discard)

	return writer
}
