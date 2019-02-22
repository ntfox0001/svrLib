package userSystem

import (
	"fmt"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/userSystem/userDefine"

	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ntfox0001/svrLib/log"

	jsoniter "github.com/json-iterator/go"
)

var (
	// 微信开放平台通过code获取
	wxOpenGetTokenByCodeUrl = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	wxOpenGetUserInfoUrl    = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
)

// 微信已经登陆，接受一个userdata
func (u *UserService) wxmpLoginProcess(w http.ResponseWriter, r *http.Request) {

	log.Debug("+ user http arrived.")
	s, _ := ioutil.ReadAll(r.Body)

	req := userDefine.WxMpLoginReq{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(s, &req); err != nil {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, err.Error())
		return
	}

	waitTokenChan := make(chan userDefine.GenerateTokenResp, 1)

	gtReq := userDefine.GenerateTokenReq{
		WaitTokenChan: waitTokenChan,
		UserData:      req.UserData,
	}

	u.userMgr.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("GenerateTokenReq", nil, gtReq))

	tokenResp := <-waitTokenChan

	if tokenResp.Token == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	loginResp := userDefine.WxMpLoginResp{
		MsgId:   "WxMpLoginResp",
		Token:   tokenResp.Token,
		UserId:  tokenResp.UserData.UserId,
		ErrorId: "0",
	}

	if s, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(loginResp); err != nil {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, err.Error())
	} else {

		io.WriteString(w, string(s))
	}

	log.Debug("- user http left.")
}

/*
微信通过code返回授权信息
{
"access_token":"ACCESS_TOKEN",
"expires_in":7200,
"refresh_token":"REFRESH_TOKEN",
"openid":"OPENID",
"scope":"SCOPE",
"unionid":"o6_bmasdasdsad6_2sgVt7hMZOPfL"
}
{"errcode":40029,"errmsg":"invalid code"}
*/
type wxmpLoginByCodeResp struct {
	access_token  string
	expires_in    int
	refresh_token string
	openid        string
	scope         string
	unionid       string
	errcode       int
	errmsg        string
}

/*
微信获取用户信息
{
"openid":"OPENID",
"nickname":"NICKNAME",
"sex":1,
"province":"PROVINCE",
"city":"CITY",
"country":"COUNTRY",
"headimgurl": "http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
"privilege":[
"PRIVILEGE1",
"PRIVILEGE2"
],
"unionid": " o6_bmasdasdsad6_2sgVt7hMZOPfL"
}
{
"errcode":40003,"errmsg":"invalid openid"
}
*/

type wxmpLoginGetUserInfoResp struct {
	openid     string
	nickname   string
	sex        int
	province   string
	city       string
	country    string
	headimgurl string
	privilege  []string
	unionid    string
	errcode    int
	errmsg     string
}

// 客户端获得code，服务器去微信验证
func (u *UserService) wxmpCodeLoginProcess(w http.ResponseWriter, r *http.Request) {

	log.Debug("+ user http arrived.")
	s, _ := ioutil.ReadAll(r.Body)

	req := userDefine.WxMpCodeLoginReq{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(s, &req); err != nil {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, err.Error())
		return
	}

	info, ok := u.appLoginInfos[req.AppId]
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "appId does not exist:"+req.AppId)
		log.Warn("wxmpCodeLoginProcess:appId does not exist", "appId", req.AppId)
		return
	}

	// 获得微信登陆授权
	gtbcUrl := fmt.Sprintf(wxOpenGetTokenByCodeUrl, info.AppId, info.Secret, req.Code)
	if gtbcwxrtstr, err := network.SyncHttpGet(gtbcUrl); err != nil {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, err.Error())
		return
	} else {
		gtbcresp := wxmpLoginByCodeResp{}
		if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(gtbcwxrtstr), &gtbcresp); err != nil {
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, err.Error())
			return
		}
		if gtbcresp.errcode != 0 {
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, err.Error())
			log.Warn("wxmpCodeLoginProcess: wxlogin err", gtbcresp.errcode, gtbcresp.errmsg)
			return
		}

		// 查询用户信息
		guiUrl := fmt.Sprintf(wxOpenGetUserInfoUrl, gtbcresp.access_token, gtbcresp.openid)

		if guiwxrtstr, err := network.SyncHttpGet(guiUrl); err != nil {
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, err.Error())
			return
		} else {
			guiresp := wxmpLoginGetUserInfoResp{}
			if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(guiwxrtstr), &guiresp); err != nil {
				w.WriteHeader(http.StatusForbidden)
				io.WriteString(w, err.Error())
				return
			}
			if guiresp.errcode != 0 {
				w.WriteHeader(http.StatusForbidden)
				io.WriteString(w, err.Error())
				log.Warn("wxmpCodeLoginProcess: wxlogin err", guiresp.errcode, guiresp.errmsg)
				return
			}
			usrData := userDefine.UserData{
				UserId:       0,
				UnionId:      guiresp.unionid,
				OpenId:       guiresp.openid,
				Nickname:     guiresp.nickname,
				Sex:          guiresp.sex,
				Language:     "",
				City:         guiresp.city,
				Province:     guiresp.province,
				Country:      guiresp.country,
				Headimgurl:   guiresp.headimgurl,
				AccessToken:  gtbcresp.access_token,
				RefreshToken: gtbcresp.refresh_token,
				ExpiresIn:    gtbcresp.expires_in,
			}
			waitTokenChan := make(chan userDefine.GenerateTokenResp, 1)

			gtReq := userDefine.GenerateTokenReq{
				WaitTokenChan: waitTokenChan,
				UserData:      usrData,
			}

			u.userMgr.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("GenerateTokenReq", nil, gtReq))

			tokenResp := <-waitTokenChan

			if tokenResp.Token == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			loginResp := userDefine.WxMpLoginResp{
				MsgId:   "WxMpLoginResp",
				Token:   tokenResp.Token,
				UserId:  tokenResp.UserData.UserId,
				ErrorId: "0",
			}

			if s, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(loginResp); err != nil {
				w.WriteHeader(http.StatusForbidden)
				io.WriteString(w, err.Error())
			} else {

				io.WriteString(w, string(s))
			}
		}

	}
	log.Debug("- user http left.")
}
