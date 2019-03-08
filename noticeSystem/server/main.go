package main

// 运行一个access中控
import (
	"flag"
	"io"
	"io/ioutil"
	"net/http"

	"os"
	"os/signal"
	"strings"

	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/noticeSystem/wxAccessRefreshServer/wxAccessRefMsg"
	"github.com/ntfox0001/svrLib/selectCase"
	"github.com/ntfox0001/svrLib/timerSystem"

	jsoniter "github.com/json-iterator/go"
	"github.com/ntfox0001/svrLib/log"
)

var (
	configFile string
	selectLoop *selectCase.SelectLoop

	server *network.Server

	wxApps []*WxApp

	c chan os.Signal
)

func init() {
	flag.StringVar(&configFile, "config", "config.json", "config filename")

}

func main() {
	flag.Parse()

	h1 := log.CallerFileHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stdout, log15Ex.LogfmtFormat())))

	h2 := log.Must.FileHandler("watrSvr.log", log15Ex.LogfmtFormat())
	log.Root().SetHandler(log.MultiHandler(h1, h2))

	if err := InitApplicationConfig(configFile); err != nil {
		log.Error("failed to config file")
		return
	}
	timerSystem.Instance().Initial()

	for _, v := range Config.WxMp {
		va := WxApp{
			AppId:                v.AppId,
			Secret:               v.Secret,
			AccessToken:          "",
			AccessTokenExpiresIn: 0,
			Ticket:               "",
			TicketExpiresIn:      0,
		}
		va.Initial()
		wxApps = append(wxApps, &va)
	}

	server = network.NewServer(Config.ListenIp, Config.Port)

	// 注册对外路由
	server.RegisterRouter("/GetWxAccessToken", network.RouterHandler{ProcessHttpFunc: getWxAccessToken})
	server.RegisterRouter("/GetWxTicket", network.RouterHandler{ProcessHttpFunc: getWxTicket})

	server.Start()

	// signal
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

runable:
	for {
		select {
		case <-c:
			break runable
		}
	}

	server.Close()
	timerSystem.Instance().Release()
}

func getWxAccessToken(w http.ResponseWriter, r *http.Request) {
	// 只有白名单里的ip才能访问
	nofind := true
	for _, v := range Config.WhiteIp {
		if strings.Contains(r.RemoteAddr, v) {
			nofind = false
			break
		}
	}

	if nofind == true {
		return
	}
	s, _ := ioutil.ReadAll(r.Body)
	req := wxAccessRefMsg.WxAccessTokenReq{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(s, &req); err != nil {
		log.Error("getWxAccessToken", "error", err.Error())
		io.WriteString(w, err.Error())
		return
	}

	for _, v := range wxApps {
		if v.AppId == req.AppId {
			resp := wxAccessRefMsg.WxAccessTokeyResp{Token: v.GetAccessToken()}

			rt, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resp)

			io.WriteString(w, string(rt))
		}
	}
}

func getWxTicket(w http.ResponseWriter, r *http.Request) {
	// 只有白名单里的ip才能访问
	nofind := true
	for _, v := range Config.WhiteIp {
		if strings.Contains(r.RemoteAddr, v) {
			nofind = false
			break
		}
	}

	if nofind == true {
		return
	}
	s, _ := ioutil.ReadAll(r.Body)
	req := wxAccessRefMsg.WxTicketReq{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(s, &req); err != nil {
		log.Error("getWxTicket", "error", err.Error())
		io.WriteString(w, err.Error())
		return
	}

	for _, v := range wxApps {
		if v.AppId == req.AppId {
			resp := wxAccessRefMsg.WxTicketResp{Ticket: v.GetTicket()}

			rt, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resp)

			io.WriteString(w, string(rt))
		}
	}
}
