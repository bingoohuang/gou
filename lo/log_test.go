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
	viper.Set("PrintCallerInfo", true)

	viper.Set("logMaxBackups", 3)
	viper.Set("logTimeFormat", "20060102150405")

	SetupLog()

	for i := 0; i < 5; i++ {
		logrus.Info("abc", "efg")
		logrus.Info("abc", "efg")
		logrus.Info("abc", "efg")
		logrus.Info("abc", "efg")
		logrus.Info("abc", "efg")
		logrus.Info("abc", "efg")

		time.Sleep(time.Second)
	}
}
