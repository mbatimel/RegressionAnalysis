package errors

var (
	AccessDeniedError     = func() *Error { return New("access denied") }
	ForbiddenError        = func() *Error { return New("forbidden") }
	MethodNotAllowedError = func() *Error { return New("method not allowed") }
	InternalServerError   = func() *Error { return New("internal server error") }
	NotFound              = func() *Error { return New("not found") }
	ObjectInUse           = func() *Error { return New("object in use") }
	AlreadyExists         = func() *Error { return New("key already exists") }
	BadMetaValue          = func() *Error { return New("bad meta") }
	InvalidRequest        = func() *Error { return New("invalid request") }
	NotEnoughMoney        = func() *Error { return New("not enough money") }
	NotInCache            = func() *Error { return New("value not found in cache") }
)

const (
	Errors = "errors"
)

type TrParams struct {
	TrKey  string                 `json:"trKey"`
	Params map[string]interface{} `json:"params,omitempty"`
}
