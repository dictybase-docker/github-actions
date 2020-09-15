package github

import (
	"path"

	"github.com/google/go-github/v32/github"
)

func uniqueFiles(sl []string) []string {
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

func CommittedFiles(event *github.CommitsComparison, skipDeleted bool) []string {
	var files []string
	for _, f := range event.Files {
		if skipDeleted {
			if f.GetStatus() == "deleted" {
				continue
			}
		}
		files = append(files, f.GetFilename())
	}
	return uniqueFiles(files)
}
