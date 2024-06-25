package dagger

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

const (
	owner = "dagger"
	repo  = "dagger"
)

func SetupDagger(clt *cli.Context) error {
	dver, err := fetchDaggerVersion()
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	gclient := github.NewClient(nil)
	rel, err := fetchDaggerRelease(gclient, dver)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	checksum, err := fetchDaggerCheckSum(clt, gclient, rel)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	binFileName, err := fetchDaggerBinary(
		clt.String("dagger-file"),
		dver,
		gclient,
		rel,
	)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	fmt.Println(checksum)
	fmt.Println(binFileName)
	return nil
}

func fetchDaggerRelease(
	gclient *github.Client,
	ver string,
) (*github.RepositoryRelease, error) {
	rel, _, err := gclient.Repositories.GetReleaseByTag(
		context.Background(),
		owner, repo, ver,
	)
	if err != nil {
		return nil, handleError("error in fetching release %s", err)
	}
	return rel, nil
}

func fetchDaggerBinary(
	fileSuffix, ver string,
	gclient *github.Client,
	rel *github.RepositoryRelease,
) (string, error) {
	tarballName := fmt.Sprintf("dagger_%s_%s", ver, fileSuffix)
	idx, err := findTarballIndex(rel, tarballName)
	if err != nil {
		return "", err
	}
	reader, err := downloadReleaseAsset(gclient, rel.Assets[idx].GetID())
	if err != nil {
		return "", err
	}
	defer reader.Close()
	binDir, err := createDaggerBinDir()
	if err != nil {
		return "", err
	}
	binFileName := filepath.Join(binDir, "dagger")
	if err := extractTarball(reader, binFileName); err != nil {
		return "", err
	}
	return binFileName, nil
}

func findTarballIndex(
	rel *github.RepositoryRelease,
	tarballName string,
) (int, error) {
	idx := slices.IndexFunc(rel.Assets, func(ast *github.ReleaseAsset) bool {
		return ast.GetName() == tarballName
	})
	if idx == -1 {
		return -1, handleError(
			tarballName,
			errors.New("could not find dagger tarball file"),
		)
	}
	return idx, nil
}

func downloadReleaseAsset(
	gclient *github.Client,
	assetID int64,
) (io.ReadCloser, error) {
	reader, _, err := gclient.Repositories.DownloadReleaseAsset(
		context.Background(),
		owner, repo,
		assetID,
		http.DefaultClient,
	)
	if err != nil {
		return nil, handleError("error in downloading asset %s", err)
	}
	return reader, nil
}

func extractTarball(reader io.ReadCloser, binFileName string) error {
	uncompressedStream, err := gzip.NewReader(reader)
	if err != nil {
		return handleError("extractTarGz: NewReader failed: %w", err)
	}
	defer uncompressedStream.Close()
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return handleError("extractTarGz: Next() failed: %w", err)
		}
		if header.Name != "dagger" {
			continue
		}
		writer, err := os.OpenFile(
			binFileName,
			os.O_CREATE|os.O_RDWR,
			os.FileMode(0755),
		)
		if err != nil {
			return handleError(
				"error in creating dagger bin file in temp dir %s",
				err,
			)
		}
		for {
			_, err := io.CopyN(writer, tarReader, 1024)
			if err != nil {
				if err == io.EOF {
					break
				}
				return handleError(
					"error in writing dagger bin file in temp dir %s",
					err,
				)
			}
		}
		defer writer.Close()
	}
	return nil
}

func fetchDaggerCheckSum(
	clt *cli.Context,
	gclient *github.Client,
	rel *github.RepositoryRelease,
) (string, error) {
	idx := slices.IndexFunc(rel.Assets, func(ast *github.ReleaseAsset) bool {
		return ast.GetName() == clt.String("checksum-file")
	})
	if idx == -1 {
		return "", handleError(
			clt.String("checksum-file"),
			errors.New("could not find checksum file"),
		)
	}
	reader, _, err := gclient.Repositories.DownloadReleaseAsset(
		context.Background(),
		owner, repo,
		rel.Assets[idx].GetID(),
		http.DefaultClient,
	)
	if err != nil {
		return "", handleError("error in downloading asset %s", err)
	}
	var line string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), clt.String("dagger-file")) {
			line = scanner.Text()

			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", handleError("error in reading checksum file %s", err)
	}

	return strings.Split(line, " ")[0], nil
}

func fetchDaggerVersion() (string, error) {
	resp, err := http.Get("https://dl.dagger.io/dagger/latest_version")
	if err != nil {
		return "", handleError("error in fetching dagger version %s", err)
	}
	defer resp.Body.Close()
	bcont, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", handleError("error in reading response body", err)
	}

	return fmt.Sprintf(
		"v%s",
		RemoveInvalidControlChars(strings.Trim(string(bcont), "")),
	), nil
}

func handleError(msg string, err error) error {
	githubactions.Errorf(msg, err)

	return fmt.Errorf(msg, err)
}

func RemoveInvalidControlChars(strc string) string {
	var builder strings.Builder
	for _, rtc := range strc {
		if rtc >= 32 && rtc != 127 {
			builder.WriteRune(rtc)
		}
	}

	return builder.String()
}

func createDaggerBinDir() (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "dagger-of-dcr")
	if err != nil {
		return "", handleError("failed to create temp dir: %w", err)
	}
	binDir := filepath.Join(tempDir, "bin")
	err = os.Mkdir(binDir, 0755)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", handleError("failed to create bin subfolder: %w", err)
	}
	return binDir, nil
}