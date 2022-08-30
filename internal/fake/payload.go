package fake

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func OntoReportWithEmptyError() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get current dir %s", err)
	}

	return filepath.Join(
		filepath.Dir(dir),
		"../testdata",
		"pheno_report.json",
	), nil
}

func OntoErrorFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get current dir %s", err)
	}

	return filepath.Join(
		filepath.Dir(dir),
		"../testdata",
		"report.json",
	), nil
}

func PullReqPayload(name string) (io.Reader, error) {
	var r io.Reader
	dir, err := os.Getwd()
	if err != nil {
		return r, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir),
		"../testdata",
		name,
	)
	fhd, err := os.Open(path)
	if err != nil {
		return fhd, fmt.Errorf("error in opening file %s", err)
	}

	return fhd, nil
}

func PushPayload() (io.Reader, error) {
	var r io.Reader
	dir, err := os.Getwd()
	if err != nil {
		return r, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(filepath.Dir(dir), "../testdata", "push.json")
	fhd, err := os.Open(path)
	if err != nil {
		return fhd, fmt.Errorf("error in opening file %s", err)
	}

	return fhd, nil
}
