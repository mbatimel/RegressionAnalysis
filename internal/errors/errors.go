package errors

import (
	"fmt"
)

const (
	noneValue = "None"

	CauseErrDescription = "description"
)

type Error struct {
	ErrorText  string
	Cause      map[string]interface{}
	statusCode int

	internalError error
}

func (e *Error) SetStatusCode(code int) *Error {
	e.statusCode = code
	return e
}

func (e *Error) GetStatusCode() int {
	return e.statusCode
}

func (e *Error) SetOuterError(err interface{}) *Error {
	e.internalError = fmt.Errorf("%v", err)
	return e
}

func (e *Error) GetOuterError() error {
	return e.internalError
}

func (e *Error) Error() string {
	var cause string
	if e.Cause != nil {
		cause = fmt.Sprintf(". Causes: %v", e.Cause)
	}
	return e.ErrorText + cause
}

func Is(errOne, errTwo error) bool {
	custErrOne, ok := errOne.(*Error)
	custErrTwo, ok2 := errTwo.(*Error)

	if ok && ok2 {
		return custErrOne.ErrorText == custErrTwo.ErrorText
	} else {
		return errOne.Error() == errTwo.Error()
	}
}

func (e *Error) AddTrErrors(trError TrParams) *Error {
	if e.Cause == nil {
		e.Cause = make(map[string]interface{}, 1)
	}

	errors, ok := e.Cause[Errors].([]TrParams)
	if !ok {
		e.Cause[Errors] = []TrParams{{
			TrKey:  trError.TrKey,
			Params: trError.Params,
		}}

		return e
	}

	e.Cause[Errors] = append(errors, trError)

	return e
}

func (e *Error) AddCause(args ...string) *Error {
	if e.Cause == nil {
		e.Cause = make(map[string]interface{})
	}

	for i := 0; i < len(args); i += 2 {
		strKey := args[i]
		e.Cause[strKey] = noneValue
		if i+1 < len(args) {
			e.Cause[strKey] = args[i+1]
		}
	}

	return e
}

func New(msg string) *Error {
	return &Error{ErrorText: msg}
}

type BadRequestError struct {
	StatusCode int
	Body       []byte
}

func (err *BadRequestError) Error() string {
	return fmt.Sprintf("status code %d, data '%s'", err.StatusCode, err.Body)
}
