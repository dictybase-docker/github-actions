package github

import (
	"path"
	"strings"

	"github.com/google/go-github/v32/github"
)

type ChangedFiles struct {
	Name   string
	Change string
}

func FilterSuffix(sl []string, suffix string) []string {
	var a []string
	for _, v := range sl {
		if strings.HasSuffix(v, suffix) {
			a = append(a, v)
			continue
		}
	}
	return a
}

func FilterDeleted(sl []string, isDeleted bool) []string {

}

func UniqueFiles(sl []string) []string {
	if len(sl) == 1 {
		return sl
	}
	m := make(map[string]int)
	var a []string
	for _, v := range sl {
		n := path.Base(v)
		if _, ok := m[n]; !ok {
			a = append(a, v)
			m[n] = 1
		}
	}
	return a
}

func CommittedFiles(event *github.CommitsComparison) []*ChangedFiles {
	var files []*ChangedFiles
	for _, f := range event.Files {
		files = append(
			files,
			&ChangedFiles{Name: f.GetFilename(), Change: f.GetStatus()},
		)
	}
	return files
}
