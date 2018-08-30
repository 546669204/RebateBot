package main

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"

	"github.com/546669204/RebateBot/common"
	httpdo "github.com/546669204/golang-http-do"
	"github.com/fsnotify/fsnotify"
	"github.com/tidwall/gjson"
)

var methods = map[string]interface{}{
	"ExitHook":    ExitHook,
	"msgprocess":  msgprocess,
	"islogin":     islogin,
	"parsexls":    parsexls,
	"initcontact": initcontact,
}
var Client *common.Client

var MemberList map[string]MemberListModel

var TemplateList map[string]TemplateModel
var EmojiTo gjson.Result
var sess sqlbuilder.Database
var StatusTo map[string]int64

type SettingModel struct {
	TuLing bool //图灵机器人
}

var Setting SettingModel

func main() {

	go initWebApi()

	StatusTo = make(map[string]int64)
	StatusTo["订单付款"] = 0
	StatusTo["订单结算"] = 1
	StatusTo["订单失效"] = 2

	var err error
	common.LogToFile()

	MemberList = make(map[string]MemberListModel)
	TemplateList = make(map[string]TemplateModel)
	InitConfig()
	var watcher *fsnotify.Watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case <-watcher.Events:
				//log.Println("修改文件：" + event.Name)
				//log.Println("修改类型：" + event.Op.String())
				log.Println("配置更新")
				InitConfig()
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filepath.Join(filepath.Dir(os.Args[0]), "template.json"))
	if err != nil {
		log.Fatal(err)
	}

	//setting
	setjsstr, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "database.json"))
	if err != nil {
		log.Println("读取配置项出错.程序自动退出")
		os.Exit(1)
	}
	Setting.TuLing = gjson.ParseBytes(setjsstr).Get("tuling").Bool()

	//Mysql 连接
	initMysqlSetting()
	sess, err = mysql.Open(mysqlsettings)
	if err != nil {
		log.Fatal("Conn ==> ", err)
	}
	defer sess.Close()

	//版本检测
	databaseVersioCheck()
	//
	common.LogToFile()
	Client = common.InitClient()
	defer Client.Conn.Close()
	Client.InitMethods(methods)
	Client.SericeInit()
	go func() {
		Client.HeartBeat()
	}()
	Client.ServiceHandle(Client.Conn)
}
func InitConfig() {
	//初始化配置
	jsstr, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "template.json"))
	if err != nil {
		log.Println("读取配置项出错.程序自动退出")
		os.Exit(1)
	}
	err = json.Unmarshal(jsstr, &TemplateList)
	if err != nil {
		log.Println("解析配置项出错.程序自动退出" + err.Error())
		os.Exit(1)
	}
}
func ExitHook() {
	//退出程序钩子
}
func msgprocess(data common.Msg) {
	b := gjson.Parse(data.Data)

	if b.Get("type").String() == "text" {
		takler := b.Get("takler").String()
		content := b.Get("content").String()

		if strings.Index(takler, "@chatroom") == -1 {

			serachkey := ParseContentUrl(content)
			if serachkey == "" {
				serachkey = ParseContentTkl(content)
			}
			if serachkey != "" || content == "签到" {

				if takler == "" || !isregisteruser(takler) {
					nickname := getnicknamefromusername(data.To, takler)
					registeruser(takler, nickname)
				}
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1500)+800))
			if serachkey != "" {
				var msg common.Msg
				msg.To = data.To
				fanli, _ := tktuiguang(serachkey, takler)
				/*msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, "../"+pic, takler)
				msg.Method = "sendfilemsg"
				Client.ConnWrite(msg)*/
				//time.Sleep(200 * time.Microsecond)
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, fanli, takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				//os.Remove(pic)
				return
			}
			if content == "签到" {
				nickname := getnicknamefromusername(data.To, takler)
				var msg common.Msg
				msg.To = data.To
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(checkin(data.To, takler, nickname)), takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				return
			}
			if content == "帮助" {
				var msg common.Msg
				msg.To = data.To

				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(templatebuild("help", nil)), takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				return
			}
			if content == "提现" {
				var msg common.Msg
				msg.To = data.To
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(withdrawmoney(takler)), takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				return
			}
			if content == "余额" {
				var msg common.Msg
				msg.To = data.To
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(userinfo(takler)), takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				return
			}
			reg := regexp.MustCompile(`\d{17,}`)
			matches := reg.FindStringSubmatch(content)
			if len(matches) == 1 {
				var msg common.Msg
				msg.To = data.To
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(bindorder(takler, matches[0])), takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				return
			}
			if content == "newfriend" {
				var msg common.Msg
				msg.To = data.To
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, fmt.Sprintf(`你好%s欢迎来到返利机器人`, takler), takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				return
			}
			if Setting.TuLing {
				var msg common.Msg
				msg.To = data.To
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, TuLingRobot(takler, content), takler)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
				return
			}
			log.Println(takler, content)
		}
	}
}
func tktuiguang(s, weid string) (string, string) {
	var data common.Msg
	data.Data = fmt.Sprintf(`{"keyword":"%s"}`, s)
	data.Method = "tbsearch"
	p := Client.ConnWriteReturn(data)
	if gjson.Parse(p.Data).Get("status").Int() != 1 {
		return `这个商品暂时没有返利哦`, ""
	}
	pro := gjson.Parse(p.Data).Get("products")
	GoodsID := pro.Get("ID").Int()
	/*
	   条件 商品id 生成时间小于24小时
	   获取 createtime weid pid  msgtext
	   1.如果 weid 相同 直接调用msgtext
	   2.如果 weid 不同 获取pid 排除pid 生成新的链接
	   3.如果没有数据 随意生成
	*/
	var AllFanli []MysqlFanli
	var where db.Cond
	where = make(db.Cond)
	where["goodsid"] = GoodsID
	where["create_time >="] = time.Now().Add(time.Hour * -24).Unix()
	err := sess.Collection("xm_fanli").Find(where).OrderBy("create_time desc").All(&AllFanli)
	if err != nil {
		log.Fatal("tktuiguang", err)
	}

	var ppp []string
	if len(AllFanli) >= 1 {
		for _, value := range AllFanli {
			if value.Weid == weid {
				return value.Msgtext, ""
			} else {
				ppp = append(ppp, value.Pid)
			}
		}

	}

	data.Data = fmt.Sprintf(`{"id":"%d","pid":"%s"}`, GoodsID, strings.Join(ppp, ","))
	data.Method = "settuiguang"
	l := Client.ConnWriteReturn(data)
	if gjson.Parse(l.Data).Get("status").Int() != 1 {
		return `这个商品暂时没有返利哦`, ""
	}
	link := gjson.Parse(l.Data).Get("link")
	pid := gjson.Parse(l.Data).Get("pid").String()
	tkl := link.Get("TaoToken").String()
	tlj := link.Get("ShortLinkUrl").String()
	if len(link.Get("CouponLinkTaoToken").String()) > 3 {
		tkl = link.Get("CouponLinkTaoToken").String()
		tlj = link.Get("CouponShortLinkUrl").String()
	}
	jiage, _ := strconv.ParseFloat(pro.Get("Jiage").String(), 10)

	youhui, _ := strconv.ParseFloat(pro.Get("Youhui").String(), 10)

	/*
		2018-05-28 暂时关闭下载
		//下载图片
		op := httpdo.Default()
		op.Url = "http:" + pro.Get("Pic").String()
		httpbyte, _ := httpdo.HttpDo(op)
		filename := fmt.Sprintf(`%d.jpg`, time.Now().UnixNano())
		file, _ := os.Create(filename)
		defer file.Close()
		file.Write(httpbyte)
	*/

	t := time.Now().Unix()
	bili := pro.Get("Bili").Float()
	if bili > 1 {
		bili = bili / 2
	}
	var tdata = make(map[string]interface{})
	tdata["Title"] = pro.Get("Title").String()
	tdata["TaoKouLing"] = tkl
	tdata["YuanJia"] = jiage
	tdata["FuKuan"] = Round(jiage-youhui, 2)
	tdata["FanLi"] = Round((jiage-youhui)*bili/100, 2)
	tdata["YouHui"] = youhui
	tdata["Link"] = tlj
	msgtext := templatebuild("fanli", tdata)

	var NewData MysqlFanli
	NewData.Weid = weid
	NewData.GoodsID = GoodsID
	NewData.Goodstitle = pro.Get("Title").String()
	NewData.CreateTime = t
	NewData.Pid = pid
	NewData.Msgtext = msgtext
	NewData.Money = Round((jiage-youhui)*bili/100, 2)
	NewData.Bili = bili

	_, err = sess.Collection("xm_fanli").Insert(NewData)
	if err != nil {
		log.Println("插入fanli错误", err)
	}

	return msgtext, "" //filename
}
func isaliduanurl(s string) bool {
	u, _ := url.Parse(s)
	if u == nil {
		return false
	}
	if u.Host == "www.dwntme.com" || u.Host == "dwntme.com" || (len(s) >= 19 && s[:17] == "http://m.tb.cn/h.") {
		return true
	}
	return false
}
func GetOriginalFromdwntme(s string) string {
	op := httpdo.Default()
	op.Url = s
	httpbyte, _ := httpdo.HttpDo(op)
	reg := regexp.MustCompile(`var url = '(.+)';`)
	matches := reg.FindStringSubmatch(string(httpbyte))
	if len(matches) == 2 {
		return GetOriginalUrl(matches[1])
	}
	return ""
}
func GetOriginalUrl(s string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", s, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	return resp.Header.Get("Location")
}

