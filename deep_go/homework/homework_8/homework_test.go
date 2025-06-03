package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	var builder strings.Builder

	if len(e.errors) == 0 {
		return ""
	}
	if len(e.errors) == 1 {
		return e.errors[0].Error()
	}
	builder.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errors)))
	
	for _, err := range e.errors {
		builder.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}
	builder.WriteString("\n")
	
	return builder.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError
	
	if mulErr, ok := err.(*MultiError); ok {
		multiErr = mulErr
	} else if err != nil {
		multiErr = &MultiError{errors: []error{err}}
	} else {
		multiErr = &MultiError{errors: []error{}}
	}
	
	for _, e := range errs {
		if e != nil {
			if mulErr, ok := e.(*MultiError); ok {
				multiErr.errors = append(multiErr.errors, mulErr.errors...)
			} else {
				multiErr.errors = append(multiErr.errors, e)
			}
		}
	}
	if len(multiErr.errors) == 0 {
		return nil
	}
	
	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedmessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedmessage)
}