package orderSystem

import (
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/orderSystem/orderData"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"github.com/ntfox0001/svrLib/util"

	"github.com/ntfox0001/svrLib/log"
)

type OrderDBPersistData struct {
	sqlProvider      orderData.IPersistDataSql
	dbSys            *database.DatabaseSystem
	selectLoopHelper selectCaseInterface.ISelectLoopHelper
	pdCallback       orderData.IPersistDataCallback
	discardDay       uint
}

// 创建数据库持久化对象
func NewOrderDBPersistData(dbCfg database.DbConfig, discardDay uint, sqlProvider orderData.IPersistDataSql) (orderData.IPersistData, error) {
	dbSys := &database.DatabaseSystem{}
	err := dbSys.Initial(dbCfg.Ip, dbCfg.Port, dbCfg.User, dbCfg.Passwd, dbCfg.DbName, 50, 50)
	if err != nil {
		return nil, err
	}

	pd := &OrderDBPersistData{
		dbSys:       dbSys,
		discardDay:  discardDay,
		sqlProvider: sqlProvider,
	}

	return pd, nil
}

func (c *OrderDBPersistData) Query(key string) {
	op := c.dbSys.NewOperation(c.sqlProvider.GetQuerySql(), key)
	cb := c.selectLoopHelper.NewCallbackHandler("OrderData_Query_Resp", key)

	c.dbSys.ExecOperationForCB(cb, op)
}

func (c *OrderDBPersistData) orderClientDataQueryByCustomIdResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	dbrt := msg.Content.(*database.DataResult)
	customId := msg.UserData.(string)

	if dbrt.Err != nil {
		// 数据库错误
		c.pdCallback.OnQuery(customId, nil, commonError.NewCommErr(dbrt.Err.Error(), orderData.Error_PersistError))
	} else {
		dbset := dbrt.FirstSet()
		if len(dbset) == 0 {
			// item 不存在
			c.pdCallback.OnQuery(customId, nil, nil)
		} else {
			d := orderData.OrderClientData{}
			if err := util.I2Stru(dbset[0], &d); err != nil {
				c.pdCallback.OnQuery(customId, nil, err)
			} else {
				c.pdCallback.OnQuery(customId, d, nil)
			}
		}
	}
	return true
}

func (c *OrderDBPersistData) Insert(data interface{}) {
	// 优化空间
	// 大量相同数据库调用，可以使用一个缓存的op
	item := data.(orderData.IOrderData)
	op := c.dbSys.NewOperation(c.sqlProvider.GetInsertSql(), item.GetInsertParams()...)

	// 发送数据库
	cb := c.selectLoopHelper.NewCallbackHandler("OrderData_Insert_Resp", item.GetCustomId())
	c.dbSys.ExecOperationForCB(cb, op)
}

func (c *OrderDBPersistData) orderClientDataInsertResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	dbrt := msg.Content.(*database.DataResult)
	customId := msg.UserData.(string)

	// 返回消息中userdata用于传送err信息
	c.pdCallback.OnInsert(customId, dbrt.Err)
	return true
}

func (c *OrderDBPersistData) Initial(slHelper selectCaseInterface.ISelectLoopHelper, pdCallback orderData.IPersistDataCallback) {
	c.selectLoopHelper = slHelper
	c.pdCallback = pdCallback
	slHelper.RegisterEvent("OrderData_Insert_Resp", c.orderClientDataInsertResp)
	slHelper.RegisterEvent("OrderData_Query_Resp", c.orderClientDataQueryByCustomIdResp)
	slHelper.RegisterEvent("OrderData_Update_Resp", c.orderClientDataUpdateStatusByCustomIdResp)
	slHelper.RegisterEvent("OrderData_Initial_Resp", c.orderClientDataLoadResp)

	op := c.dbSys.NewOperation(c.sqlProvider.GetInitialSql(), c.pdCallback.GetName(), c.discardDay)
	cb := c.selectLoopHelper.NewCallbackHandler("OrderData_Initial_Resp", nil)
	c.dbSys.ExecOperationForCB(cb, op)
}

func (c *OrderDBPersistData) orderClientDataLoadResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	dbrt := msg.Content.(*database.DataResult)

	if dbrt.Err != nil {
		c.pdCallback.OnInitial(nil, dbrt.Err)
		return true
	}
	dset := dbrt.FirstSet()
	if dset != nil {
		rdmap := make(map[string]orderData.OrderClientData)
		for _, v := range dset {
			var d orderData.OrderClientData
			if err := util.I2Stru(v, &d); err != nil {
				log.Error("orderClientDataLoadResp error", "err", err.Error())
				c.pdCallback.OnInitial(nil, err)
				return true
			}
			rdmap[d.CustomId] = d
		}
		c.pdCallback.OnInitial(rdmap, nil)
	} else {
		c.pdCallback.OnInitial(nil, commonError.NewStringErr("orderClientDataLoadResp error:no data set."))
	}
	return true
}

func (c *OrderDBPersistData) orderClientDataUpdateStatusByCustomIdResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	dbrt := msg.Content.(*database.DataResult)
	customId := msg.UserData.(string)

	c.pdCallback.OnUpdate(customId, dbrt.Err)
	return true
}
func (c *OrderDBPersistData) Update(key string, status int) {
	op := c.dbSys.NewOperation(c.sqlProvider.GetUpdateSql(), key, status)
	cb := c.selectLoopHelper.NewCallbackHandler("OrderData_Update_Resp", key)
	c.dbSys.ExecOperationForCB(cb, op)
}
