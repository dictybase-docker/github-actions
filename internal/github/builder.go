package github

import (
	"path"
	"strings"
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

func (b *ChangedFilesBuilder) FilterUniqueByName() *ChangedFilesBuilder {
	if len(b.files) <= 1 {
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

func (b *ChangedFilesBuilder) List() []string {
	var sl []string
	for _, v := range b.files {
		sl = append(sl, v.Name)
	}
	return sl
}

func FileNames(s []string) []string {
	var a []string
	for _, f := range s {
		a = append(a, path.Base(f))
	}
	return a
}
