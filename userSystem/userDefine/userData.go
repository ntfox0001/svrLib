package userDefine

// user数据库结构
type UserData struct {
	UserId       int    `json:"userId,string"`
	UnionId      string `json:"unionid"`
	OpenId       string `json:"openid"`
	Nickname     string `json:"nickname"`
	Sex          int    `json:"sex"`
	Language     string `json:"language"`
	City         string `json:"city"`
	Province     string `json:"province"`
	Country      string `json:"country"`
	Headimgurl   string `json:"headimgurl"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
