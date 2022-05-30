package data

import "fmt"

type RecordNotFoundError struct{}

func (e *RecordNotFoundError) Error() string {
	return "Not Found"
}

type InternalError struct {
	err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error occurred: %v", e.err)
}
