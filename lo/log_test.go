package lo

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

func TestSetupLog(t *testing.T) {
	_ = os.RemoveAll("./logs")

	viper.Set("logrus", true)
	viper.Set("logDebug", true)
	viper.Set("logdir", "./logs")
	viper.Set("contextHookSkip", 8)

	viper.Set("logMaxBackups", 3)
	viper.Set("logTimeFormat", "20060102150405")

	SetupLog()

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
