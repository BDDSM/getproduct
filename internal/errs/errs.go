package errs

import "encoding/json"

type Error struct {
	error string
}

func New(err error) *Error {
	return &Error{error: err.Error()}
}

func (e *Error) Error() string {
	return e.error
}

func (e *Error) MarshalJSON() ([]byte, error) {

	errormap := make(map[string]interface{})
	errormap["error"] = e.error

	return json.Marshal(errormap)

}
