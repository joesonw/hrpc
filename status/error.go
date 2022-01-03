package status

import (
	"errors"
	"fmt"
	"net/http"
)

type Status struct {
	httpStatus int
	err        error
}

func (e *Status) Unwrap() error {
	return e.err
}

func (e *Status) Error() string {
	return e.err.Error()
}

func (e *Status) Status() int {
	return e.httpStatus
}

func New(httpStatus int, msg string) *Status {
	return &Status{
		httpStatus: httpStatus,
		err:        errors.New(msg),
	}
}

func Wrap(httpStatus int, err error) *Status {
	return &Status{
		httpStatus: httpStatus,
		err:        err,
	}
}

func Errorf(httpStatus int, format string, a ...interface{}) *Status {
	return Wrap(httpStatus, fmt.Errorf(format, a...))
}

func IsStatus(e error) bool {
	var s *Status
	return errors.As(e, &s)
}

func HTTPStatus(e error) int {
	var s *Status
	if !errors.As(e, &s) {
		return http.StatusInternalServerError
	}
	return s.Status()
}
