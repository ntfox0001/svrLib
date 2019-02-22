package payDataStruct

// 支付系统回调函数的参数
type PaySystemNotify struct {
	ExtentData interface{}
	UserId     int    `json:"userId"`
	ProductId  string `json:"productId"` //商户产品id
	PayType    int    // 1：wx,2： apple
}

const (
	PaySystemNotify_PayType_Wx    = 1
	PaySystemNotify_PayType_Apple = 2
)