func tklParse(str string) string {
	op := httpdo.Default()
	op.Url = "http://api.chaozhi.hk/tb/tklParse"
	op.Method = "POST"
	//str = strings.Replace(str,"￥","%EF%BF%A5",-1)
	op.Data = "tkl=" + str
	httpbyte, err := httpdo.HttpDo(op)
	op.Header = "Content-Type:application/json; charset=utf-8"
	if err != nil {
		log.Println("tklParse", err)
	}
	returnjson := gjson.ParseBytes(httpbyte)
	if returnjson.Get("data").Get("suc").Bool() {
		taobaourl := returnjson.Get("data").Get("url").String()
		return GetOriginalUrl(taobaourl)
	}
	log.Println("tklParse", "口令无效", str)
	return ""
}
func ParseContentUrl(str string) string {
	reg := regexp.MustCompile(`((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+`)
	matches := reg.FindStringSubmatch(str)
	if len(matches) >= 3 {
		if isaliduanurl(matches[0]) {
			return GetOriginalFromdwntme(matches[0])
		}
		if IsNormalUrl(matches[0]) {
			return matches[0]
		}
	}
	return ""
}
func IsNormalUrl(str string) bool {
	u, _ := url.Parse(str)
	if u == nil {
		return false
	}
	if u.Host == "detail.tmall.com" || u.Host == "item.taobao.com" {
		return true
	}
	return false
}
func ParseContentTkl(str string) string {
	reg := regexp.MustCompile(`￥(\w{11})￥`)
	matches := reg.FindStringSubmatch(str)
	if len(matches) >= 2 {
		return tklParse(matches[0])
	}
	return ""
}

