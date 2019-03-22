package commonError

import (
	"errors"
	"fmt"
)

func NewStringErr(str string) error {
	return errors.New(str)
}
func NewStringErr2(args ...interface{}) error {
	return errors.New(fmt.Sprint(args...))
}
