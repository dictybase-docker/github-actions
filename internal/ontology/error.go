package ontology

import "fmt"

type ViolationNotFound struct {
	Level string
}

func (v *ViolationNotFound) Error() string {
	return fmt.Sprintf("violation %s is not found", v.Level)

}

func IsViolationNotFound(err error) bool {
	if _, ok := err.(*ViolationNotFound); ok {
		return true
	}
	return false
}
