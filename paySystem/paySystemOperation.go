package paySystem

import (
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/paySystem/payDataStruct"
	"github.com/ntfox0001/svrLib/util"
	"time"
)

/*
	支付数据库操作
	插入预支付id记录
	插入支付记录
	更新支付记录状态(4个状态)

	这里的函数都是同步函数
*/

// 创建新订单，微信等是被动验证，苹果是主动验证
// extentInfo 可以保存appid或者receipt
func (*PaySystem) PayRecord_NewBill(userId int, billId string, productId string, fee int, billType string, extentInfo string, passive bool) error {
	status := payDataStruct.PayStatusWaitForUserPay
	if !passive {
		status = payDataStruct.PayStatusWaitForVerification
	}
	// inuserId int,inbillId varchar(256),intransactionId text,inproductId varchar(32),intotalFee int,instatus varchar(32),increateTime int,infinishTime int,inextentInfo text,inbillType varchar(32)
	op := database.Instance().NewOperation("call PayBillTable_Insert(?,?,?,?,?,?,?,?,?,?)",
		userId, billId, "", productId, fee, status, time.Now().Unix(), 0, extentInfo, billType)

	_, err := database.Instance().SyncExecOperation(op)
	return err
}

// 设置订单支付成功
func (*PaySystem) PayRecord_SetPayStatusSuccess(billId string, transactionId string) error {
	//inbillId varchar(256), intransactionId varchar(256),instatus tinyint
	op := database.Instance().NewOperation("call PayBillTable_PaySuccess(?,?,?)", billId, transactionId, payDataStruct.PayStatusSuccess)

	_, err := database.Instance().SyncExecOperation(op)
	return err
}

// 设置一个订单状态为错误，错误信息不能超过32个字符
func (*PaySystem) PayRecord_SetError(billId string, errInfo string) error {
	op := database.Instance().NewOperation("call PayBillTable_UpdateStatusByBillId(?,?)", billId, errInfo)
	_, err := database.Instance().SyncExecOperation(op)
	return err
}



// 查询订单
func (*PaySystem) PayRecord_Query(billId string) (payDataStruct.PayBillData, error) {
	op := database.Instance().NewOperation("call PayBillTable_QueryByBillId(?)", billId)
	rt, err := database.Instance().SyncExecOperation(op)
	if err != nil {
		return payDataStruct.PayBillData{}, err
	}
	payBillDS := rt.FirstSet()
	if len(payBillDS) != 1 {
		return payDataStruct.PayBillData{}, commonError.NewStringErr2("bill dataset length must be 1. len:", len(payBillDS))
	}

	pd := payDataStruct.PayBillData{}
	if err := util.I2Stru(payBillDS[0], &pd); err != nil {
		return payDataStruct.PayBillData{}, commonError.NewStringErr2("bill data format error")
	}

	return pd, nil
}
