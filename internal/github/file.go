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

type ChangedFilesBuilder struct {
	files []*ChangedFiles
}

func (b *ChangedFilesBuilder) FilterSuffix(suffix string) *ChangedFilesBuilder {
	if len(b.files) == 0 {
		return b
	}
	var a []*ChangedFiles
	for _, v := range b.files {
		if strings.HasSuffix(v.Name, suffix) {
			a = append(a, v)
			continue
		}
	}
	return &ChangedFilesBuilder{files: a}
}

func (b *ChangedFilesBuilder) FilterDeleted(isDeleted bool) *ChangedFilesBuilder {
	if len(b.files) == 0 {
		return b
	}
	var a []*ChangedFiles
	for _, v := range b.files {
		if v.Change == "deleted" {
			continue
		}
		a = append(a, v)
	}
	return &ChangedFilesBuilder{files: a}
}

func (b *ChangedFilesBuilder) UniqueFiles() *ChangedFilesBuilder {
	if len(b.files) == 0 {
		return b
	}
	if len(b.files) == 1 {
		return b
	}
	m := make(map[string]int)
	var a []*ChangedFiles
	for _, v := range b.files {
		n := path.Base(v.Name)
		if _, ok := m[n]; !ok {
			a = append(a, v)
			m[n] = 1
		}
	}
	return &ChangedFilesBuilder{files: a}
}

func CommittedFiles(event *github.CommitsComparison) *ChangedFilesBuilder {
	var a []*ChangedFiles
	for _, f := range event.Files {
		a = append(
			a,
			&ChangedFiles{Name: f.GetFilename(), Change: f.GetStatus()},
		)
	}
	return &ChangedFilesBuilder{files: a}
}
