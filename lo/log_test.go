package lo

import (
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

func TestSetupLog(t *testing.T) {
	_ = os.RemoveAll("./logs")

	viper.Set(LogrusKey, true)
	viper.Set(LogDebugKey, true)
	viper.Set(LogdirKey, "./logs")

	viper.Set(LogMaxBackupsDaysKey, 3)
	viper.Set(LogTimeFormatKey, "20060102150405")

	SetupLog()

	err := errors.New("error")
	err = errors.Wrap(err, "open failed")
	err = errors.Wrap(err, "read config failed")

	// 使用`%+v`来打印日志堆栈
	logrus.Warnf("errors %+v", err)

	for i := 0; i < 5; i++ {
		logrus.WithFields(map[string]interface{}{"key1": "value10", "key2": "value\n20"}).Info("abc", "efg")
		logrus.WithFields(map[string]interface{}{"key1": "value11", "key2": "value\r21"}).Info("abc", "efg")
		logrus.WithFields(map[string]interface{}{"key1": "value12", "key2": "value'=22"}).Info("abc", "efg")
		logrus.WithFields(map[string]interface{}{"key1": "value13", "key2": "value:23"}).Info("abc", "efg")
		logrus.WithFields(map[string]interface{}{"key1": "value14", "key2": "value\"24"}).Info("abc", "efg")
		logrus.WithFields(map[string]interface{}{"key1": "value15", "key2": "value 25"}).Info("abc", "efg")

		time.Sleep(time.Second)
	}
}