func islogin(data common.Msg) {
	data.Data = `{"isrun":1,"islogin":0,"name":"","lastlogin":"","runid":` + os.Args[1] + `}`
	data.Method = "msgreturn"
	Client.ConnWrite(data)
}

func isregisteruser(weid string) bool {
	count, err := sess.Collection("xm_user").Find("weid", weid).Count()
	if err != nil {
		log.Fatal(err)
		return true
	}
	if count == 0 {
		return false
	}
	return true
}

func registeruser(weid, wename string) bool {
	t := time.Now().Unix()
	var NewData MysqlUser
	NewData.Weid = weid
	NewData.Wename = wename
	NewData.CreateTime = t
	NewData.UpdateTime = t
	_, err := sess.Collection("xm_user").Insert(NewData)
	if err != nil {
		log.Fatal("registeruser", err)
		return false
	}
	return true
}

func getnicknamefromusername(runid, username string) string {
	for _, item := range MemberList[runid].MemberItem {
		if item.UserName == username {
			return item.NickName
		}
	}
	return ""
}

func checkin(runid, weid, nickname string) string {
	var OneUser MysqlUser
	FindUser := sess.Collection("xm_user").Find("weid", weid)
	err := FindUser.One(&OneUser)
	if err != nil {
		log.Fatal("getnewweid", err)
		return ""
	}
	if (time.Now().Unix()-OneUser.CheckTime)/(60*60) < 24 {
		return "您上次签到到现在还没有24小时呢 请稍后重试。"
	}

	var seed int64
	binary.Read(crand.Reader, binary.LittleEndian, &seed)
	rand.Seed(seed)
	m := rand.Float64()
	m = Round(m, 2)
	t := time.Now().Unix()
	OneUser.Money = Round(m+OneUser.Money, 2)
	OneUser.CheckTime = t
	err = FindUser.Update(OneUser)
	if err != nil {
		log.Fatal("checkin", err)
		return ""
	}
	var data = make(map[string]interface{})
	data["UserName"] = nickname
	data["CheckMoney"] = m
	data["Money"] = OneUser.Money
	return templatebuild("check", data)

}

