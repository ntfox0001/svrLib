package litjson_test

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.funplus.io/jiang.liu/gameSvr/litjson"
)

func Test1(t *testing.T) {
	a := make([]int, 10, 10)

	var b interface{}
	b = a

	switch b.(type) {
	case []interface{}:
		fmt.Print("aaa")
	}

	fmt.Print(reflect.TypeOf(a))
}
func Test2(t *testing.T) {
	jd := litjson.NewJsonDataFromJson(`{"aaa":"ffff", "f32":223.22, "f64":3.2414212222, "int32":123213213, "uint64":327326342543254321}`)
	jd.SetKey("aaa", "eee")
	jdarray := litjson.NewJsonData()
	jd.SetKey("obej", jdarray)
	fmt.Println(jd.ToJson())
	jdarray.Append(123)
	fmt.Println(jd.ToJson())
	jdarray.Append(321)
	fmt.Println(jd.ToJson())
	jdarray.Append(4444)
	jdarray.Append(66666)
	jdarray.SetIndex(4, 999)
	fmt.Println(jd.ToJson())
	jdarray.RemoveId(1)
	fmt.Println(jd.ToJson())
	jd.RemoveKey("f32")
	fmt.Println(jd.ToJson())

}
