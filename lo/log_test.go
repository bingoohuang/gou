package lo

import (
	"github.com/spf13/viper"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetupLog(t *testing.T) {
	viper.Set("PrintCallerInfo", true)
	SetupLog()

	logrus.Info("abc", "efg")
	logrus.Info("abc", "efg")
	logrus.Info("abc", "efg")
	logrus.Info("abc", "efg")
	logrus.Info("abc", "efg")
	logrus.Info("abc", "efg")
}
