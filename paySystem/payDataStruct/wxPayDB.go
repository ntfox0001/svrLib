package payDataStruct

import "github.com/ntfox0001/svrLib/database/dbtools/dbtoolsData"

const (
	PayStatusWaitForUserPay      = "WaitForUserPay"      // 等待用户支付
	PayStatusWaitForVerification = "WaitForVerification" // 等待服务器验证结果
	PayStatusSuccess             = "PaySuccess"          // 支付成功
)

// 支付数据库结构
type PayBillData struct {
	PayBillTable                      dbtoolsData.TableName
	UserId                            int                         `json:"userId,string" dbdef:"int"`
	BillId                            string                      `json:"billId" dbdef:"varchar(256),prim"`  // 商户订单号
	TransactionId                     string                      `json:"transactionId" dbdef:"text"`        // 第三方收据或订单号
	ProductId                         string                      `json:"productId" dbdef:"varchar(32)"`     // 商户产品id
	TotalFee                          int                         `json:"totalFee,string" dbdef:"int"`       // 订单价格
	Status                            string                      `json:"status" dbdef:"varchar(32),update"` // 订单状态
	CreateTime                        int64                       `json:"createTime,string" dbdef:"int"`     // 订单创建时间
	FinishTime                        int64                       `json:"finishTime,string" dbdef:"int"`     // 订单完成时间
	ExtentInfo                        string                      `json:"extentInfo" dbdef:"text"`           // 保存额外信息，比如第三方appid，receipt
	BillType                          string                      `json:"billType" dbdef:"varchar(32)"`      // 订单类型
	procedure_PayBillTable_PaySuccess dbtoolsData.CreateProcedure `dbsql:"create procedure PayBillTable_PaySuccess(inbillId varchar(256), intransactionId varchar(256),instatus tinyint) begin update PayBillTable set transactionId=intransactionId, status=instatus, finishTime=UNIX_TIMESTAMP() where billId=inbillId;end"`
	// load只加载过去3天的未完成订单，只有主动订单才需要加载 (WaitForVerification)
	procedure_PayBillTable_Load dbtoolsData.CreateProcedure `dbsql:"create procedure PayBillTable_Load() begin select * from PayBillTable where createTime>(UNIX_TIMESTAMP() - 60 * 60 * 24 * 3) and status = 'WaitForVerification';end"`
}
