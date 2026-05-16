package logger

import (
	"flag"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func newLoggerCtx(logFormat, logLevel string) *cli.Context {
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "log-format", Value: "json"},
		cli.StringFlag{Name: "log-level", Value: "info"},
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
			logFormat:     "json",
			logLevel:      "info",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: "02/Jan/2006:15:04:05"},
			wantLevel:     logrus.InfoLevel,
		},
		{
			name:          "text formatter",
			logFormat:     "text",
			logLevel:      "info",
			wantFormatter: &logrus.TextFormatter{TimestampFormat: "02/Jan/2006:15:04:05"},
			wantLevel:     logrus.InfoLevel,
		},
		{
			name:          "debug level",
			logFormat:     "json",
			logLevel:      "debug",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: "02/Jan/2006:15:04:05"},
			wantLevel:     logrus.DebugLevel,
		},
		{
			name:          "warn level",
			logFormat:     "json",
			logLevel:      "warn",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: "02/Jan/2006:15:04:05"},
			wantLevel:     logrus.WarnLevel,
		},
		{
			name:          "error level",
			logFormat:     "json",
			logLevel:      "error",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: "02/Jan/2006:15:04:05"},
			wantLevel:     logrus.ErrorLevel,
		},
		{
			name:          "fatal level",
			logFormat:     "json",
			logLevel:      "fatal",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: "02/Jan/2006:15:04:05"},
			wantLevel:     logrus.FatalLevel,
		},
		{
			name:          "panic level",
			logFormat:     "json",
			logLevel:      "panic",
			wantFormatter: &logrus.JSONFormatter{TimestampFormat: "02/Jan/2006:15:04:05"},
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
