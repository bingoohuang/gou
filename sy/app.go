package sy

import (
	"io"

	"github.com/spf13/viper"

	"github.com/bingoohuang/gou/cnf"
	"github.com/bingoohuang/gou/htt"
	"github.com/bingoohuang/gou/lo"
)

// AppOption defines the application options.
type AppOption struct {
	EnvPrefix   string
	LogLevel    string
	ConfigBeans []interface{}

	LogWriter io.Writer
}

// SetupApp setup the application.
func SetupApp(appOption *AppOption) {
	lo.DeclareLogPFlags()
	cnf.DeclarePflags()
	cnf.DeclarePflagsByStruct(appOption.ConfigBeans...)

	pprofAddr := htt.PprofAddrPflag()

	if err := cnf.ParsePflags(appOption.EnvPrefix); err != nil {
		panic(err)
	}

	htt.StartPprof(*pprofAddr)

	if appOption.LogLevel != "" {
		viper.Set("loglevel", appOption.LogLevel)
	}

	appOption.LogWriter = lo.SetupLog()

	_ = UpdatePidFile()

	cnf.LoadByPflag(appOption.ConfigBeans...)
}
