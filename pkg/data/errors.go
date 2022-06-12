package data

import "fmt"

type RecordNotFoundError struct {
	err error
}

func (e *RecordNotFoundError) Error() string {
	return "Not Found"
}

type InternalError struct {
	err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error occurred: %v", e.err)
}

type ForeignKeyError struct {
	err error
}

func (e *ForeignKeyError) Error() string {
	return fmt.Sprintf("Foreign key constraint failed: %v", e.err)
}
