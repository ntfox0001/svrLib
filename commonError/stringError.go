package commonError

import (
	"fmt"
)

type StringError struct {
	errStr string
}

func (ce StringError) Error() string {
	return ce.errStr
}

func NewStringErr(str string) StringError {
	return StringError{errStr: str}
}
func NewStringErr2(args ...interface{}) StringError {
	return StringError{errStr: fmt.Sprint(args...)}
}
