package storage

import (
	"fmt"
	"os"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/minio/minio-go"
	"github.com/urfave/cli"
)

func getS3Host(c *cli.Context) string {
	if len(c.String("s3-server-port")) > 0 {
		return fmt.Sprintf("%s:%s", c.String("s3-server"), c.String("s3-server-port"))
	}
	return c.String("s3-server")
}

func SaveInS3(c *cli.Context) error {
	s3Client, err := minio.New(
		getS3Host(c),
		c.String("access-key"),
		c.String("secret-key"),
		true,
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting minio client %s", err),
			2,
		)
	}
	reader, err := os.Open(c.String("input"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to read the file %s %s", c.String("input"), err),
			2,
		)
	}
	defer reader.Close()
	fstat, err := reader.Stat()
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to stat the file %s", err),
			2,
		)
	}
	l := logger.GetLogger(c)
	l.Debugf("upload size of file %s is %d byte", c.String("input"), fstat.Size())
	path := c.String("upload-path")
	if len(path) == 0 {
		path = c.String("input")
	}
	l.Debugf("upload path %s", path)
	_, err = s3Client.PutObject(
		c.String("s3-bucket"),
		path,
		reader,
		fstat.Size(),
		minio.PutObjectOptions{ContentType: "application/text"},
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to upload file %s", err),
			2,
		)
	}
	l.Infof("save file %s to s3 storage", c.String("input"))
	return nil
}
