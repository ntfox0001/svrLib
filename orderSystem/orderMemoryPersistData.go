package orderSystem

import (
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/orderSystem/orderData"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"time"
)

type OrderMemoryPersistData struct {
	pdCallback       orderData.IPersistDataCallback
	orderDataMap     map[string]orderData.IOrderData
	oldOrderDataMap  map[string]orderData.IOrderData // 由于map的delete并不真的删除item的key，所以map会越来越大，所以每次插入的时候检查是否放弃老map
	preTime          int64
	discardTime      int64
	selectLoopHelper selectCaseInterface.ISelectLoopHelper
}

func NewOrderMemoryPersistData(discardTime int64) orderData.IPersistData {
	return &OrderMemoryPersistData{
		orderDataMap:    make(map[string]orderData.IOrderData),
		oldOrderDataMap: make(map[string]orderData.IOrderData),
		preTime:         time.Now().Unix(),
		discardTime:     discardTime,
	}
}

func (c *OrderMemoryPersistData) Query(key string) {
	var item orderData.IOrderData
	if v, ok := c.orderDataMap[key]; !ok {
		if v, ok := c.oldOrderDataMap[key]; !ok {
			c.pdCallback.OnQuery(key, nil, nil)
			return
		} else {
			item = v
		}
	} else {
		item = v
	}
	c.pdCallback.OnQuery(key, item, nil)
}

func (c *OrderMemoryPersistData) Insert(data interface{}) {
	item := data.(orderData.IOrderData)

	if _, ok := c.orderDataMap[item.GetCustomId()]; ok {
		c.pdCallback.OnInsert(item.GetCustomId(), commonError.NewCommErr("item already exist.", orderData.Error_ItemAlreadyExist))
		return
	}
	if _, ok := c.oldOrderDataMap[item.GetCustomId()]; !ok {
		c.pdCallback.OnInsert(item.GetCustomId(), commonError.NewCommErr("item already exist.", orderData.Error_ItemAlreadyExist))
		return
	}
	// 达到丢弃时间，那么舍弃之前的map
	if time.Now().Unix()-c.preTime > c.discardTime {
		c.oldOrderDataMap = c.orderDataMap
		c.orderDataMap = make(map[string]orderData.IOrderData)
	}
	c.orderDataMap[item.GetCustomId()] = item

	c.pdCallback.OnInsert(item.GetCustomId(), nil)
}

func (c *OrderMemoryPersistData) Initial(slHelper selectCaseInterface.ISelectLoopHelper, pdCallback orderData.IPersistDataCallback) {
	c.pdCallback = pdCallback
	c.selectLoopHelper = slHelper

	c.pdCallback.OnInitial(c.orderDataMap, nil)
}

func (c *OrderMemoryPersistData) Update(key string, status int) {
	if v, ok := c.orderDataMap[key]; !ok {
		if v, ok := c.oldOrderDataMap[key]; !ok {
			c.pdCallback.OnUpdate(key, commonError.NewCommErr("not exist item.", orderData.Error_ItemNotExist))
			return
		} else {
			c.oldOrderDataMap[key].SetStatus(status)

			c.pdCallback.OnUpdate(v.GetCustomId(), nil)
			return
		}
	} else {
		c.orderDataMap[key].SetStatus(status)

		c.pdCallback.OnUpdate(v.GetCustomId(), nil)
		return
	}
}
