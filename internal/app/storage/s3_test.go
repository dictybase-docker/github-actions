package storage

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestGetS3Host(t *testing.T) {
	t.Parallel()

	t.Run("with port", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("s3-server", "minio.example.com", "")
		flagSet.String("s3-server-port", "9000", "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "s3-server"},
			cli.StringFlag{Name: "s3-server-port"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		host := getS3Host(ctx)
		assert.Equal(t, "minio.example.com:9000", host)
	})

	t.Run("without port", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("s3-server", "minio.example.com", "")
		flagSet.String("s3-server-port", "", "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "s3-server"},
			cli.StringFlag{Name: "s3-server-port"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		host := getS3Host(ctx)
		assert.Equal(t, "minio.example.com", host)
	})
}
