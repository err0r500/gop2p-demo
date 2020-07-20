package domain

import (
	"fmt"
)

// todo maybe refactor this to have more generic types (like notFound, conflict...)
// struct is actually only really needed for ErrMalformed

// ErrResourceNotFound is used when a resource is not found
type ErrResourceNotFound struct{}

func (ErrResourceNotFound) Error() string { return "user not found" }

// ErrTechnical is used when a tech error happens
type ErrTechnical struct{}

func (ErrTechnical) Error() string { return "a technical error happened" }

// ErrMalformed is used when invalid params are provided to usecases
type ErrMalformed struct {
	Details []string
}

func (e ErrMalformed) Error() string {
	return fmt.Sprintf("%v", e.Details)
}
