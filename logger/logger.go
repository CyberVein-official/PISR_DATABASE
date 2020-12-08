package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	customFormatter := new(Formatter)
	Log.SetReportCaller(true)
	Log.SetFormatter(customFormatter)
	f, err := os.OpenFile("../log/server.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		Log.Error(err)
		return
	}
	Log.SetOutput(f)
	Log.SetLevel(logrus.DebugLevel)
}
