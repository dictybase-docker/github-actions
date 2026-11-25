package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	jsonFormat = "json"
	timeFormat = "02/Jan/2006:15:04:05"
)

func GetLogger(clt *cli.Context) *logrus.Entry {
	log := logrus.New()
	log.Out = os.Stderr

	switch clt.GlobalString("log-format") {
	case "text":
		log.Formatter = &logrus.TextFormatter{
			TimestampFormat: timeFormat,
		}
	case jsonFormat:
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: timeFormat,
		}
	}

	l := clt.GlobalString("log-level")
	switch l {
	case "debug":
		log.Level = logrus.DebugLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	}

	return logrus.NewEntry(log)
}
