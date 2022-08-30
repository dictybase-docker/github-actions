package ontology

import "fmt"

type ViolationNotFoundError struct {
	Level string
}

func (v *ViolationNotFoundError) Error() string {
	return fmt.Sprintf("violation %s is not found", v.Level)
}

func IsViolationNotFound(err error) bool {
	if _, ok := err.(*ViolationNotFoundError); ok {
		return true
	}

	return false
}
