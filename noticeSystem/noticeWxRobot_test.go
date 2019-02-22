package noticeSystem

import (
	"testing"
)

func TestNoticeWxRobot(t *testing.T) {
	templ := make([]WxRobotTemplateCfg, 0)
	templ = append(templ, WxRobotTemplateCfg{
		Type: "ttt",
		Msg:  "我擦我擦我擦啊{ss}",
	})
	robot, _ := newWxRobot("192.168.1.117", "8888", 60, templ)

	data := make(map[string]string)
	data["{type}"] = "ttt"
	data["{ss}"] = "[我是参数]"
	robot.roomSend(data)
	data2 := make(map[string]string)
	data2["{type}"] = "ttt"
	data2["{ss}"] = "[憨笑]"
	robot.roomSend(data2)
	robot.close()
}
