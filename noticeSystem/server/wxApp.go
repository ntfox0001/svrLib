package main

import (
	"fmt"
	"oryxserver/logic/cmd/wxAccessRefreshServer/wxAccessRefMsg"
	"oryxserver/network"
	"oryxserver/timerSystem"
	"sync"
	"time"

	log "github.com/inconshreveable/log15"
	jsoniter "github.com/json-iterator/go"
)

const (
	RefreshAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	RefreshTicketUrl      = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

type WxApp struct {
	AppId                string
	Secret               string
	AccessToken          string
	AccessTokenExpiresIn int64
	Ticket               string
	TicketExpiresIn      int64
	AccessTokenLock      sync.RWMutex
	TicketLock           sync.RWMutex
}

func (w *WxApp) Initial() {
	w.RefreshAccessToken(nil, nil)
	t := timerSystem.NewTimerItemLoopByFunc(Config.RefreshTime, w.RefreshAccessToken, nil)
	timerSystem.Instance().AddTimer(t)
}

func (w *WxApp) RefreshAccessToken(data interface{}, time *time.Time) {
	url := fmt.Sprintf(RefreshAccessTokenUrl, w.AppId, w.Secret)
	//fmt.Println(url)
	if result, err := network.SyncHttpGet(url); err == nil {

		log.Info(result)
		var resp wxAccessRefMsg.WxMpRefreshAccessTokenResp
		if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(result), &resp); err == nil {

			log.Debug("WxMp access refresh", "appid", w.AppId, "old", w.AccessToken, "new", resp.AccessToken)
			w.AccessTokenLock.Lock()
			w.AccessToken = resp.AccessToken
			w.AccessTokenExpiresIn = int64(resp.ExpiresIn)
			w.AccessTokenLock.Unlock()

			w.RefreshTicket(w.AccessToken)
		} else {
			log.Error("Failed to RefreshAccessToken", "result", result, "err", err.Error())
		}
	}
}

func (w *WxApp) RefreshTicket(access_token string) {
	url := fmt.Sprintf(RefreshTicketUrl, access_token)
	//fmt.Println(url)
	if result, err := network.SyncHttpGet(url); err == nil {

		log.Info(result)
		var resp wxAccessRefMsg.WxMpRefreshTickedResp
		if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(result), &resp); err == nil {

			log.Debug("WxMp ticket refresh", "appid", w.AppId, "old", w.Ticket, "new", resp.Ticket)
			w.TicketLock.Lock()
			w.Ticket = resp.Ticket
			w.TicketExpiresIn = int64(resp.ExpiresIn)
			w.TicketLock.Unlock()

		} else {
			log.Error("Failed to RefreshAccessToken", "result", result, "err", err.Error())
		}
	}
}

func (w *WxApp) GetAccessToken() (rt string) {
	w.AccessTokenLock.RLock()
	rt = w.AccessToken
	w.AccessTokenLock.RUnlock()
	return
}

func (w *WxApp) GetTicket() (rt string) {
	w.TicketLock.RLock()
	rt = w.Ticket
	w.TicketLock.RUnlock()
	return
}
