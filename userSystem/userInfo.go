package userSystem

import (
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"github.com/ntfox0001/svrLib/userSystem/userDefine"
)

type UserInfo struct {
	usrData *userDefine.UserData
}

// 异步加载 for userManagerManager
func asyncNewUserInfo(cb *selectCaseInterface.CallbackHandler, usrData *userDefine.UserData) *UserInfo {

	op := database.Instance().NewOperation("call userInsert(?,?,?,?,?)",
		usrData.UnionId, usrData.OpenId, usrData.Nickname, usrData.Headimgurl, usrData.RefreshToken)
	op.UserData = usrData
	database.Instance().ExecOperationForCB(cb, op)

	return newUserInfoForUserData(usrData)
}

// 创建新用户
func newUserInfoForUserData(ud *userDefine.UserData) *UserInfo {
	ui := &UserInfo{
		usrData: ud,
	}

	return ui
}

func (u *UserInfo) GetUserData() *userDefine.UserData {
	return u.usrData
}
