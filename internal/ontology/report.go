package ontology

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

func ParseViolations(path string, level string) ([]string, error) {
	cont, err := gabs.ParseJSONFile(path)
	if err != nil {
		return []string{}, fmt.Errorf("error in parsing json file %s", err)
	}
	hasLevel := false
	var violCont *gabs.Container
	for _, child := range cont.S("").Children() {
		if !child.ExistsP(level) {
			continue
		}
		violCont = child
		hasLevel = true
		break
	}
	if !hasLevel {
		return []string{}, &ViolationNotFound{Level: level}
	}
	var s []string
	for k := range violCont.S("violations").ChildrenMap() {
		s = append(s, strings.ReplaceAll(k, "_", " "))
	}
	return s, nil
}
