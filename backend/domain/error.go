package domain

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

type Error struct {
	technicalValue error // Техническая ошибка
	userValue      error // Пользовательская ошибка
	httpCode       int
	extraCode      int
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.userValue != nil {
		return e.userValue.Error()
	}
	if e.technicalValue != nil {
		return e.technicalValue.Error()
	}
	return "error"
}

func ErrUnauthorized(msg string) *Error {
	return &Error{userValue: errors.New(msg), httpCode: http.StatusUnauthorized}
}

func ErrForbidden(msg string) *Error {
	return &Error{userValue: errors.New(msg), httpCode: http.StatusForbidden}
}

func ErrBadRequest(msg string) *Error {
	return &Error{userValue: errors.New(msg), httpCode: http.StatusBadRequest}
}

func (e *Error) Message(debug bool) string {
	if debug {
		if e.userValue == nil {
			return e.technicalValue.Error()
		}

		if e.technicalValue == nil {
			return e.userValue.Error()
		}

		return fmt.Sprintf("%s: %s", e.userValue, e.technicalValue)
	}

	if e.userValue == nil || e.userValue.Error() == "" {
		return e.technicalValue.Error()
	}

	return e.userValue.Error()
}

func (e *Error) HttpCode() int {
	return e.httpCode
}

func (e *Error) ExtraCode() int {
	return e.extraCode
}