func templatebuild(templatename string, data interface{}) string {
	t, ok := TemplateList[templatename]
	if !ok {
		log.Println("模板不存在", TemplateList)
		return "模板不存在"
	}
	for _, value := range t.Text {
		tem := template.New("fieldname example")
		tem, _ = tem.Parse(value)
		b := bytes.NewBuffer(make([]byte, 0))
		tem.Execute(b, data)
		return b.String()
	}
	return "模板不存在"
}

func MsgTextProcess(s string) string {
	if !EmojiTo.Exists() {
		jsstr, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "emoji.json"))
		if err != nil {
			log.Println("读取配置项出错.程序自动退出")
			os.Exit(1)
		}
		EmojiTo = gjson.ParseBytes(jsstr)
	}

	reg := regexp.MustCompile(`<([^>]+)>`)
	matchs := reg.FindAllString(s, -1)
	ret := s
	sort.Strings(matchs)

	for _, value := range Duplicate(matchs) {
		u := EmojiTo.Get("QQFaceMap").Get(value.(string)).String()
		b := EmojiTo.Get("EmojiCodeMap").Get(u).String()
		ret = strings.Replace(ret, value.(string), b, -1)
	}
	return ret
}

func Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

func withdrawmoney(weid string) string {
	var OneUser MysqlUser
	tx, _ := sess.NewTx(nil)
	FindUser := tx.Collection("xm_user").Find("weid", weid)
	err := FindUser.One(&OneUser)
	if err != nil {
		log.Fatal("withdrawmoney", err)
		return ""
	}
	if OneUser.Money < 5 {
		var data = make(map[string]interface{})
		data["Money"] = OneUser.Money
		return templatebuild("withdrawoff", data)
	}
	t := time.Now().Unix()
	var NewData MysqlWithdraw
	NewData.Money = OneUser.Money
	NewData.Status = 0
	NewData.Weid = weid
	NewData.CreateTime = t

	OneUser.Money = 0.00
	OneUser.UpdateTime = t

	FindUser.Update(OneUser)
	sess.Collection("xm_withdraw").Insert(NewData)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err)
		return `申请失败，未知错误。`
	}

	var data = make(map[string]interface{})
	data["Time"] = time.Now().Format("2006-01-02 15:04:05")
	data["Money"] = NewData.Money
	return templatebuild("withdrawok", data)
}

