package logger

import (
	"fmt"
	"os"
	"strings"

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

type Formatter struct{}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var out string
	level := strings.ToUpper(entry.Level.String()[:4])
	timeFmt := entry.Time.Format("2006-01-02 15:03:04.000")
	if entry.HasCaller() {
		if strings.Contains(entry.Caller.Function, "middleware") {
			out = fmt.Sprintf("[%s][%s] | [GIN] |%s\n", level, timeFmt, entry.Message)
			return []byte(out), nil
		}
		index := strings.LastIndex(entry.Caller.Function, "cybervein")
		out = fmt.Sprintf("[%s][%s] | [APP] | %s(%d) | %s\n", level, timeFmt, entry.Caller.Function[index+9:len(entry.Caller.Function)], entry.Caller.Line, entry.Message)
	} else {
		out = fmt.Sprintf("[%s][%s] | [APP] | %s\n", level, timeFmt, entry.Message)
	}
	return []byte(out), nil
}
