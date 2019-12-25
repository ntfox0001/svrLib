package litjson_test

import (
	"fmt"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/ntfox0001/svrLib/litjson"
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

type TestJsonData struct {
	A   int               `json:"a"`
	Baa *litjson.JsonData `json:"_baa_"`
}

func Test9(t *testing.T) {
	tjd := TestJsonData{
		A:   1,
		Baa: litjson.NewJsonData(),
	}
	tjd.Baa.SetKey("bb", "cc")
	newtjd := litjson.NewJsonDataFromObject(tjd)
	fmt.Println(newtjd.ToJson())

	tjd = TestJsonData{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(newtjd.ToJson(), &tjd); err != nil {
		fmt.Println(err.Error())
	} else {
		newtjd := litjson.NewJsonDataFromObject(tjd)
		fmt.Println(newtjd.ToJson())
	}

}

func Test10(t *testing.T) {
	i := 0
	f := func(a int) {
		switch a {
		case 0:
			i = 1
		case 1:
			i = 2
		}
		fmt.Println(i)
	}
	f(0)

}

func Test11(t *testing.T) {
	for i := 0; i < 100000; i++ {
		jd := litjson.NewJsonDataFromJson(`{"msgId":"UpdatePlayerDataReq","dataRoot":{"PlayerData":"{\"StaminaByInt\":32,\"dataVersion\":\"95558eba3d19ad4156db36c2f3674649c35e5f6e\",\"gold\":55,\"stamina\":32.0,\"adTicket\":0,\"hintFindout\":0,\"sceneDataMap\":{\"cat\":{\"sceneName\":\"cat\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"dinner\":{\"sceneName\":\"dinner\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"dog\":{\"sceneName\":\"dog\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"endline\":{\"sceneName\":\"endline\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"eve\":{\"sceneName\":\"eve\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"firework\":{\"sceneName\":\"firework\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"food\":{\"sceneName\":\"food\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"football\":{\"sceneName\":\"football\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"gift\":{\"sceneName\":\"gift\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":1553596944,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":3,\"enter\":0,\"success\":1,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"gift1\":{\"sceneName\":\"gift1\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"hatch\":{\"sceneName\":\"hatch\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"lighting\":{\"sceneName\":\"lighting\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"makeup\":{\"sceneName\":\"makeup\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"masterline\":{\"sceneName\":\"masterline\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"ninihome\":{\"sceneName\":\"ninihome\",\"keyState\":{\"New Game Object\":{\"beFound\":true,\"foundTime\":0,\"hint\":false},\"ball\":{\"beFound\":true,\"foundTime\":0,\"hint\":false},\"tape\":{\"beFound\":true,\"foundTime\":0,\"hint\":false}},\"unlockByGold\":false,\"lastFinishTime\":1552911465,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":5,\"enter\":0,\"success\":1,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"package\":{\"sceneName\":\"package\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"pyramidmaster\":{\"sceneName\":\"pyramidmaster\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"setup\":{\"sceneName\":\"setup\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"station\":{\"sceneName\":\"station\",\"keyState\":{\"station_clue1\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue10\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue11\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue12\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue13\":{\"beFound\":true,\"foundTime\":0,\"hint\":false},\"station_clue14\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue15\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue16\":{\"beFound\":true,\"foundTime\":0,\"hint\":false},\"station_clue2\":{\"beFound\":true,\"foundTime\":0,\"hint\":false},\"station_clue3\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue4\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue5\":{\"beFound\":true,\"foundTime\":0,\"hint\":false},\"station_clue6\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue7\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue8\":{\"beFound\":false,\"foundTime\":0,\"hint\":false},\"station_clue9\":{\"beFound\":false,\"foundTime\":0,\"hint\":false}},\"unlockByGold\":false,\"lastFinishTime\":1552911565,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":1,\"enter\":0,\"success\":1,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0},\"warriors\":{\"sceneName\":\"warriors\",\"keyState\":{},\"unlockByGold\":false,\"lastFinishTime\":0,\"canSkipNextCD\":false,\"masterPoint\":0,\"moonLevel\":0,\"enter\":0,\"success\":0,\"minTime\":0.0,\"averageTime\":0.0,\"totalTime\":0.0}},\"achievementPoint\":230,\"device\":\"Z270-HD3 (Gigabyte Technology Co., Ltd.)\",\"platform\":\"WindowsEditor\",\"timeZone\":28800,\"lastSceneId\":0,\"lastLoginTime\":0,\"lastLogoutTime\":0,\"watchADForGoldCount\":0,\"lastStaminaUpdate\":0,\"showReviewDialog\":true,\"infiniteAdTicket\":false,\"infiniteStamina\":false,\"doubleGold\":false,\"unlockAll\":false,\"lastAutoStaminaUpdate\":0,\"wxShareSessionTime\":0,\"wxShareTimelineTime\":0,\"signData\":0,\"lastSignDay\":26}","AchievementData":"{\"AchievementPlayRecord\":{\"1\":1,\"15\":1,\"2\":1,\"3\":1,\"4\":1,\"8\":1},\"OnActorFindout_Normal\":\"7\",\"OnOpeningWindow\":\"5\",\"OnOpeningWindow__byDay__\":\"3/26/2019\",\"OnSceneMissionCompleted_Special\":\"8\",\"OnSceneMissionCompleted_gift\":\"1\",\"OnSceneMissionCompleted_ninihome\":\"1\",\"OnSceneMissionCompleted_station\":\"1\",\"OnSwitchSceneTime\":\"178\",\"OnUseAdTicket\":\"1\",\"OnWatchAd\":\"40\"}"}}`)
		obj := jd.ToObject()
		jd2 := litjson.NewJsonDataFromObject(obj)
		s := jd2.ToJson()
		fmt.Println(s)
	}
}

type TestStruct struct {
	AA int
	BB string
}

func Test12(t *testing.T) {
	jd := litjson.NewJsonData()
	jd.SetKey("aaa", "bbb")

	a := &TestStruct{}
	ss := `{"AA":111, "BB":"fdsfds"}`

	println(a)
	litjson.Conv2Obj(ss, a)
	println(a)
	fmt.Println(a)
}