func parsexls(data common.Msg) {
	jdata := gjson.Parse(data.Data).Array()
	if jdata == nil {
		log.Println("解析xls 失败", data.Data)
	}
	var where db.Cond
	for i := 0; i < len(jdata); i++ {
		where = make(db.Cond)
		where["orderid"] = jdata[i].Get("订单编号").Int()
		where["goodsid"] = jdata[i].Get("商品ID").Int()
		FindOrders := sess.Collection("xm_orders").Find(where)
		count, err := FindOrders.Count()
		if err != nil {
			log.Println("parsexls sql错误", err)
		}
		if count > 0 {
			var NewData MysqlOrders
			NewData.Status = StatusTo[jdata[i].Get("订单状态").String()]

			FindOrders.Update(NewData)
			if NewData.Status == 1 { //订单完成 增加余额
				var OneOrders MysqlOrders
				FindOrders.One(&OneOrders)
				var OneUser MysqlUser
				FindUser := sess.Collection("xm_user").Find("weid", OneOrders.Weid)
				FindUser.One(&OneUser)
				OneUser.Money = Round(OneUser.Money+OneOrders.Buymoney, 2)
				FindUser.Update(OneUser)
				if len(OneOrders.Weid) == 0 {
					continue
				}
				rid, _ := weidTousername(OneOrders.Weid)
				var msg common.Msg
				var tdata = make(map[string]interface{})
				tdata["Title"] = OneOrders.Goodsname
				tdata["Buymoney"] = OneOrders.Buymoney
				tdata["Money"] = OneUser.Money
				tdata["OrderID"] = OneOrders.OrderID
				content := templatebuild("shouhuo", tdata)
				msg.To = rid
				msg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(content), OneOrders.Weid)
				msg.Method = "sendtextmsg"
				Client.ConnWrite(msg)
			}
			if NewData.Status == 2 {
				var OneOrders MysqlOrders
				FindOrders.One(&OneOrders)
				if OneOrders.Status == 1 {
					var OneUser MysqlUser
					FindUser := sess.Collection("xm_user").Find("weid", OneOrders.Weid)
					OneUser.Money = Round(OneUser.Money-OneOrders.Buymoney, 2)
					FindUser.Update(OneUser)
					OneOrders.Status = 2
					FindOrders.Update(OneOrders)
				}

			}
			continue
		}
		t := time.Now().Unix()

		var NewData MysqlOrders
		NewData.OrderID = jdata[i].Get("订单编号").Int()
		NewData.GoodsID = jdata[i].Get("商品ID").Int()
		NewData.Goodsname = jdata[i].Get("商品信息").String()
		NewData.Paymoney = jdata[i].Get("付款金额").Float()
		NewData.Pubmoney = jdata[i].Get("效果预估").Float()
		reg := regexp.MustCompile(`(\%| )`)
		bili, _ := strconv.ParseFloat(reg.ReplaceAllString("2.50 %", ""), 10)
		NewData.Pubbili = bili
		NewData.Buymoney = 0
		NewData.Buybili = 0
		NewData.Income = NewData.Pubmoney
		NewData.Status = StatusTo[jdata[i].Get("订单状态").String()]
		NewData.Pid = jdata[i].Get("来源媒体ID").String() + "-" + jdata[i].Get("广告位ID").String()
		NewData.UpdateTime = t
		NewData.CreateTime = t

		//自动绑定

		where = make(db.Cond)
		where["goodsid"] = NewData.GoodsID
		where["create_time >="] = time.Now().Add(time.Hour * -24).Unix()
		where["pid"] = NewData.Pid
		var AllFanli []MysqlFanli
		sess.Collection("xm_fanli").Find(where).All(&AllFanli)
		if len(AllFanli) > 0 {
			if len(AllFanli) == 1 {
				NewData.Buymoney = Round(AllFanli[0].Bili*NewData.Paymoney/100, 2)
				NewData.Buybili = AllFanli[0].Bili
				NewData.Income = Round(NewData.Pubmoney-NewData.Buymoney, 2)
				NewData.Weid = AllFanli[0].Weid
				var tdata = make(map[string]interface{})
				tdata["Title"] = NewData.Goodsname
				tdata["Paymoney"] = NewData.Paymoney
				tdata["Money"] = NewData.Buymoney
				tdata["OrderID"] = NewData.OrderID
				content := templatebuild("autobindorder", tdata)
				rid, _ := weidTousername(NewData.Weid)
				var bindmsg common.Msg
				bindmsg.To = rid
				bindmsg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(content), NewData.Weid)
				bindmsg.Method = "sendtextmsg"
				Client.ConnWrite(bindmsg)
			} else {
				var tdata = make(map[string]interface{})
				tdata["Title"] = NewData.Goodsname
				tdata["Paymoney"] = NewData.Paymoney
				content := templatebuild("bindorder", tdata)
				var bindmsg common.Msg
				for _, value := range AllFanli {
					rid, _ := weidTousername(value.Weid)
					bindmsg.To = rid
					bindmsg.Data = fmt.Sprintf(`{"content":"%s","to":"%s"}`, MsgTextProcess(content), value.Weid)
					bindmsg.Method = "sendtextmsg"
					Client.ConnWrite(bindmsg)
				}
			}
		}
		sess.Collection("xm_orders").Insert(NewData)
	}
}

