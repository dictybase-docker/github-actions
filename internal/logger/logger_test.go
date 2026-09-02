package logger

import (
	"flag"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

const infoLevelName = "info"

func newLoggerCtx(logFormat, logLevel string) *cli.Context {
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "log-format", Value: jsonFormat},
		cli.StringFlag{Name: "log-level", Value: infoLevelName},
	}

	flagSet.String("log-format", logFormat, "")
	flagSet.String("log-level", logLevel, "")
	_ = flagSet.Parse([]string{})

	return cli.NewContext(app, flagSet, nil)
}

func TestGetLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		logFormat     string
		logLevel      string
		wantFormatter any
		wantLevel     logrus.Level
	}{
		{
			name:          "default",
			logFormat:     jsonFormat,
			logLevel:      infoLevelName,
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: timeFormat},
			wantLevel:     logrus.InfoLevel,
		},
		{
			name:          "text formatter",
			logFormat:     "text",
			logLevel:      infoLevelName,
			wantFormatter: &logrus.TextFormatter{TimestampFormat: timeFormat},
			wantLevel:     logrus.InfoLevel,
		},
		{
			name:          "debug level",
			logFormat:     jsonFormat,
			logLevel:      "debug",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: timeFormat},
			wantLevel:     logrus.DebugLevel,
		},
		{
			name:          "warn level",
			logFormat:     jsonFormat,
			logLevel:      "warn",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: timeFormat},
			wantLevel:     logrus.WarnLevel,
		},
		{
			name:          "error level",
			logFormat:     jsonFormat,
			logLevel:      "error",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: timeFormat},
			wantLevel:     logrus.ErrorLevel,
		},
		{
			name:          "fatal level",
			logFormat:     jsonFormat,
			logLevel:      "fatal",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: timeFormat},
			wantLevel:     logrus.FatalLevel,
		},
		{
			name:          "panic level",
			logFormat:     jsonFormat,
			logLevel:      "panic",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: timeFormat},
			wantLevel:     logrus.PanicLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newLoggerCtx(tt.logFormat, tt.logLevel)
			entry := GetLogger(ctx)
			assert.NotNil(t, entry)
			assert.Equal(t, tt.wantLevel, entry.Logger.Level)
			assert.IsType(t, tt.wantFormatter, entry.Logger.Formatter)
		})
	}
}
