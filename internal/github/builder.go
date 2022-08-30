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
	var afl []*ChangedFiles
	for _, v := range b.files {
		if strings.HasSuffix(v.Name, suffix) {
			afl = append(afl, v)

			continue
		}
	}

	return &ChangedFilesBuilder{files: afl}
}

func (b *ChangedFilesBuilder) FilterDeleted(
	isDeleted bool,
) *ChangedFilesBuilder {
	if len(b.files) == 0 {
		return b
	}
	afl := make([]*ChangedFiles, 0)
	for _, vfl := range b.files {
		if vfl.Change == "deleted" {
			continue
		}
		afl = append(afl, vfl)
	}

	return &ChangedFilesBuilder{files: afl}
}

func (b *ChangedFilesBuilder) FilterUniqueByName() *ChangedFilesBuilder {
	if len(b.files) <= 1 {
		return b
	}
	mnt := make(map[string]int)
	afl := make([]*ChangedFiles, 0)
	for _, v := range b.files {
		n := path.Base(v.Name)
		if _, ok := mnt[n]; !ok {
			afl = append(afl, v)
			mnt[n] = 1
		}
	}

	return &ChangedFilesBuilder{files: afl}
}

func (b *ChangedFilesBuilder) List() []string {
	slc := make([]string, 0)
	for _, v := range b.files {
		slc = append(slc, v.Name)
	}

	return slc
}

func FileNames(s []string) []string {
	afl := make([]string, 0)
	for _, f := range s {
		afl = append(afl, path.Base(f))
	}

	return afl
}