func ToJson(j interface{}) string {
	js, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(js)
}

func bindorder(weid, orderid string) string {
	var OneOrder MysqlOrders
	FindOrders := sess.Collection("xm_orders").Find("orderid", orderid)
	err := FindOrders.One(&OneOrder)
	if err == db.ErrNoMoreRows {
		return ""
	}
	if err != nil {
		log.Println(err)
		return ""
	}
	if OneOrder.Weid != "" {
		return "该笔订单已经被人绑定，如需投诉请联系技术。"
	}

	var OneUser MysqlUser
	sess.Collection("xm_user").Find("weid", weid).One(&OneUser)
	var where db.Cond
	where = make(db.Cond)
	where["weid"] = weid
	where["goodsid"] = OneOrder.GoodsID
	var OneFanli MysqlFanli
	sess.Collection("xm_fanli").Find(where).OrderBy("create_time desc").One(&OneFanli)

	OneOrder.Buymoney = Round(OneFanli.Bili*OneOrder.Paymoney/100, 2)
	OneOrder.Buybili = OneFanli.Bili
	OneOrder.Income = Round(OneOrder.Pubmoney-OneOrder.Buymoney, 2)
	OneOrder.Weid = weid
	FindOrders.Update(OneOrder)

	var data = make(map[string]interface{})
	data["Title"] = OneOrder.Goodsname
	data["Buymoney"] = OneOrder.Buymoney
	return templatebuild("bindorderok", data)
}

func userinfo(weid string) string {
	var OneUser MysqlUser
	var AllOrder []MysqlOrders
	var where db.Cond
	where = make(db.Cond)
	sess.Collection("xm_user").Find("weid", weid).One(&OneUser)
	where["weid"] = weid
	where["status"] = 0
	err := sess.Collection("xm_orders").Find(where).All(&AllOrder)
	var totalmoney float64
	totalmoney = 0
	if err != nil {
		log.Println("userinfo", err)
		return ""
	}
	for _, k := range AllOrder {
		totalmoney += k.Buymoney
	}
	var data = make(map[string]interface{})
	data["Name"] = OneUser.Wename
	data["Money"] = OneUser.Money
	data["Dmoney"] = Round(totalmoney, 2)
	return templatebuild("money", data)
}

func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func weidTousername(weid string) (string, string) {
	for k, v := range MemberList {
		for _, v2 := range v.MemberItem {
			if v2.RemarkName == weid {
				return k, v2.UserName
			}
		}
	}
	return "", ""
}

func TuLingRobot(id string, msg string) string {
	id = strings.Replace(id, "_", "", -1)
	op := httpdo.Default()
	op.Url = "http://openapi.tuling123.com/openapi/api/v2"
	op.Method = "post"
	op.Data = "{\"reqType\":0,\"perception\":{\"inputText\":{\"text\":\"" + msg + "\"}},\"userInfo\":{\"apiKey\":\"141f94237af141918cbfdeaa0323d480\",\"userId\":\"" + id + "\"}}"
	result, _ := httpdo.HttpDo(op)
	return gjson.Parse(string(result)).Get("results").Array()[0].Get("values").Get("text").String()
}
func initcontact(data common.Msg) {
	b := gjson.Parse(data.Data)
	var list MemberListModel
	list.MemberCount = len(b.Array())
	for _, v := range b.Array() {
		var item MemberItemModel
		item.UserName = v.Get("username").String()
		item.NickName = v.Get("nickname").String()

		list.MemberItem = append(list.MemberItem, item)
	}

	MemberList[data.To] = list
}
