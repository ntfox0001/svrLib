package paySystem

import (
	"fmt"
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/log"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/paySystem/payDataStruct"
	"github.com/ntfox0001/svrLib/util"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	sandUrl  = "https://sandbox.itunes.apple.com/verifyReceipt"
	buyUrl   = "https://buy.itunes.apple.com/verifyReceipt"
	netError = 1
)

var (
	repeatTimes = []int64{0, 5, 10, 15, 60, 300, 3600} // 重试次数对应的时间
)

type applePayItem struct {
	billId        string // 订单id
	userId        int
	receipt       string
	productId     string
	repeatCount   int // 重试次数
	updateManager *util.UpdateManager
	createTime    int64 // 用于计算更新时间
	closeCh       chan interface{}
}

func newApplePayItem(userId int, receipt string, productId string) *applePayItem {
	item := &applePayItem{
		billId:      util.GetUniqueId(),
		userId:      userId,
		receipt:     receipt,
		productId:   productId,
		repeatCount: 0,

		createTime: time.Now().Unix(),
	}

	item.updateManager = util.NewUpdateManager2("applePay", time.Second, item._update)

	return item
}

// 开始
func (i *applePayItem) run() {
	i.updateManager.Run(time.Second)
}

// 等待这个支付过程退出
func (i *applePayItem) waitForClose() {
	<-i.closeCh
}
func (i *applePayItem) _getRepeatTime() int64 {
	rtLen := len(repeatTimes)
	count := i.repeatCount
	if i.repeatCount >= rtLen {
		count = rtLen - 1
	}
	return repeatTimes[count]
}

func (i *applePayItem) _update() {
	if time.Now().Unix()-i.createTime >= i._getRepeatTime() {
		i.repeatCount++

		i.beginPay()
	}
}

func (i *applePayItem) beginPay() {
	// 保存订单数据
	if err := _self.PayRecord_NewBill(i.userId, i.billId, i.productId, 0, "APPLE", i.receipt, false); err != nil {
		log.Warn("apple NewBill failed.", "receipt", i.receipt, "userId", i.userId)
		return
	}

	// 获得结果
	if resp, err := i.getReceiptResp(i.receipt); err != nil {
		log.Warn(err.Error(), "userId", i.userId)
		if err.(commonError.CommError).GetType() != netError {
			// 无法挽回的错误，关闭
			i.close()
		}
		return
	} else {
		// 检查
		if err := i.validatingReceipt(i.userId, i.billId, i.productId, resp); err == nil {
			if _self.callback == nil {
				log.Error("WxCallback is nil")
				return
			}

			notify := payDataStruct.PaySystemNotify{
				ExtentData: nil,
				ProductId:  i.productId,
				UserId:     i.userId,
				PayType:    payDataStruct.PaySystemNotify_PayType_Apple,
			}

			// 发送回调
			_self.callback.SendReturnMsgNoReturn(notify)
		}
		// 这里的错误都是无法挽回的，所以不管是否成功，都退出
		i.close()
	}
}

// 关闭支付过程
func (i *applePayItem) close() {
	i.updateManager.Close()
	i.closeCh <- struct{}{}
}

// 验证收据是否有效
func (i *applePayItem) getReceiptResp(receipt string) (payDataStruct.IapPayDataResp, error) {
	req := payDataStruct.IapPayDataReq{
		Receipt_data: receipt,
	}

	reqStr, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)
	if err != nil {
		log.Warn("req marshal failed.", "req", reqStr)
		return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("req marshal failed.", 0)
	}

	sandRespStr, err := network.SyncHttpPost(sandUrl, reqStr, network.ContentTypeJson)
	if err != nil {
		log.Warn("sandbox http post failed", "err", err.Error())
		return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("sandbox http post failed", netError)
	}

	resp := payDataStruct.IapPayDataResp{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(sandRespStr, &resp); err != nil {
		log.Warn("the format of sandbox's resp failed.")
		return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("the format of sandbox's resp failed.", 0)
	}

	// check status
	if resp.Status == 21007 {
		// 这是一个正式服的收据
		buyRespStr, err := network.SyncHttpPost(buyUrl, reqStr, network.ContentTypeJson)
		if err != nil {
			log.Warn("http post failed", "err", err.Error())
			return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("http post failed.", netError)
		}

		resp = payDataStruct.IapPayDataResp{}
		if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(buyRespStr, &resp); err != nil {
			log.Warn("the format of resp failed.")
			return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("the format of resp failed.", 0)
		}
	}

	return resp, nil
}

// 验证收据是否正确
func (i *applePayItem) validatingReceipt(userId int, billId string, productId string, resp payDataStruct.IapPayDataResp) error {
	if resp.Status == 0 {
		transaction_id := ""
		for _, v := range resp.Receipt.In_app {
			if productId == v.Product_id {
				transaction_id = v.Transaction_id
				break
			}
		}
		if transaction_id == "" {
			log.Error("transaction does not found.")
			return commonError.NewStringErr("transaction does not found.")
		}
		if err := _self.PayRecord_SetPayStatusSuccess(billId, transaction_id); err != nil {
			log.Error("apple SetPayStatusSuccess failed", "err", err.Error())
			return commonError.NewStringErr("apple SetPayStatusSuccess failed")
		}
	} else {
		if err := _self.PayRecord_SetError(billId, fmt.Sprint(resp.Status)); err != nil {
			log.Error("apple SetError failed", "err", err.Error())
			return commonError.NewStringErr("apple SetError failed")
		}

	}
	return nil
}
