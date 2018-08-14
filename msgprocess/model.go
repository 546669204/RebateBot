package main

type MemberItemModel struct {
	UserName   string `json:"UserName"`
	NickName   string `json:"NickName"`
	RemarkName string `json:"RemarkName"`
	Province   string `json:"Province"`
	City       string `json:"City"`
	Sex        int    `json:"Sex"`
	Signature  string `json:"Signature"`
}

type MemberListModel struct {
	MemberCount int               `json:"MemberCount"`
	MemberItem  []MemberItemModel `json:"MemberList"`
}

type TemplateModel struct {
	Text  []string          `json:"text"`
	Props map[string]string `json:"props"`
}

/*
  数据库模型
*/
type MysqlUser struct {
	ID         int64   `db:"id" json:"id"`
	Weid       string  `db:"weid" json:"weid"`
	Wename     string  `db:"wename" json:"wename"`
	Money      float64 `db:"money" json:"money"`
	CheckTime  int64   `db:"check_time" json:"check_time"`
	UpdateTime int64   `db:"update_time" json:"update_time"`
	CreateTime int64   `db:"create_time" json:"create_time"`
}

type MysqlFanli struct {
	ID         int64   `db:"id" json:"id"`
	Weid       string  `db:"weid" json:"weid"`
	GoodsID    int64   `db:"goodsid" json:"goodsid"`
	Goodstitle string  `db:"goodstitle" json:"goodstitle"`
	Money      float64 `db:"money" json:"money"`
	Bili       float64 `db:"bili" json:"bili"`
	Pid        string  `db:"pid" json:"pid"`
	Msgtext    string  `db:"msgtext" json:"msgtext"`
	UpdateTime int64   `db:"update_time" json:"update_time"`
	CreateTime int64   `db:"create_time" json:"create_time"`
}

type MysqlOrders struct {
	ID         int64   `db:"id" json:"id"`
	OrderID    int64   `db:"orderid" json:"orderid"`
	Weid       string  `db:"weid" json:"weid"`
	Goodsname  string  `db:"goodsname" json:"goodsname"`
	GoodsID    int64   `db:"goodsid" json:"goodsid"`
	Paymoney   float64 `db:"paymoney" json:"paymoney"`
	Pubmoney   float64 `db:"pubmoney" json:"pubmoney"`
	Pubbili    float64 `db:"pubbili" json:"pubbili"`
	Buymoney   float64 `db:"buymoney" json:"buymoney"`
	Buybili    float64 `db:"buybili" json:"buybili"`
	Income     float64 `db:"income" json:"income"`
	Status     int64   `db:"status" json:"status"`
	Pid        string  `db:"pid" json:"pid"`
	UpdateTime int64   `db:"update_time" json:"update_time"`
	CreateTime int64   `db:"create_time" json:"create_time"`
}

type MysqlWithdraw struct {
	ID         int64   `db:"id" json:"id"`
	Weid       string  `db:"weid" json:"weid"`
	Money      float64 `db:"money" json:"money"`
	Status     int64   `db:"status" json:"status"`
	UpdateTime int64   `db:"update_time" json:"update_time"`
	CreateTime int64   `db:"create_time" json:"create_time"`
}

type MysqlConfig struct {
	ID    int64  `db:"id" json:"id"`
	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`
}
