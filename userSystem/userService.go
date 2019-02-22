package userSystem

import (
	"github.com/ntfox0001/svrLib/network"

	"github.com/ntfox0001/svrLib/log"
)

type UserService struct {
	userMgr           *UserManager
	listenip          string
	port              string
	server            *network.Server
	ssl               bool
	certFile, keyFile string
	callback          IServiceCallback
	appLoginInfos     map[string]AppLoginInfo
}
type AppLoginInfo struct {
	AppId  string
	Secret string
}
type UserServiceParams struct {
	Listenip string `json:"listenIp"`
	Port     string `json:"port"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
	AppId    string `json:"appId"`
	Secret   string
	Servcb   IServiceCallback     `json:"-"`
	UsrMgrcb IUserManagerCallback `json:"-"`
	Usrcb    IUserCallback        `json:"-"`
}

func NewUserService(params UserServiceParams) *UserService {

	usrServ := UserService{
		userMgr:       NewUserManager(params.Listenip, params.Port, params.UsrMgrcb, params.Usrcb),
		listenip:      params.Listenip,
		port:          params.Port,
		server:        nil,
		ssl:           true,
		certFile:      params.CertFile,
		keyFile:       params.KeyFile,
		callback:      params.Servcb,
		appLoginInfos: make(map[string]AppLoginInfo),
	}

	return &usrServ
}

func (u *UserService) Initial() {
	if u.ssl {
		u.server = network.NewServerSsl(u.listenip, u.port, u.certFile, u.keyFile)
	} else {
		u.server = network.NewServer(u.listenip, u.port)
	}

	// 用户长连接
	wsr := network.NewRouterWSHandler(u.userMgr)
	wsr.DisableCheckOrigin(false)
	u.server.RegisterRouter("/user", wsr)

	// 注册微信登陆，外面有php已经登陆时使用，多用于网页服务，cookie保持登录状态
	u.server.RegisterRouter("/wxmpLogin", network.RouterHandler{ProcessHttpFunc: u.wxmpLoginProcess})
	// 注册微信code方式登陆,用于长连接用户登陆
	u.server.RegisterRouter("/wxmpCodeLogin", network.RouterHandler{ProcessHttpFunc: u.wxmpCodeLoginProcess})

	// 初始化回调
	if err := u.callback.OnInitial(u.server); err != nil {
		log.Error("callback initial error", "err", err.Error())
		return
	}

	if err := u.server.Start(); err != nil {
		log.Error("http server error", "err", err.Error())
		return
	}
}

func (u *UserService) Release() {
	u.callback.OnRelease()
	u.server.Close()
	u.userMgr.Release()

	log.Debug("UserService release.")
}

func (u *UserService) AddAppLoginInfo(info AppLoginInfo) {
	u.appLoginInfos[info.AppId] = info
}
