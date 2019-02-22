package orderData

import "github.com/ntfox0001/svrLib/database/dbtools/dbtoolsData"

type OrderServerData struct {
	OrderServerDataTable dbtoolsData.TableName
	CustomId             string `json:"customId" dbdef:"varchar(128),prim"`
	GroupId              string `json:"groupId" dbdef:"varchar(20),prim" dbcomment:"分组名称"`
	SequenceId           uint64 `json:"SequenceId" dbdef:"int unsigned" dbcomment:"发送消息的序列id"`
	Origin               string `json:"origin" dbdef:"varchar(20)" dbcomment:"订单来源"`
	Target               string `json:"target" dbdef:"varchar(20)" dbcomment:"订单目标"`
	CreateTime           string `json:"createTime" dbdef:"datetime"`
	Status               int    `json:"status,string" dbdef:"tinyint" dbcomment:"0：等待确认，-1确认发送,>0表示重试次数"`
	Data                 string `json:"data" dbdef:"TEXT"`
	SendTime             int64  `json:"-"` // 上一次尝试发送的时间
}

func (d OrderServerData) GetCustomId() string {
	return d.CustomId
}
func (d OrderServerData) GetStatus() int {
	return d.Status
}
func (d OrderServerData) SetStatus(s int) {
	d.Status = s
}
func (d OrderServerData) GetInsertParams() []interface{} {
	return []interface{}{}
}

type OrderClientData struct {
	OrderClientDataTable dbtoolsData.TableName
	CustomId             string `json:"customId" dbdef:"varchar(128),prim"`
	Name                 string `json:"name" dbdef:"varchar(20),query" dbcomment:"客户端名字应全局唯一"`
	SequenceId           uint64 `json:"SequenceId" dbdef:"int unsigned" dbcomment:"发送消息的序列id"`
	Target               string `json:"target" dbdef:"varchar(20)" dbcomment:"订单目标"`
	CreateTime           string `json:"createTime" dbdef:"datetime"`
	Status               int    `json:"status,string" dbdef:"tinyint,update" dbcomment:"0：等待确认，-1确认发送,>0表示重试次数"`
	Data                 string `json:"data" dbdef:"TEXT"`
	SendTime             int64  `json:"-"` // 上一次尝试发送的时间
	// 定制存储过程
	// load最近三天的未完成订单
	procedure1 dbtoolsData.CreateProcedure `dbsql:"create procedure OrderClientData_Load (inname varchar(20), inday int) begin select * from OrderClientData where name=inname and createtime < DATE_FORMAT(date_add(now(), interval -inday day), '%Y-%m-%d');end"`
}

func (d OrderClientData) GetCustomId() string {
	return d.CustomId
}
func (d OrderClientData) GetStatus() int {
	return d.Status
}
func (d OrderClientData) SetStatus(s int) {
	d.Status = s
}
func (d OrderClientData) GetInsertParams() []interface{} {
	return []interface{}{d.CustomId, d.Name, d.SequenceId, d.Target, d.CreateTime, d.Status, d.Data}
}
