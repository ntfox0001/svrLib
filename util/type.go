package util

import (
	"reflect"
)

func GetTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	return t.Name()
}
