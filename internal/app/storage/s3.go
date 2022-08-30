package storage

import (
	"fmt"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/minio/minio-go"
	"github.com/urfave/cli"
)

func getS3Host(clt *cli.Context) string {
	if len(clt.String("s3-server-port")) > 0 {
		return fmt.Sprintf(
			"%s:%s",
			clt.String("s3-server"),
			clt.String("s3-server-port"),
		)
	}

	return clt.String("s3-server")
}

func SaveInS3(clt *cli.Context) error {
	s3Client, err := minio.New(
		getS3Host(clt),
		clt.String("access-key"),
		clt.String("secret-key"),
		true,
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting minio client %s", err),
			2,
		)
	}
	log := logger.GetLogger(clt)
	path := clt.String("upload-path")
	if len(path) == 0 {
		path = clt.String("input")
	}
	log.Debugf("upload path %s", path)
	_, err = s3Client.FPutObject(
		clt.String("s3-bucket"),
		path,
		clt.String("input"),
		minio.PutObjectOptions{ContentType: "application/text"},
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to upload file %s", err),
			2,
		)
	}
	log.Infof("save file %s to s3 storage", clt.String("input"))

	return nil
}
