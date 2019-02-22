package userSystem

import (
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
)

type IUserCallback interface {
	OnInitial(helper selectCaseInterface.ISelectLoopHelper) error
	OnRelease()
}
type IUserManagerCallback interface {
	OnInitial(helper selectCaseInterface.ISelectLoopHelper) error
	OnInitUser(usr *User) // 在server启动时，创建user
	OnNewUser(usr *User)  // 在运行时新user创建
	OnRelease()
}
type IServiceCallback interface {
	//
	OnInitial(server *network.Server) error
	OnRelease()
}
