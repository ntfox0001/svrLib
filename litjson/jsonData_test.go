package litjson_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/ntfox0001/svrLib/litjson"
	"github.com/ntfox0001/svrLib/util"
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

func Test3(t *testing.T) {
	jd := litjson.NewJsonDataFromJson(`{"int":-1}`)
	fmt.Println(jd.ToJson())

	fmt.Println(jd.Get("int").GetUInt32())
	fmt.Println(jd.Get("int").GetUInt64())
	fmt.Println(jd.Get("int").GetInt32())
	fmt.Println(jd.Get("int").GetInt64())
}

func Test4(t *testing.T) {
	jd := litjson.NewJsonDataFromJson(`{"bb":true}`)
	fmt.Println(jd.ToJson())

	fmt.Println(jd.Get("bb").GetBool())
}

func Test5(t *testing.T) {
	jd := util.JdMsg("LoginReq")
	jd.SetKey("keep", litjson.NewJsonData())
	jd.Get("keep").SetKey("gameId", 1)
	jd.SetKey("list", litjson.NewJsonData())
	jd.Get("list").Append("fff")
	jd.Get("list").Append(123)
	jd.Get("list").Append(true)
	fmt.Println(jd.ToJson())
}
func Test6(t *testing.T) {
	jd := litjson.NewJsonDataFromJson(`{"int":18446744073709551615}`)
	fmt.Println(jd.Get("int").GetFloat32())
	fmt.Println(jd.Get("int").GetFloat64())
	fmt.Println(jd.Get("int").GetInt32())
	fmt.Println(jd.Get("int").GetUInt32())
	fmt.Println(jd.Get("int").GetInt64())
	fmt.Println(jd.Get("int").GetUInt64())
	fmt.Println(jd.Get("int").GetUInt())
	fmt.Println(^uint64(0))
}

func Test7(t *testing.T) {
	jd := litjson.NewJsonDataFromJson(`{"int":5}`)
	njd, err := jd.Safe_Get("fff")
	fmt.Println(njd, err)

	njd1, err1 := jd.Safe_Get("int")
	fmt.Println(njd1.ToJson(), err1)
}

type TestInt struct {
	A int `json:"a"`
	B float32
	C string
}

func Test8(t *testing.T) {
	extra.RegisterFuzzyDecoders()
	jd := litjson.NewJsonDataFromJson(`{"a":"5","b":2,"c":"ff","d":true}`)
	bb := TestInt{}

	err := jd.Conv2Obj(&bb)

	fmt.Println(bb, err)
}
