package ontology

import (
	"errors"
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
	for _, child := range cont.Children() {
		v, ok := child.Search("level").Data().(string)
		if !ok {
			return []string{""}, errors.New("incompatible report format, level key not found")
		}
		if v != level {
			continue
		}
		violCont = child
		hasLevel = true

		break
	}
	children := violCont.S("violations").Children()
	if !hasLevel || len(children) == 0 {
		return []string{}, &ViolationNotFoundError{Level: level}
	}
	var slc []string
	for _, child := range children {
		for k := range child.ChildrenMap() {
			slc = append(slc, strings.ReplaceAll(k, "_", " "))
		}
	}

	return slc, nil
}
